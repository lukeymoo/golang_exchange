package routes

/**

Browser based routes are contained in this file

*/

import (
	"fmt"
	"net/http"
	"github.com/julienschmidt/httprouter"
	session "../session"
	models "../models"
	form "../forms"
	"html/template"
	"strings"
)

// DELETE LATER
func Debug(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	fmt.Fprint(res, "This path does nothing lol!!")
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
		http.Redirect(res, req, "/register?err=logged_in", 302)
		return
	}
	// Ensure we received all register form elements
	formError := req.ParseForm()
	if formError != nil {
		fmt.Println("Error =>", formError)
		fmt.Fprint(res, "Server error please try again")
		return
	}

	// Ensure we recieved all form fields
	if !form.ValidRegisterForm(req) {
		http.Redirect(res, req, "/register?err=invalid_form", 302)
		return
	}

	err := ""
	query := "&f=" + req.Form.Get("f") + "&l=" + req.Form.Get("l") +
			"&u=" + req.Form.Get("u") + "&e=" + req.Form.Get("e") + "&z=" + req.Form.Get("z") +
			"&g=" + req.Form.Get("g")
	
	// Validate names
	if !form.ValidName(req.Form.Get("f")) {
		err += "F|"
	}
	if !form.ValidName(req.Form.Get("l")) {
		err += "L|"
	}
	// Validate Username
	if !form.ValidUsername(req.Form.Get("u")) {
		err += "U|"
	}
	// Password
	if !form.ValidPassword(req.Form.Get("p")) {
		err += "P|"
	}
	// Password again
	if req.Form.Get("p") != req.Form.Get("pa") {
		err += "PM|"
	}
	// Email
	if !form.ValidEmail(req.Form.Get("e")) {
		err += "E|"
	}
	// Email again
	e_lower := strings.ToLower(req.Form.Get("e"))
	ea_lower := strings.ToLower(req.Form.Get("ea"))
	if e_lower != ea_lower {
		err += "EM|"
	}
	// Zipcode
	if !form.ValidZipcode(req.Form.Get("z")) {
		err += "Z|"
	}
	// If errors redirect with error codes
	if err != "" {
		err = err[:len(err)-1]
		http.Redirect(res, req, "/register?err=" + err + query, 302)
		return
	}

	// If username is in use, redirect with error
	err = ""
	if models.DoesUsernameExist(req.Form.Get("u")) {
		err += "UIN|"
	}
	if models.DoesEmailExist(req.Form.Get("e")) {
		err += "EIN|"
	}
	if err != "" {
		err = err[:len(err)-1]
		http.Redirect(res, req, "/register?err=" + err + query, 302)
		return
	}
	// Save the user, set session & redirect
	var user models.User
	user.Firstname 	= req.Form.Get("f")
	user.Lastname 	= req.Form.Get("l")
	user.Username 	= req.Form.Get("u")
	user.Password 	= req.Form.Get("p")
	user.Email 		= req.Form.Get("e")
	user.Zipcode 	= req.Form.Get("z")
	if models.SaveUser(user) {
		// Get Saved user
		user = models.FindUserByUsername(req.Form.Get("u"))
		// Set session
		session.SetSession(req.Form.Get("u"), req.Form.Get("e"), user.Id.String())
		// Redirect
		http.Redirect(res, req, "/", 302)
		return
	} else {
		// Server error
		fmt.Fprint(res, "Server error => ", err)
	}
	return
}


func ForgotPage(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	fmt.Fprint(res, "Does nothing yet")
	return
}

func AccountPage(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// Must be logged in
	if !session.IsSession() {
		http.Redirect(res, req, "/?err=need_login&next=/account", 302)
		return
	}
	fmt.Fprint(res, "Account page")
	return
}


/**
	GET
	/login/logout
	Logs the user out
*/
func Logout(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	session.ClearSession();
	http.Redirect(res, req, "/", 302)
	return
}


/** HELPER **/

