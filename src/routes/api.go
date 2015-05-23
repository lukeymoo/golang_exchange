package routes

/**

Ajax based routes are contained here

*/

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	models "../models"
	session "../session"
	form "../forms"
	json "encoding/json"
	helper "../helper"
)

/**
	API Endpoint
	POST
	/login/process
	Processes login form & sets session
*/
func ProcessLogin(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// Ensure not logged in
	if session.IsSession() {
		http.Redirect(res, req, "/", 302)
		return
	}
	// Parse form
	formError := req.ParseForm()
	if formError != nil {
		fmt.Println("Error authenticating => ", formError)
		fmt.Fprint(res, "Server error please try again")
		return
	}
	// Ensure proper form
	if !form.ValidLoginForm(req) {
		helper.ApiSend(res, "DX-REJECTED", "Invalid form, missing fields")
		return
	}
	u_type := ""
	// Is valid username or email
	if form.ValidUsername(req.Form.Get("u")) {
		u_type = "username"
	} else if form.ValidEmail(req.Form.Get("u")) {
		u_type = "email"
	}

	// Invalid username or email
	if u_type == "" {
		helper.ApiSend(res, "DX-REJECTED", "Invalid username or email")
		return
	}

	// Validate password
	if !form.ValidPassword(req.Form.Get("p")) {
		helper.ApiSend(res, "DX-REJECTED", "Password must be 2-32 characters")
		return
	}

	// If username was supplied
	if u_type == "username" {
			// And if the username + password combo worked, say DX-OK
		if models.ValidUsernameLogin(req.Form.Get("u"), req.Form.Get("p")) {
			var user models.User // placeholder
			user = models.FindUserByUsername(req.Form.Get("u"))
			session.SetSession(req.Form.Get("u"), user.Email, user.Id.String())
		
			helper.ApiSend(res, "DX-OK", "Logged in")
			return
		} else { // If invalid username + password combo
			helper.ApiSend(res, "DX-REJECTED", "Invalid username/password combo")
			return
		}
	} else { // If an email was supplied
		// If the email + password combo worked, say DX-OK
		if models.ValidEmailLogin(req.Form.Get("u"), req.Form.Get("p")) {
			var user models.User // placeholder
			user = models.FindUserByEmail(req.Form.Get("u"))
			session.SetSession(user.Username, req.Form.Get("e"), user.Id.String())
			
			helper.ApiSend(res, "DX-OK", "Logged in")
			return
		} else { // If invalid email + password combo
			helper.ApiSend(res, "DX-REJECTED", "Invalid email/password combo")
			return
		}
	}

	return
}

func SessionState(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	var temp helper.SessionStateResponse

	logged_in := session.GetVariable("LOGGED_IN")

	if logged_in == "true" {
		temp.LOGGED_IN = true
		temp.USERNAME = session.GetVariable("USERNAME")
	} else {
		temp.LOGGED_IN = false
		temp.USERNAME = ""
	}

	temp_formatted, err := json.Marshal(temp)
	if err != nil {
		fmt.Fprint(res, "Server error")
		fmt.Println("Error creating JSON SessionState Response => ", err)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.Write(temp_formatted)
	return
}