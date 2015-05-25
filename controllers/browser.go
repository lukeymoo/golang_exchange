package controllers

/**

Browser based routes are contained in this file

*/

import (
	models "../models"
	session "../session"
	"../util"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"net/http"
	"strings"
)

// DELETE LATER
func Debug(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	fmt.Fprint(res, "Nothing to debug\nBut thanks for the visit!")
	return
}

/**
GET
/
Display home page
*/
func IndexPage(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	var context session.TemplateContext

	session.CreateSessionObj(&context)

	context.TITLE = "Home"
	context.PAGE = "HOME"

	tmp, err := template.ParseFiles("templates/layout/master.html", "templates/index.html")
	if err != nil {
		fmt.Println("Error =>", err)
		fmt.Fprint(res, "Server error occurred")
		return
	}

	err = tmp.Execute(res, context)
	if err != nil {
		fmt.Println("Error =>", err)
	}
	return
}

/**
GET
/register
Display register page
*/
func RegisterPage(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	tmp, err := template.ParseFiles("templates/layout/master.html", "templates/index.html")
	if err != nil {
		fmt.Println("Error =>", err)
		fmt.Fprint(res, "Server error occurred")
		return
	}
	var context session.TemplateContext

	context.TITLE = "Register"
	context.PAGE = "REGISTER"

	session.CreateSessionObj(&context)

	err = tmp.Execute(res, context)
	if err != nil {
		fmt.Println("Server error occurred")
	}
	return
}

/**
Cannot be authenticated
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

	// Create form
	var registerForm = &util.RegisterForm{
		Firstname:     strings.ToLower(req.FormValue("f")),
		Lastname:      strings.ToLower(req.FormValue("l")),
		Username:      strings.ToLower(req.FormValue("u")),
		Password:      req.FormValue("p"),
		PasswordAgain: req.FormValue("pa"),
		Email:         strings.ToLower(req.FormValue("e")),
		EmailAgain:    strings.ToLower(req.FormValue("ea")),
		Zipcode:       req.FormValue("z"),
		Gender:        strings.ToLower(req.FormValue("g")),
		Tos:           req.FormValue("tos"),
	}

	err := ""
	query := registerForm.CreateQuery()

	// Validate the form
	err = registerForm.ValidateForm()
	if err != "" {
		err = err[:len(err)-1]
		http.Redirect(res, req, "/register?err="+err+query, 302)
		return
	}

	// Save the user, set session & redirect
	var user models.User
	user.Firstname = registerForm.Firstname
	user.Lastname = registerForm.Lastname
	user.Username = registerForm.Username
	user.Password = registerForm.Password
	user.Email = registerForm.Email
	user.Zipcode = registerForm.Zipcode
	if models.SaveUser(user) {
		// Get Saved user
		user = models.FindUserByUsername(registerForm.Username)
		// Set session
		session.SetSession(user.Username, user.Email, user.Id.String())
		// Redirect
		http.Redirect(res, req, "/", 302)
		return
	} else {
		// Server error
		fmt.Fprint(res, "Server error => ", err)
	}
	return
}

/**
Needs auth
GET
/p/new
Create new post page
*/
func CreatePost(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	if !session.IsSession() {
		http.Redirect(res, req, "/?err=need_login&next=/p/new", 302)
		return
	}

	tmp, err := template.ParseFiles("templates/layout/master.html", "templates/index.html")
	if err != nil {
		fmt.Println("Error =>", err)
		fmt.Fprint(res, "Server error occurred")
		return
	}
	var context session.TemplateContext

	context.TITLE = "Create Post"
	context.PAGE = "CREATEPOST"

	session.CreateSessionObj(&context)

	err = tmp.Execute(res, context)
	if err != nil {
		fmt.Println("Error create post template :: ", err)
	}
	return
}

/**
needs auth
POST
/p/process
Process create post form
*/
func ProcessPost(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	if !session.IsSession() {
		http.Redirect(res, req, "/?err=need_login&next=/p/new", 302)
		return
	}
	/** Save files **/
	fmt.Fprint(res, "This route is incomplete")
	return
}

func ForgotPage(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	fmt.Fprint(res, "Does nothing yet")
	return
}

/*
Needs auth
GET
/account
Account settings page
*/
func AccountPage(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	if !session.IsSession() {
		http.Redirect(res, req, "/?err=need_login&next=/account", 302)
		return
	}
	fmt.Fprint(res, "Account page")
	return
}

/*
GET
/login/logout
Clear session variables
*/
func Logout(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	session.ClearSession()
	http.Redirect(res, req, "/", 302)
	return
}
