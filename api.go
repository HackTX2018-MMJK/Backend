package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

/* HOST to listen on */
var HOST string = "0.0.0.0"

/* PORT to bind to */
var PORT int = 8080

/* mappings represents token->seat mappings */
var tokens map[string]string

func handle_requests(mappings map[string]string) {
	log.Println("Starting Web Server at " + HOST + ":" + strconv.Itoa(PORT))
	tokens = mappings

	router := mux.NewRouter()

	/* Routes */
	router.HandleFunc("/seats/{uuid}", SeatHandler).Methods("GET")

	/* Web Server */
	timeout := time.Second * 15
	server := &http.Server{
		Addr: HOST + ":" + strconv.Itoa(PORT),

		/* Handle Timeouts */
		WriteTimeout: timeout,
		ReadTimeout:  timeout,
		IdleTimeout:  timeout,

		/* Router/Handler */
		Handler: router,
	}

	/* Run Server in non-blocking Goroutine */
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	/* Shutdown Gracefully */
	quit_sig := make(chan os.Signal, 1)
	signal.Notify(quit_sig, os.Interrupt)

	/* Suspend until exit */
	<-quit_sig
	log.Println("Closing Web Server")

	/* Handle remaining connections */
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	/* Close Server */
	server.Shutdown(ctx)
	os.Exit(0)
}

func SeatHandler(w http.ResponseWriter, r *http.Request) {
	/* Extract values */
	vars := mux.Vars(r)
	token, exists := tokens[vars["uuid"]]

	if exists {
		respones := StringResponse{Response: token}

		b, err := json.Marshal(respones)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type:", "application/json")
		w.Write(b)
	} else {
		http.NotFound(w, r)
	}
}

type StringResponse struct {
	Response string `json "response"`
}
