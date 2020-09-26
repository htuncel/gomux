package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"

	"main/controllers"
	_ "main/docs"
)

var (
	hub *Hub
)

// @title Swagger Example API
// @version 1.0
// @description This is a sample server Example server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @BasePath /api/v1
func main() {
	hub = newHub()
	go hub.run()

	router := mux.NewRouter()
	router.PathPrefix("/documentation/").Handler(httpSwagger.WrapHandler)

	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	fs := http.FileServer(http.Dir("./public/views/"))
	router.PathPrefix("/chat/").Handler(http.StripPrefix("/chat/", fs))

	setupRoutes(router)

	srv := &http.Server{
		Handler: handlers.CORS(
			handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
			handlers.AllowedMethods([]string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"}),
			handlers.AllowedOrigins([]string{"*"}))(router),
		Addr: "127.0.0.1:8080",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func setupRoutes(router *mux.Router) {
	bookController := controllers.NewBookController()
	exampleController := controllers.NewExampleController()

	v1 := router.PathPrefix("/api/v1").Subrouter()

	v1.Use(loggingMiddleware)

	example := v1.PathPrefix("/example").Subrouter()
	example.HandleFunc("/token", exampleController.TokenHandler).Methods("GET")
	example.HandleFunc("/multifile", exampleController.MultipleFileUpload).Methods("POST")

	auth := v1.PathPrefix("/auth").Subrouter()
	auth.Use(Authenticate)
	auth.HandleFunc("", homeHandler).Methods("GET")

	books := v1.PathPrefix("/books").Subrouter()
	books.HandleFunc("", bookController.FindBooks).Methods("GET")
	books.HandleFunc("/{id}", bookController.FindBook).Methods("GET")
	books.HandleFunc("", bookController.CreateBook).Methods("POST")
	books.HandleFunc("/{id}", bookController.UpdateBook).Methods("PATCH")
	books.HandleFunc("/{id}", bookController.DeleteBook).Methods("DELETE")
}

// homeHandler godoc
// @Summary Get welcome message, get token before trying this route
// @Description Get welcome message
// @Tags home
// @Accept json
// @Produce json
// @Header 200 {string} Token "qwerty"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /auth [get]
// @Security ApiKeyAuth
// @param headkey header string true "headkey"
func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	trueClients := 0
	falseClients := 0
	for _, client := range hub.clients {
		if client {
			trueClients++
		} else {
			falseClients++
		}
	}

	socketMessage := []byte("hello")
	hub.broadcast <- socketMessage

	json.NewEncoder(w).Encode(map[string]string{
		"data":   "Hello World",
		"header": r.Header.Get("headkey"),
		"true":   strconv.Itoa(trueClients),
		"false":  strconv.Itoa(falseClients),
	})
}
