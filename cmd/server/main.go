package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ChaitanyaSaiV/Event-Ingestion/internal/handlers"
	"github.com/ChaitanyaSaiV/Event-Ingestion/internal/router"
	"github.com/ChaitanyaSaiV/Event-Ingestion/internal/storage"
)

func main() {
	// ─────────────────────────────────────────────────────
	// STEP 1: Wire up dependencies (what you already had)
	// ─────────────────────────────────────────────────────

	// Create the storage layer
	storeType := flag.String("store", "memory", "store backend: memory or file")
	flag.Parse()
	var store handlers.IncidentStore
	switch *storeType {
	case "memory":
		store = storage.NewInMemoryStore()
	case "file":
		fs, err := storage.NewFileStorage("incidents.json")
		if err != nil {
			log.Fatal(err)
		}
		store = fs
	default:
		log.Fatalf("unknown store type: %s", *storeType)
	}

	// Create the handler, inject the store
	handler := handlers.NewIncidentHandler(store)

	// Get the configured server (router.NewServer returns *http.Server)
	server := router.NewServer(":8080", handler)

	// ─────────────────────────────────────────────────────
	// STEP 2: Start the server in a BACKGROUND goroutine
	// ─────────────────────────────────────────────────────
	//
	// Why a goroutine?
	// server.ListenAndServe() BLOCKS forever — it never returns while
	// the server runs. If we called it directly in main, the next lines
	// of code would never execute. We need main to be free to listen
	// for shutdown signals, so we run the server in the background.
	//
	// "go" before a function call means "run this in a new goroutine
	// (lightweight thread) and don't wait for it to finish."

	go func() {
		log.Printf("server starting on %s", server.Addr)

		// ListenAndServe normally blocks. It only returns when:
		//   1. There's a real error (port already in use, etc.)
		//   2. server.Shutdown() is called — in which case it returns
		//      a SPECIAL error: http.ErrServerClosed
		//
		// We want to crash on real errors, but ignore ErrServerClosed
		// (because that's our intentional shutdown, not a failure).

		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}
	}()

	// ─────────────────────────────────────────────────────
	// STEP 3: Set up signal handling
	// ─────────────────────────────────────────────────────
	//
	// The OS sends "signals" to running programs. Two important ones:
	//   - os.Interrupt (SIGINT)  — what Ctrl+C sends
	//   - syscall.SIGTERM        — what Docker/Kubernetes send when
	//                              they want a graceful stop
	//
	// We want to catch BOTH and shut down cleanly.

	// Create a channel that can hold 1 signal value.
	// Channels are how goroutines communicate in Go.
	// We need a buffered channel (size 1) so that signal.Notify
	// can send the signal without blocking if we're not yet listening.
	quit := make(chan os.Signal, 1)

	// Tell Go: "when these signals come in, send them to my 'quit' channel"
	signal.Notify(quit, os.Interrupt, syscall.SIGINT)

	// ─────────────────────────────────────────────────────
	// STEP 4: Block here until a signal arrives
	// ─────────────────────────────────────────────────────
	//
	// <-quit means "receive a value from the quit channel."
	// This BLOCKS main until something is sent to the channel.
	// Since signal.Notify is the only sender, this effectively
	// blocks until Ctrl+C or SIGTERM.
	//
	// Meanwhile, our server (in the goroutine) keeps handling requests.

	<-quit
	log.Println("shutdown signal received")

	// ─────────────────────────────────────────────────────
	// STEP 5: Graceful shutdown with a timeout
	// ─────────────────────────────────────────────────────
	//
	// We've received the signal. Now we want to:
	//   1. Stop accepting NEW requests
	//   2. Wait for IN-FLIGHT requests to finish
	//   3. But not wait forever — give them up to 30 seconds
	//
	// Why a timeout? Imagine a request is stuck (slow database, etc.).
	// Without a timeout, the server hangs forever waiting for it.
	// With a 30-second timeout, we force the shutdown even if some
	// requests are still in progress.

	// Create a context that automatically cancels after 30 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	// 'defer' means "run this when main() exits"
	// We MUST call cancel() to release the resources associated with the
	// context, even if shutdown succeeds before the timeout.
	defer cancel()

	// server.Shutdown does the actual graceful shutdown:
	//   - Stops accepting new connections immediately
	//   - Waits for active requests to finish
	//   - If 'ctx' expires (30 seconds), it gives up and returns an error
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("shutdown error: %v", err)
	}

	log.Println("server stopped")

	// After this line, main() returns and the program exits normally.
}
