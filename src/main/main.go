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
)

func main() {
	// Initialize server
	fmt.Println("Initializing server...")
	router := httprouter.New()

	// Initialize redis client
	fmt.Println("Initializing redis...")
	session.InitRedis()

	// Initialize mongo client
	fmt.Println("Initialize mongo...")
	database.InitMgo()

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
	http.ListenAndServe(port, Log(router))

	return
}

func Log(handler http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
	handler.ServeHTTP(w, r)
    })
}