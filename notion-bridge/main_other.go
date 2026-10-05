//go:build !windows

package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
)

// Headless mode for development and tests on non-Windows systems.
func main() {
	dir := flag.String("dir", ".", "data directory")
	noNgrok := flag.Bool("no-ngrok", false, "do not start ngrok")
	flag.Parse()
	home, _ := os.UserHomeDir()
	m, err := NewManager(filepath.Join(*dir, "config.json"), home, filepath.Join(*dir, "logs"), func(s string) { fmt.Println(s) })
	if err != nil { fmt.Println(err); os.Exit(1) }
	if *noNgrok { m.mu.Lock(); m.cfg.NgrokDomain = ""; m.mu.Unlock() }
	if err := m.StartAll(); err != nil { fmt.Println("Warnungen:", err) }
	fmt.Println(m.StatusText())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	m.StopAll()
}
