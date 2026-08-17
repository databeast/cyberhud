package source

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
)

const (
	backoffInitial = 1 * time.Second
	backoffMax     = 30 * time.Second
)

func StartInstantMonitor(root string, out chan<- struct{}) error {
	if _, err := os.Stat(root); err != nil {
		return fmt.Errorf("scan root unavailable: %w", err)
	}
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("create watcher: %w", err)
	}
	if err := watcher.Add(root); err != nil {
		_ = watcher.Close()
		return fmt.Errorf("watch %s: %w", root, err)
	}

	// Watch existing immediate child subdirectories at startup (depth 1).
	watchExistingSubdirs(root, watcher)

	go watchLoop(root, watcher, out)
	return nil
}

// watchExistingSubdirs adds watches for all immediate child subdirectories
// of root that exist at startup time.
func watchExistingSubdirs(root string, watcher *fsnotify.Watcher) {
	entries, err := os.ReadDir(root)
	if err != nil {
		log.Printf("usb watcher: failed to read subdirectories of %s: %v", root, err)
		return
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		subdir := filepath.Join(root, entry.Name())
		if err := watcher.Add(subdir); err != nil {
			log.Printf("usb watcher: failed to watch subdirectory %s: %v", subdir, err)
		}
	}
}

// watchLoop is the main event loop goroutine for the filesystem watcher.
// It handles events, errors, and recovers from path unavailability using
// exponential backoff. It also manages subdirectory watches at depth 1.
func watchLoop(root string, watcher *fsnotify.Watcher, out chan<- struct{}) {
	defer watcher.Close()
	for {
		select {
		case ev, ok := <-watcher.Events:
			if !ok {
				// Events channel closed — attempt recovery.
				log.Printf("usb watcher: events channel closed, attempting recovery for %s", root)
				if !recoverWatch(root, &watcher, out) {
					return
				}
				continue
			}
			// Check if the scan root became unavailable after this event.
			if _, statErr := os.Stat(root); statErr != nil {
				log.Printf("usb watcher: scan root unavailable (%v), entering recovery", statErr)
				if !recoverWatch(root, &watcher, out) {
					return
				}
				continue
			}
			if ev.Op&(fsnotify.Create|fsnotify.Remove|fsnotify.Rename|fsnotify.Write) == 0 {
				continue
			}

			// Handle subdirectory watch management for depth-1 children.
			handleSubdirEvent(root, ev, watcher)

			select {
			case out <- struct{}{}:
			default:
			}
		case watchErr, ok := <-watcher.Errors:
			if !ok {
				// Errors channel closed — attempt recovery.
				log.Printf("usb watcher: errors channel closed, attempting recovery for %s", root)
				if !recoverWatch(root, &watcher, out) {
					return
				}
				continue
			}
			log.Printf("usb watcher: error on %s: %v", root, watchErr)
			// Check if the root is still accessible after the error.
			if _, statErr := os.Stat(root); statErr != nil {
				log.Printf("usb watcher: scan root unavailable after error (%v), entering recovery", statErr)
				if !recoverWatch(root, &watcher, out) {
					return
				}
			}
			// Otherwise log and continue monitoring.
		}
	}
}

// handleSubdirEvent manages subdirectory watches based on filesystem events.
// It adds a watch when a new immediate child directory is created under root,
// and removes the watch when a subdirectory is removed. Only depth 1 is monitored.
func handleSubdirEvent(root string, ev fsnotify.Event, watcher *fsnotify.Watcher) {
	// Only manage watches for direct children of the scan root.
	parentDir := filepath.Dir(ev.Name)
	if parentDir != root {
		// Event is from a watched subdirectory or deeper — no watch management needed.
		// The signal is still sent (handled by the caller).
		return
	}

	switch {
	case ev.Op&fsnotify.Create != 0:
		// New entry created directly under root — check if it's a directory.
		info, err := os.Stat(ev.Name)
		if err != nil {
			// Transient file or already gone — ignore gracefully.
			return
		}
		if !info.IsDir() {
			return
		}
		// Add watch for the new subdirectory (depth 1 only).
		if err := watcher.Add(ev.Name); err != nil {
			log.Printf("usb watcher: failed to watch new subdirectory %s: %v", ev.Name, err)
		}

	case ev.Op&(fsnotify.Remove|fsnotify.Rename) != 0:
		// Entry removed or renamed directly under root — attempt to remove watch.
		// watcher.Remove is a no-op if the path isn't watched, so this is safe
		// for both files and directories.
		_ = watcher.Remove(ev.Name)
	}
}

// recoverWatch implements exponential backoff retry to re-establish the watch
// on root. It replaces the watcher pointer on success. Returns true if recovery
// succeeded, false if unrecoverable (should not happen under normal operation
// since it retries indefinitely, but kept as a safety valve).
func recoverWatch(root string, watcher **fsnotify.Watcher, out chan<- struct{}) bool {
	// Close old watcher if still open.
	_ = (*watcher).Close()

	backoff := backoffInitial
	for {
		time.Sleep(backoff)

		// Check if path is available.
		if _, err := os.Stat(root); err != nil {
			log.Printf("usb watcher: recovery attempt failed, path %s still unavailable (%v), retrying in %v", root, err, backoff)
			backoff = nextBackoff(backoff)
			continue
		}

		// Path is available — create new watcher and add root.
		newWatcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.Printf("usb watcher: recovery failed to create watcher (%v), retrying in %v", err, backoff)
			backoff = nextBackoff(backoff)
			continue
		}
		if err := newWatcher.Add(root); err != nil {
			log.Printf("usb watcher: recovery failed to watch %s (%v), retrying in %v", root, err, backoff)
			_ = newWatcher.Close()
			backoff = nextBackoff(backoff)
			continue
		}

		// Recovery succeeded — replace watcher and signal immediate rescan.
		*watcher = newWatcher
		log.Printf("usb watcher: recovery successful, re-watching %s", root)

		// Re-watch existing subdirectories after recovery (depth 1).
		watchExistingSubdirs(root, newWatcher)

		select {
		case out <- struct{}{}:
		default:
		}
		return true
	}
}

// nextBackoff doubles the interval up to backoffMax.
func nextBackoff(current time.Duration) time.Duration {
	next := current * 2
	if next > backoffMax {
		return backoffMax
	}
	return next
}
