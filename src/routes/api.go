package routes

/**

Ajax based routes are contained here

*/

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	//models "../models"
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

	// Validate form
	var loginForm = &form.LoginForm {
		UsernameOrEmail: 	req.FormValue("u"),
		Password: 			req.FormValue("p"),
	}

	err := loginForm.ValidateForm()

	// Send errors to user
	if err != "" {
		jsonResponse, err := json.Marshal(struct{Error string}{Error: err})
		if err != nil {
			fmt.Fprint(res, "Server Error")
			return
		}
		res.Header().Set("Content-Type", "application/json")
		res.Write(jsonResponse)
		fmt.Println("sent")
		return
	}

	fmt.Fprint(res, "done")
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