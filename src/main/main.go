package main
/**
	Entry point for server
*/

import (
	"fmt"
	"net/http"
	"log"
	"github.com/julienschmidt/httprouter"
	routes "../routes"
)

func main() {
	// Initialize server
	router := httprouter.New()

	// Port
	port := ":3000"

	// Static files root path
	static := "/home/busx/Documents/goweb/public/"

	// Static files
	router.ServeFiles("/js/*filepath", http.Dir(static + "/js")) // Javascript
	router.ServeFiles("/css/*filepath", http.Dir(static + "/css")) // CSS
	router.ServeFiles("/img/*filepath", http.Dir(static + "/img")) // Img
	router.ServeFiles("/font/*filepath", http.Dir(static + "/font")) // Font

	// Routes
	router.GET("/", routes.Index)

	// Start listening
	fmt.Println("Server listening on port " + port)
	log.Fatal(http.ListenAndServe(port, router))

	return
}