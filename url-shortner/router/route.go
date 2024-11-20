package router

import (
	"net/http"
	"urlshortner/controller"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/", controller.ConvertUrl).Methods("POST")
	router.HandleFunc("/{shortUrl}", controller.GetLongUrl).Methods("GET")

	return router
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}
