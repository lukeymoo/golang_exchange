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
);

func main() {
	// Initialize server
	fmt.Println("Initializing server...");
	router := httprouter.New();

	// Initialize redis client
	fmt.Println("Initializing redis...");
	routes.InitRedis();

	// Port
	port := ":3000";

	// Static files root path
	static := "/home/busx/Documents/goweb/public/";

	// Static files
	router.ServeFiles("/js/*filepath", http.Dir(static + "/js")); // Javascript
	router.ServeFiles("/css/*filepath", http.Dir(static + "/css")); // CSS
	router.ServeFiles("/img/*filepath", http.Dir(static + "/img")); // Img
	router.ServeFiles("/font/*filepath", http.Dir(static + "/font")); // Font

	// Routes
	router.GET("/", routes.IndexPage); // Index ( Home ) page
	router.GET("/register", routes.RegisterPage); // Display registration page
	router.POST("/register/process", routes.ProcessRegister); // Process registration form
	router.GET("/login/logout", routes.Logout); // Logout

	// Debug
	router.GET("/fakelogin", routes.FakeLogin);
	router.GET("/debug", routes.Debug);

	// Start listening
	fmt.Println("Server listening on port " + port);
	log.Fatal(http.ListenAndServe(port, router));

	return;
}