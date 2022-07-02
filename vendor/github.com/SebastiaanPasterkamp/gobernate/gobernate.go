// Package gobernate provides an easy HTTP Handler containing all end-points
// required to run a golang service in Kubernetes.
// This code is roughly based on:
// https://blog.gopheracademy.com/advent-2017/kubernetes-ready-service/
package gobernate

import (
	"github.com/SebastiaanPasterkamp/gobernate/handlers"
	"github.com/SebastiaanPasterkamp/gobernate/version"

	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// Gobernate provides an easy HTTP Handler containing all end-points required to
// run a golang service in Kubernetes. This code is roughly based on:
// https://blog.gopheracademy.com/advent-2017/kubernetes-ready-service/
type Gobernate struct {
	// Router is the mux router, so extra end-points can be added.
	Router   *mux.Router
	port     string
	info     version.Info
	listener net.Listener
	srv      *http.Server
	shutdown chan bool
	isReady  *atomic.Value
}

// New creates a Gobernate instace with port and version info
func New(port, name, release, commit, buildTime string) *Gobernate {
	info := version.Info{
		Name:      name,
		Release:   release,
		Commit:    commit,
		BuildTime: buildTime,
	}
	isReady := &atomic.Value{}
	isReady.Store(false)
	shutdown := make(chan bool)
	router := handlers.Router(info, isReady, shutdown)
	srv := &http.Server{
		Handler: router,
	}
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
	}
	return &Gobernate{
		port:     port,
		info:     info,
		listener: listener,
		srv:      srv,
		Router:   router,
		shutdown: shutdown,
		isReady:  isReady,
	}
}

// Launch runs the gobernate service on the port, and serves the version info.
// Returns a shutdown channel to either block while serving, or to close to take
// down the service
func (g *Gobernate) Launch() chan bool {
	log.Printf("Starting %s, Version %s, Commit %s, Build Time %s",
		g.info.Name, g.info.Release, g.info.Commit, g.info.BuildTime)

	go g.shutdownOnSignals()
	go g.serve()

	return g.shutdown
}

// Ready signals that the service is ready to serve. Call once all initialization
// has completed.
func (g *Gobernate) Ready() {
	g.isReady.Store(true)
}

// URL returns the address:port on which the service is listening.
func (g *Gobernate) URL() string {
	return "http://" + g.listener.Addr().String()
}

func (g *Gobernate) serve() {
	log.Print("The service is ready to listen and serve.")
	if err := g.srv.Serve(g.listener); err != http.ErrServerClosed {
		// Error starting or closing listener:
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}
}

func (g *Gobernate) shutdownOnSignals() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	select {
	case sig := <-sigs:
		// We received an interrupt/terminate signal; shut down.
		log.Errorf("Received signal: %+v. Shutting down...", sig)
		g.Shutdown()
		break
	case <-g.shutdown:
		log.Error("Received shutdown command. Shutting down...")
		break
	}
}

// Shutdown gracefully terminates the service. Automatically called when
// receiving SIGINT, or SIGTERM.
func (g *Gobernate) Shutdown() {
	if err := g.srv.Shutdown(context.Background()); err != nil {
		// Error from closing listeners, or context timeout:
		log.Printf("HTTP server Shutdown: %v", err)
	}

	select {
	case <-g.shutdown:
		// already closed
		break
	default:
		close(g.shutdown)
	}

}
