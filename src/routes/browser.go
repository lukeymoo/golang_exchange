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
	"sync"
)

// DELETE LATER
func Debug(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	var registerForm = &form.RegisterForm {
		Firstname: "Lukeymoo",
		Lastname: req.FormValue("l"),
		Username: "FakeUsername!",
	}
	fmt.Fprint(res, "Errors => " + registerForm.CreateQuery())
	return
}

/**
	GET
	/
	Display home page
*/
func IndexPage(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	var context session.TemplateContext
	
	context.TITLE = "Home"
	context.PAGE = "HOME"
	
	var wg sync.WaitGroup

	wg.Add(1)
	
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		session.CreateSessionObj(&context);
	}(&wg)

	tmp, err := template.ParseFiles("template/layout/master.html", "template/index.html")
	if err != nil {
		fmt.Println("Error =>", err)
		fmt.Fprint(res, "Server error occurred")
		return
	}

	// Wait
	wg.Wait()

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

	// Create form
	var registerForm = &form.RegisterForm {
		Firstname: 		req.FormValue("f"),
		Lastname: 		req.FormValue("l"),
		Username: 		req.FormValue("u"),
		Password:		req.FormValue("p"),
		PasswordAgain: 	req.FormValue("pa"),
		Email:			strings.ToLower(req.FormValue("e")),
		EmailAgain:		strings.ToLower(req.FormValue("ea")),
		Zipcode:		req.FormValue("z"),
		Gender:			req.FormValue("g"),
		Tos:			req.FormValue("tos"),
	}

	err := ""
	query := registerForm.CreateQuery()

	// Validate the form
	err = registerForm.ValidateForm()
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

