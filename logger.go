package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
)

func loggingMiddleware(next http.Handler) http.Handler {
	t := time.Now().Format("02-01-2006")
	file, errFile := os.OpenFile(t+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if errFile != nil {
		log.Println("Error opening .log file")
	}
	return handlers.LoggingHandler(io.MultiWriter(file, os.Stdout), next)
}

/*
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// before request
		t := time.Now()

		next.ServeHTTP(w, r)

		// after request
		latency := time.Since(t)
		if !strings.HasPrefix(r.URL.String(), "/documentation/") {
			log.Println(r.URL.String() + "\t" + r.RemoteAddr + "\t" + latency.String())
		}
	})
}
*/
