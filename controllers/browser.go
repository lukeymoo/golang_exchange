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
	"sync"
	"os"
	"io"
	"github.com/satori/go.uuid"
	"image"
)

// DELETE LATER
func Debug(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	fmt.Fprintf(res, "LOGGED_IN => %v\nUSERNAME => %s\n", session.IsLoggedIn(), session.GetVariable("USERNAME"))
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
	
	// Validate post type
	if req.FormValue("posttype") != "sale" && req.FormValue("posttype") != "general" {
		fmt.Fprint(res, "Invalid post type")
		return
	}

	// Validate description
	if len(req.FormValue("postdescription")) < 2 || len(req.FormValue("postdescription")) > 2500 {
		fmt.Fprint(res, "Invalid post description")
		return
	}

	// Parse Multipart form
	err := req.ParseMultipartForm(0)
	if err != nil {
		fmt.Fprint(res, "Server error occurred", err)
		return
	}

	formData := req.MultipartForm

	imageCount := 0

	var wg sync.WaitGroup

	// read form data
	photo1 := formData.File["photo1"]
	if photo1 != nil {
		imageCount++
	}
	photo2 := formData.File["photo2"]
	if photo2 != nil {
		imageCount++
	}
	photo3 := formData.File["photo3"]
	if photo3 != nil {
		imageCount++
	}
	photo4 := formData.File["photo4"]
	if photo4 != nil {
		imageCount++
	}

	if imageCount == 0 {
		if req.FormValue("posttype") == "sale" {
			fmt.Fprint(res, "Sale posts need at least 1 image")
			return
		}
	} else if imageCount > 0 {
		wg.Add(imageCount)
		// Parse images
		if photo1 != nil {
			go func(wg *sync.WaitGroup) {
				defer wg.Done()
				fmt.Fprintf(res, "Photo1 => %t\n", photo1)
			}(&wg)
		}
		if photo2 != nil {
			go func(wg *sync.WaitGroup) {
				defer wg.Done()
				ph2, err:= photo2[0].Open()
				if err != nil {
					fmt.Fprintln(res, "Error occurred reading photo2")
				}
				defer ph2.Close()
				tmp, err := os.Create(os.Getenv("GO_STATIC_PATH") + "/cdn/product/" + uuid.NewV4().String())
				if err != nil {
					fmt.Fprintln(res, "Failed to create temp file for photo2")
					return
				}
				defer tmp.Close()
				io.Copy(tmp, ph2)

				// Validate the file

			}(&wg)
		}
		if photo3 != nil {
			go func(wg *sync.WaitGroup) {
				defer wg.Done()
			}(&wg)
		}
		if photo4 != nil {
			go func(wg *sync.WaitGroup) {
				defer wg.Done()
			}(&wg)
		}

		wg.Wait()
	}

	fmt.Fprintf(res, "Post type => %s\n", req.FormValue("posttype"))
	fmt.Fprintf(res, "Post description :: \n\n%s\n\n", req.FormValue("postdescription"))
	fmt.Fprintln(res, "File upload count => ", imageCount)
	fmt.Fprintln(res, "photo1 => ", photo1)
	fmt.Fprintln(res, "photo2 => ", photo2)
	fmt.Fprintln(res, "photo3 => ", photo3)
	fmt.Fprintln(res, "photo4 => ", photo4)
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
