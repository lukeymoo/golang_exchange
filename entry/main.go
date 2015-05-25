package main
/**
	Entry point for server
*/

import (
	"fmt"
	"net/http"
	"log"
	"github.com/julienschmidt/httprouter"
	controller "../controllers"
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

	// controller
	router.GET("/", controller.IndexPage) // Index ( Home ) page
	router.GET("/register", controller.RegisterPage) // Display registration page
	router.POST("/register/process", controller.ProcessRegister) // Process registration form
	router.POST("/login/process", controller.ProcessLogin) // Process login form ( Ajax Request )
	/** INCOMPLETE **/
	router.GET("/forgot", controller.ForgotPage) // Present user with form to reset password
	/** INCOMPLETE **/
	router.GET("/account", controller.AccountPage) // Account settings page
	router.GET("/p/new", controller.CreatePost) // Create Post page
	/** INCOMPLETE **/
	router.POST("/p/process", controller.ProcessPost) // Process Post Form
	router.GET("/login/logout", controller.Logout) // Logout

	/** API Endpoints **/
	router.GET("/api/session/state", controller.SessionState)

	// Debug
	router.GET("/debug", controller.Debug)

	// Start listening
	fmt.Println("Server listening on port " + port)
	log.Fatal(http.ListenAndServe(port, router));

	return
}