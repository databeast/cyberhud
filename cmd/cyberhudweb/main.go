// Command cyberhudweb serves a minimal web control UI for the cyberhud daemon.
//
// It is a peer to cyberhudctl: it speaks the same console Unix-domain socket
// protocol and is just another client of the same commands. It renders a
// server-side page showing display regions, their current and available modes,
// and a to-scale layout of the virtual display, and lets the user switch modes.
//
// Security model: this is a Pi-local tool. There is no authentication — binding
// the address is the access control ("run it or don't"). Bind to a loopback
// address (the default) unless you intend LAN access.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	socketPath := flag.String("socket", "/run/cyberhudd/console.sock",
		"path to the cyberhudd console Unix domain socket")
	addr := flag.String("addr", "127.0.0.1:8080",
		"HTTP bind address (host:port); use 0.0.0.0:PORT for LAN access")
	timeout := flag.Duration("timeout", 2*time.Second,
		"per-command socket timeout")
	dumpState := flag.Bool("dump-state", false,
		"connect to the daemon, print the parsed read model, and exit (diagnostic)")
	flag.Parse()

	client := NewClient(*socketPath, *timeout)

	if *dumpState {
		if err := runDumpState(client); err != nil {
			fmt.Fprintf(os.Stderr, "ERR %v\n", err)
			os.Exit(1)
		}
		return
	}

	srv := NewServer(client)
	log.SetFlags(log.Ldate | log.Ltime | log.Lmsgprefix)
	log.SetPrefix("[cyberhudweb] ")
	log.Printf("serving on http://%s (daemon socket: %s)", *addr, *socketPath)
	if err := srv.ListenAndServe(*addr); err != nil {
		log.Fatalf("http server: %v", err)
	}
}

// runDumpState fetches and prints the read model for diagnostics (Task 3 demo).
func runDumpState(f Fetcher) error {
	st, err := f.FetchState()
	if err != nil {
		return err
	}
	fmt.Printf("virtual display: %dx%d\n", st.Geometry.VirtualW, st.Geometry.VirtualH)
	fmt.Printf("regions (%d):\n", len(st.Regions))
	for _, r := range st.Regions {
		fmt.Printf("  [%d] %s (id=%s controller=%s) current=%s\n",
			r.Index, r.Name, r.ID, r.Controller, r.Current)
		for _, m := range r.Modes {
			if m.Title != "" {
				fmt.Printf("      - %s  %q\n", m.ID, m.Title)
			} else {
				fmt.Printf("      - %s\n", m.ID)
			}
		}
	}
	for _, g := range st.Geometry.Regions {
		fmt.Printf("  geom %s: x=%d y=%d w=%d h=%d\n", g.Name, g.X, g.Y, g.W, g.H)
	}
	return nil
}
