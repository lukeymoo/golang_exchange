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
	session "../session"
	database "../database"
	"os"
	"runtime"
)

func main() {

	// Server using
	fmt.Printf("CPU Usage %d\n", runtime.NumCPU());

	// Initialize redis client
	fmt.Println("Initializing session handler...")
	session.InitRedis()

	// Initialize mongo client
	fmt.Println("Initializing database connection...")
	database.InitMgo()

	// Initialize router
	fmt.Println("Initializing multiplexor...")
	router := httprouter.New()

	// Port
	port := ":3000"

	// Static files root path
	static := os.Getenv("GO_STATIC_PATH")

	// Static files
	router.ServeFiles("/js/*filepath", http.Dir(static + "/js")) // Javascript
	router.ServeFiles("/css/*filepath", http.Dir(static + "/css")) // CSS
	router.ServeFiles("/img/*filepath", http.Dir(static + "/img")) // Img
	router.ServeFiles("/font/*filepath", http.Dir(static + "/font")) // Font
	router.ServeFiles("/h/*filepath", http.Dir(static + "/h")) // Static html files

	// Routes
	router.GET("/", routes.IndexPage) // Index ( Home ) page
	router.GET("/register", routes.RegisterPage) // Display registration page
	router.POST("/register/process", routes.ProcessRegister) // Process registration form
	router.POST("/login/process", routes.ProcessLogin) // Process login form ( Ajax Request )
	router.GET("/forgot", routes.ForgotPage) // Present user with form to reset password
	router.GET("/account", routes.AccountPage) // Account settings page
	router.GET("/login/logout", routes.Logout) // Logout

	/** API Endpoints **/
	router.GET("/api/session/state", routes.SessionState)

	// Debug
	router.GET("/debug", routes.Debug)

	// Start listening
	fmt.Println("Server listening on port " + port)
	log.Fatal(http.ListenAndServe(port, router));

	return
}