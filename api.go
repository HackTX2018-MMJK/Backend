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
	router.HandleFunc("/rr/enqueue", EnqueueHandler).Methods("POST").
		Headers("Content-Type", "application/json")
	router.HandleFunc("/rr/status/{uuid}", StatusHandler).Methods("GET")
	router.HandleFunc("/rr/size", SizeHandler).Methods("GET")
	router.HandleFunc("/rr/dequeue", DequeueHandler).Methods("POST").
		Headers("Content-Type", "application/json")

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

func DequeueHandler(w http.ResponseWriter, r *http.Request) {
	/* Extract values */
	var body StringResponse

	/* Parse POSTED data */
	if r.Body == nil {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity),
			http.StatusUnprocessableEntity)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity),
			http.StatusUnprocessableEntity)
		return
	}

	token := body.Response
	_, exists := tokens[token]

	if exists {
		Remove_waiter(token)
	} else {
		http.NotFound(w, r)
	}
}

func SizeHandler(w http.ResponseWriter, r *http.Request) {
	response := IntResponse{Response: Get_length()}

	b, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type:", "application/json")
	w.Write(b)
}

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	/* Extract values */
	vars := mux.Vars(r)

	token := vars["uuid"]
	_, exists := tokens[vars["uuid"]]

	if exists {
		response := IntResponse{Response: Get_position(token)}

		b, err := json.Marshal(response)
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

func EnqueueHandler(w http.ResponseWriter, r *http.Request) {
	/* Extract values */
	var body StringResponse

	/* Parse POSTED data */
	if r.Body == nil {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity),
			http.StatusUnprocessableEntity)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity),
			http.StatusUnprocessableEntity)
		return
	}

	token := body.Response
	_, exists := tokens[token]

	if exists {
		Add_waiter(token)
	} else {
		http.NotFound(w, r)
	}
}

func SeatHandler(w http.ResponseWriter, r *http.Request) {
	/* Extract values */
	vars := mux.Vars(r)
	token, exists := tokens[vars["uuid"]]

	if exists {
		response := StringResponse{Response: token}

		b, err := json.Marshal(response)
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

type IntResponse struct {
	Response int `json "response"`
}

type StringResponse struct {
	Response string `json "response"`
}
