package main

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type User struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var users = map[string]User{
	"1": {
		Id:   "1",
		Name: "AA",
		Age:  10,
	},
	"2": {
		Id:   "2",
		Name: "BB",
		Age:  20,
	},
}

func main() {
	port := os.Getenv("APP_PORT")

	slog.SetLogLoggerLevel(slog.LevelDebug)

	mux := http.DefaultServeMux
	mux.HandleFunc("GET /api/v1/users", func(w http.ResponseWriter, r *http.Request) {
		b, _ := json.Marshal(&users)
		w.WriteHeader(http.StatusOK)
		w.Write(b)
	})
	mux.HandleFunc("GET /api/v1/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("informed user id"))
			return
		}
		user, ok := users[id]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("user not found"))
			return
		}
		b, _ := json.Marshal(&user)
		w.WriteHeader(http.StatusOK)
		w.Write(b)
	})

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	go func() {
		slog.Info("Server is listening", "port", port)
		if err := server.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				slog.Error("begin server error", "error", err)
			}
		}
		close(shutdown)
	}()

	<-shutdown
	slog.Info("Server closed")
}
