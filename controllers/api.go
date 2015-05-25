package controllers

/**

Ajax based routes are contained here

*/

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	models "../models"
	session "../session"
	"../util"
	json "encoding/json"
	"strings"
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

	
	var loginForm = &util.LoginForm {
		UsernameOrEmail: 	strings.ToLower(req.FormValue("u")),
		Password: 			req.FormValue("p"),
	}

	// Validate form
	idType, err := loginForm.ValidateForm()

	// Send errors to user
	if err != "" {
		jsonResponse, err := json.Marshal(struct{Error string}{Error: err})
		if err != nil {
			fmt.Fprint(res, "Server Error")
			return
		}
		res.Header().Set("Content-Type", "application/json")
		res.Write(jsonResponse)
		return
	}

	// try to login
	err = loginForm.TryLogin()
	if err != "" {
		err = err[:len(err)-1]
		jsonResponse, err := json.Marshal(struct{Error string}{Error: err})
		if err != nil {
			fmt.Fprint(res, "Server error")
			return
		}
		res.Header().Set("Content-Type", "application/json")
		res.Write(jsonResponse)
		return
	}

	// Fetch account details
	var userAccount models.User
	if idType == "username" {
		userAccount = models.FindUserByUsername(loginForm.UsernameOrEmail)
	} else if idType == "email" {
		userAccount = models.FindUserByEmail(loginForm.UsernameOrEmail)
	}

	// Good login set session variables and send DX-OK
	jsonResponse, errR := json.Marshal(struct{
			Status string
			Message string
		}{Status: "DX-OK", Message: "Logged in"})
	if errR != nil {
		fmt.Fprint(res, "Server error")
		return
	}

	session.SetSession(userAccount.Username, userAccount.Email, userAccount.Id.String())
	res.Header().Set("Content-Type", "application/json")
	res.Write(jsonResponse)
	return
}

func SessionState(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	var temp util.SessionStateResponse

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