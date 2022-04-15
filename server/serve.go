package server

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	cron "github.com/robfig/cron/v3"
)

// Serve starts a server listening on "listenAddr". In parrallel, a cron-like job
// is also started, pulling the latest registry changes from the provided registry-url
// This function is blocking and can be stopped by cancelling the provided context.
func Serve(ctx context.Context, registryUrl, listenAddr, updateFreq string) error {
	l := log.Default()
	// Set up the handler and pull in all data
	handler := NewHandler(registryUrl, l)
	if err := handler.Pull(ctx); err != nil {
		return err
	}

	// create a router to handle inbound requests
	router := mux.NewRouter()
	router.HandleFunc("/", Ok).Methods("GET")
	// use some form of versioning to allow for future changes
	v1Router := router.PathPrefix("/v1").Subrouter()
	v1Router.HandleFunc("/chains", handler.Chains).Methods("GET")
	v1Router.HandleFunc("/chain/{chain}", handler.Chain).Methods("GET")
	v1Router.HandleFunc("/chain/{chain}/endpoints/{type}", handler.Endpoints).Methods("GET")
	v1Router.HandleFunc("/chain/{chain}/assets", handler.ChainAsset).Methods("GET")
	v1Router.HandleFunc("/assets", handler.Assets).Methods("GET")
	v1Router.HandleFunc("/asset/{asset}", handler.Asset).Methods("GET")
	s := http.Server{Addr: listenAddr, Handler: router}

	errs := make(chan error, 1)
	go func() {
		// If there is an error on startup catch it and pass it through
		// the channel
		errs <- s.ListenAndServe()
	}()

	l.Printf("server up on %s", s.Addr)

	crawler := cron.New(cron.WithLogger(cron.PrintfLogger(l)))
	crawler.AddFunc(updateFreq, func() {
		// update the servers local records
		if err := handler.Pull(ctx); err != nil {
			l.Print(err)
		}
	})
	crawler.Start()
	defer crawler.Stop()

	l.Printf("cron scheduler running with update frequency: %s", updateFreq)

	select {
	// Use contexts to manage the servers lifecycle
	case <-ctx.Done():
		// This will stop the other go routine if it hasn't already
		// stopped yet
		l.Print("shutting down server")
		if err := s.Close(); err != nil {
			return err
		}
		err := <-errs
		return err
	case err := <-errs:
		return err
	}
}

func Ok(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusOK)
}
