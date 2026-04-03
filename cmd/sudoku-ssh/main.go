package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"charm.land/wish/v2"
	"charm.land/wish/v2/activeterm"
	wishbubbletea "charm.land/wish/v2/bubbletea"

	"github.com/PhilippSchweizer/sudoku-engine/internal/tui"
)

func main() {
	addr := flag.String("addr", envOr("SUDOKU_SSH_ADDR", ":2222"), "listen address (e.g. :2222 or 0.0.0.0:22)")
	hostKey := flag.String("host-key", envOr("SUDOKU_SSH_HOST_KEY", ""), "path to SSH host private key (PEM); created if missing. Empty uses .sudoku-ssh/host_ed25519 in the current directory")
	flag.Parse()

	keyPath := *hostKey
	if keyPath == "" {
		keyPath = filepath.Join(".sudoku-ssh", "host_ed25519")
	}
	if err := os.MkdirAll(filepath.Dir(keyPath), 0o700); err != nil {
		log.Fatalf("host key directory: %v", err)
	}

	s, err := wish.NewServer(
		wish.WithAddress(*addr),
		wish.WithHostKeyPath(keyPath),
		// First entry becomes innermost after composition; outer runs first per session.
		wish.WithMiddleware(
			wishbubbletea.Middleware(tui.NewRemoteSessionModel),
			activeterm.Middleware(),
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("sudoku-ssh listening on ssh://%s (host key %s)", *addr, keyPath)
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
