package source

import (
	"context"
	"io"
	"os"
	"sync"
	"time"
)

// tailReader manages the lifecycle of a background goroutine that tails a
// configured data source (file, named pipe, or Unix domain socket) and pushes
// ingested bytes into a lineBuffer. It uses context-based cancellation for
// clean shutdown.
type TailReader struct {
	mu     sync.Mutex
	cancel context.CancelFunc
	active bool
	done   chan struct{}
}

// clampScanMS ensures ScanMS falls within the valid range [100, 30000].
func clampScanMS(ms int) int {
	if ms < 100 {
		return 100
	}
	if ms > 30000 {
		return 30000
	}
	return ms
}

// Start launches a background goroutine that opens the data source specified
// by policy.Source and continuously reads new bytes into buf. For regular files
// it seeks to the end before entering the read loop. On EOF or pipe close it
// sleeps for the clamped ScanMS duration then retries. On read errors it closes
// the source and retries opening at the ScanMS interval. If the source path
// does not exist at activation time, it retries at ScanMS until the path
// becomes available or the reader is stopped.
func (r *TailReader) Start(policy Policy, buf *LineBuffer) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.active {
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	r.cancel = cancel
	r.active = true
	r.done = make(chan struct{})

	scanMS := clampScanMS(policy.ScanMS)
	scanDur := time.Duration(scanMS) * time.Millisecond
	source := policy.Source

	go r.run(ctx, source, scanDur, buf)
}

// run is the main goroutine loop. It handles opening, reading, and retrying
// the data source until the context is cancelled.
func (r *TailReader) run(ctx context.Context, source string, scanDur time.Duration, buf *LineBuffer) {
	defer close(r.done)

	readBuf := make([]byte, 4096)

	for {
		if err := ctx.Err(); err != nil {
			return
		}

		f, err := r.openSource(ctx, source, scanDur)
		if err != nil {
			// Context was cancelled during open retry loop.
			return
		}

		r.readLoop(ctx, f, readBuf, buf, scanDur)
		f.Close()

		// If context is done, exit. Otherwise sleep before retrying open.
		if ctx.Err() != nil {
			return
		}
		if !r.sleep(ctx, scanDur) {
			return
		}
	}
}

// openSource attempts to open the data source file. If the file does not exist,
// it retries at the given interval until the file appears or the context is
// cancelled. Returns nil error only when a file is successfully opened.
func (r *TailReader) openSource(ctx context.Context, source string, scanDur time.Duration) (*os.File, error) {
	for {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		f, err := os.Open(source)
		if err != nil {
			if !os.IsNotExist(err) {
				// Non-existence and other errors: retry after sleep.
				if !r.sleep(ctx, scanDur) {
					return nil, ctx.Err()
				}
				continue
			}
			// File does not exist — retry at scan interval.
			if !r.sleep(ctx, scanDur) {
				return nil, ctx.Err()
			}
			continue
		}

		// Seek to end for regular files.
		info, err := f.Stat()
		if err == nil && info.Mode().IsRegular() {
			_, _ = f.Seek(0, io.SeekEnd)
		}

		return f, nil
	}
}

// readLoop continuously reads from the opened file, pushing data into the
// buffer. It returns when a non-EOF/non-temporary error occurs, or when the
// context is cancelled.
func (r *TailReader) readLoop(ctx context.Context, f *os.File, readBuf []byte, buf *LineBuffer, scanDur time.Duration) {
	for {
		if ctx.Err() != nil {
			return
		}

		n, err := f.Read(readBuf)
		if n > 0 {
			buf.Ingest(readBuf[:n])
		}

		if err != nil {
			if err == io.EOF {
				// EOF: sleep then retry read (tail behavior).
				if !r.sleep(ctx, scanDur) {
					return
				}
				continue
			}
			// Non-EOF error (connection lost, pipe closed, etc.): return to
			// trigger reopen in the outer loop.
			return
		}
	}
}

// sleep pauses for the given duration or until the context is cancelled.
// Returns true if the full duration elapsed, false if cancelled.
func (r *TailReader) sleep(ctx context.Context, d time.Duration) bool {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return false
	case <-t.C:
		return true
	}
}

// Stop cancels the background goroutine and waits for it to exit, then marks
// the reader as inactive.
func (r *TailReader) Stop() {
	r.mu.Lock()
	cancel := r.cancel
	done := r.done
	r.mu.Unlock()

	if cancel != nil {
		cancel()
	}
	if done != nil {
		<-done
	}

	r.mu.Lock()
	r.active = false
	r.cancel = nil
	r.done = nil
	r.mu.Unlock()
}

// IsActive reports whether the tail reader goroutine is currently running.
func (r *TailReader) IsActive() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.active
}
