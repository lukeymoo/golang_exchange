package routes
/**
	Functions for routes
*/

import (
	"fmt"
	"net/http"
	"github.com/julienschmidt/httprouter"
	redis "github.com/alphazero/Go-Redis"
	session "../session"
	"html/template"
	"regexp"
	"log"
)

var (
	rc redis.Client
)

/**
	Initializes redis session client
*/
func InitRedis() {
	spec := redis.DefaultSpec().Db(0).Password("9b3af6edcf71b34520a7d16412ad9325OMGOMG")
	client, err := redis.NewSynchClientWithSpec(spec)
	if err != nil {
		log.Fatal(err)
	}
	rc = client
	return
}

/**
	GET
	/
	Display home page
*/
func IndexPage(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	tmp, err := template.ParseFiles("template/layout/master.html", "template/index.html")
	if err != nil {
		fmt.Println("Error =>", err)
		fmt.Fprint(res, "Server error occurred")
		return
	}

	var context session.TemplateContext

	context.TITLE = "Home"
	context.PAGE = "HOME"

	session.CreateSessionObj(&context);


	err = tmp.Execute(res, context)
	if err != nil {
		fmt.Println("Error =>", err)
		return
	}
	return
}

/**
	GET
	/register
	Display register page
*/
func RegisterPage(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	tmp, err := template.ParseFiles("template/layout/master.html", "template/index.html")
	if err != nil {
		fmt.Println("Error =>", err)
		fmt.Fprint(res, "Server error occurred")
		return
	}
	var context session.TemplateContext

	context.TITLE = "Register"
	context.PAGE = "REGISTER"

	session.CreateSessionObj(&context);

	err = tmp.Execute(res, context)
	if err != nil {
		fmt.Println("Server error occurred")
		return
	}
	return
}

/**
	POST
	/register/process
	Processes register form & creates an account
*/
func ProcessRegister(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// Ensure not logged in
	if session.IsSession() {
		http.Redirect(res, req, "/register?err=logged_in", 301)
		return
	}
	// Ensure we received all register form elements
	formError := req.ParseForm()
	if formError != nil {
		fmt.Println("Error =>", formError)
		fmt.Fprint(res, "Server error")
		return
	}

	// Ensure we recieved all form fields
	if !validRegisterForm(req) {
		http.Redirect(res, req, "/register?err=invalid_form", 301)
		return
	}

	err := ""
	query := "&f=" + req.Form.Get("f") + "&l=" + req.Form.Get("l") +
			"&u=" + req.Form.Get("u") + "&e=" + req.Form.Get("e") + "&z=" + req.Form.Get("z") +
			"&g=" + req.Form.Get("g")
	
	// Validate names
	if !validName(req.Form.Get("f")) {
		err += "F|"
	}
	if !validName(req.Form.Get("l")) {
		err += "L|"
	}
	// Validate Username
	if !validUsername(req.Form.Get("u")) {
		err += "U|"
	}
	// Password
	if !validPassword(req.Form.Get("p")) {
		err += "P|"
	}
	// Password again
	if req.Form.Get("p") != req.Form.Get("pa") {
		err += "PM|"
	}
	// Email
	if !validEmail(req.Form.Get("e")) {
		err += "E|"
	}
	// Email again
	if req.Form.Get("e") != req.Form.Get("ea") {
		err += "EM|"
	}
	// Zipcode
	if !validZipcode(req.Form.Get("z")) {
		err += "Z|"
	}
	// If errors redirect with error codes
	if err != "" {
		err = err[:len(err)-1]
		http.Redirect(res, req, "/register?err=" + err + query, 301)
		return
	}

	/** Determine if the username and/or email is in use **/
	if session.IsSession() {
		if session.IsTimedOut() {
			fmt.Fprint(res, "Session has timed out")
		} else {
			fmt.Fprint(res, "Logged in")
		}
	} else {
		fmt.Fprint(res, "Not logged in")
	}
	return
}



/**
	GET
	/login/logout
	Logs the user out
*/
func Logout(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	session.ClearSession();
	http.Redirect(res, req, "/", 301)
	return
}





/** Helpers **/


func validRegisterForm(req *http.Request) bool {

	/** Ensure fields are set and not empty **/
	// First name
	if req.Form.Get("f") == "" {
		return false
	}
	// Last name
	if req.Form.Get("l") == "" {
		return false
	}
	// Username
	if req.Form.Get("u") == "" {
		return false
	}
	// Password
	if req.Form.Get("p") == "" {
		return false
	}
	// Password again
	if req.Form.Get("pa") == "" {
		return false
	}
	// Email
	if req.Form.Get("e") == "" {
		return false
	}
	// Email again
	if req.Form.Get("ea") == "" {
		return false
	}
	// Zipcode
	if req.Form.Get("z") == "" {
		return false
	}
	// Gender must be set
	if req.Form.Get("g") == "" {
		return false
	}
	// Gender must be 'm' or 'f'
	if req.Form.Get("g") != "m" && req.Form.Get("g") != "f" {
		return false
	}
	// Terms and conditions must be I_AGREE
	if req.Form.Get("tos") != "I_AGREE" {
		return false
	}

	return true
}

// Validate Name
func validName(name string) bool {
	// Regex
	match, err := regexp.MatchString(`^[A-Za-z]+(([\'-])?[A-Za-z]+$)`, name)
	if err != nil {
		log.Fatal(err)
	}
	if !match {
		return false
	}

	// Test length
	if len(name) < 2 || len(name) > 32 {
		return false
	}


	return true
}

// Validate Username
func validUsername(name string) bool {
	// Regex
	match, err := regexp.MatchString(`^[A-Za-z0-9_]+$`, name)
	if err != nil {
		log.Fatal(err)
	}
	if !match {
		return false
	}
	// Test length
	if len(name) < 2 || len(name) > 16 {
		return false
	}
	return true
}

// Validate Password
func validPassword(pwd string) bool {
	if len(pwd) < 2 || len(pwd) > 32 {
		return false
	}
	return true
}

// Validate email
func validEmail(email string) bool {
	match, err := regexp.MatchString(`^([a-zA-Z0-9_.+-])+\@(([a-zA-Z0-9-])+\.)+([a-zA-Z0-9]{2,4})+$`, email)
	if err != nil {
		log.Fatal(err)
	}
	if !match {
		return false
	}
	if len(email) > 64 {
		return false
	}
	return true
}

// Validate zipcode
func validZipcode(code string) bool {
	match, err := regexp.MatchString(`^[0-9]+$`, code)
	if err != nil {
		log.Fatal(err)
	}
	if !match {
		return false
	}
	if len(code) != 5 {
		return false
	}
	return true
}