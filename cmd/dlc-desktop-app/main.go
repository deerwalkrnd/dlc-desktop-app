// Command dlc-desktop-app serves the Deerwalk Learning Center video library
// to a local web browser.
//
// It scans the DLC folder next to the executable, builds the catalogue in
// memory and serves the embedded web app on http://localhost:3000.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

	"github.com/deerwalkrnd/dlc-desktop-app/internal/library"
	"github.com/deerwalkrnd/dlc-desktop-app/internal/server"
	"github.com/deerwalkrnd/dlc-desktop-app/web"
)

// version is set at build time with -ldflags "-X main.version=...".
var version = "dev"

func main() {
	if err := run(); err != nil {
		slog.Error(err.Error())
		if runtime.GOOS == "windows" {
			// Keep the console window open long enough to read the error
			// when the app was started by double-clicking.
			fmt.Fprintln(os.Stderr, "\nPress Enter to close this window.")
			fmt.Scanln()
		}
		os.Exit(1)
	}
}

func run() error {
	var (
		addr        = flag.String("addr", "127.0.0.1:3000", "address to listen on; use 0.0.0.0:3000 to share on the local network")
		libraryPath = flag.String("library", "", "video library folder (default: DLC next to the executable)")
		openBrowser = flag.Bool("open", true, "open the app in the default web browser")
		verbose     = flag.Bool("verbose", false, "log debug details")
		showVersion = flag.Bool("version", false, "print the version and exit")
	)
	flag.Parse()

	if *showVersion {
		fmt.Println(version)
		return nil
	}

	level := slog.LevelInfo
	if *verbose {
		level = slog.LevelDebug
	}
	log := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level}))
	slog.SetDefault(log)

	root, err := resolveLibrary(*libraryPath)
	if err != nil {
		return err
	}

	log.Info("starting DLC Desktop App", "version", version, "library", root)
	lib := library.New(root)
	x, err := lib.Rescan()
	if err != nil {
		return fmt.Errorf("cannot read the library folder: %w", err)
	}
	if !x.Exists {
		log.Warn("library folder not found; creating it. Put your .mp4 files there and press Rescan in the app", "path", root)
		if err := os.MkdirAll(root, 0o755); err != nil {
			log.Warn("could not create library folder", "err", err)
		}
	}
	log.Info("library scanned", "lessons", x.LessonCount(), "classes", len(x.Classes), "skipped", len(x.Skipped), "warnings", len(x.Warnings), "took", x.Duration.Round(time.Millisecond))
	for _, s := range x.Skipped {
		log.Warn("skipped file", "path", s.Path, "reason", s.Message)
	}

	hosts := server.LoopbackHosts(*addr)
	if hosts == nil {
		log.Warn("listening on a non-loopback address; the library is visible to the local network", "addr", *addr)
	}
	srv, err := server.New(lib, server.Options{UI: web.Dist(), AllowedHosts: hosts, Logger: log, Version: version})
	if err != nil {
		return err
	}

	url := browserURL(*addr)
	ln, err := net.Listen("tcp", *addr)
	if err != nil {
		if errors.Is(err, syscall.EADDRINUSE) || isAddrInUse(err) {
			if alreadyRunning(url) {
				log.Info("DLC Desktop App is already running", "url", url)
				if *openBrowser {
					open(url)
				}
				return nil
			}
			return fmt.Errorf("port %s is used by another program; close it or start with -addr 127.0.0.1:3001", *addr)
		}
		return err
	}

	httpSrv := srv.HTTPServer()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	errc := make(chan error, 1)
	go func() { errc <- httpSrv.Serve(ln) }()
	log.Info("ready - open this address in your browser", "url", url)
	if *openBrowser {
		open(url)
	}

	select {
	case err := <-errc:
		return err
	case <-ctx.Done():
	}
	log.Info("shutting down")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return httpSrv.Shutdown(shutdownCtx)
}

// resolveLibrary picks the library folder: the -library flag if given,
// otherwise DLC next to the executable, falling back to DLC in the working
// directory (useful with `go run`, whose executable lives in a temp folder).
func resolveLibrary(flagValue string) (string, error) {
	if flagValue != "" {
		return filepath.Abs(flagValue)
	}
	var candidates []string
	if exe, err := os.Executable(); err == nil {
		if exe, err = filepath.EvalSymlinks(exe); err == nil {
			candidates = append(candidates, filepath.Join(filepath.Dir(exe), "DLC"))
		}
	}
	if wd, err := os.Getwd(); err == nil {
		candidates = append(candidates, filepath.Join(wd, "DLC"))
	}
	if len(candidates) == 0 {
		return "", errors.New("cannot determine the library folder; pass -library")
	}
	for _, c := range candidates {
		if fi, err := os.Stat(c); err == nil && fi.IsDir() {
			return c, nil
		}
	}
	return candidates[0], nil
}

func browserURL(addr string) string {
	_, port, err := net.SplitHostPort(addr)
	if err != nil {
		port = "3000"
	}
	return "http://localhost:" + port
}

func isAddrInUse(err error) bool {
	var se *os.SyscallError
	if errors.As(err, &se) {
		// WSAEADDRINUSE on Windows.
		if errno, ok := se.Err.(syscall.Errno); ok && errno == 10048 {
			return true
		}
	}
	return false
}

// alreadyRunning reports whether another copy of this app answers at url.
func alreadyRunning(url string) bool {
	client := http.Client{Timeout: time.Second}
	resp, err := client.Get(url + "/healthz")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	var body struct {
		App string `json:"app"`
	}
	return json.NewDecoder(resp.Body).Decode(&body) == nil && body.App == "dlc-desktop-app"
}

func open(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	if err := cmd.Start(); err != nil {
		slog.Debug("could not open browser", "err", err)
		return
	}
	go cmd.Wait()
}
