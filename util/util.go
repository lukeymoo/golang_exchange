package util

import (
	"encoding/hex"
	"crypto/md5"
	"encoding/json"
	"net/http"
	"fmt"
	db "../database"
	"gopkg.in/mgo.v2/bson"
	"regexp"
	"log"
	models "../models"
	"sync"
)

type (

	ApiResponse struct {
		Status string
		Message string
	}

	ApiErrorResponse struct {
		Error string
	}

	SessionStateResponse struct {
		LOGGED_IN bool
		USERNAME string
	}

	PostForm struct {
		Type 		string
		Description string
		Images 		[]string
	}

	RegisterForm struct {
		Firstname 		string
		Lastname 		string
		Username 		string
		Password 		string
		PasswordAgain 	string
		Email 			string
		EmailAgain 		string
		Zipcode 		string
		Gender 			string
		Tos 			string
	}

	LoginForm struct {
		UsernameOrEmail 	string
		Password 			string
	}

)

/**
	Encodes a given string into an MD5 sum and returns it in string format
*/
func MD5String(text string) (string) {
	hashed := md5.New()
	hashed.Write([]byte(text))
	return hex.EncodeToString(hashed.Sum(nil))
}

func ApiSend(res http.ResponseWriter, status string, message string) {
	var temp ApiResponse = ApiResponse{status, message}
	temp_formatted, err := json.Marshal(temp)
	if err != nil {
		fmt.Fprint(res, "Server error")
		fmt.Println("Error creating JSON response => ", err)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.Write(temp_formatted)
	return
}

func ApiError(res http.ResponseWriter, error string) {
	var temp ApiErrorResponse
	temp.Error = error
	temp_formatted, err := json.Marshal(temp)
	if err != nil {
		fmt.Fprint(res, "Server error")
		fmt.Println("Error creating JSON Error Response => ", err)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	res.Write(temp_formatted)
	return
}



/*** Register/Login Forms ***/




/** Loads form data into struct and validates it **/
func (form *RegisterForm) ValidateForm() (string) {
	err := ""
	/** If the form is empty return invalid_form **/
	if form.IsEmpty() {
		err = "invalid_form|"
		return err
	}
	// Validate names
	if !ValidName(form.Firstname) {
		err += "F|"
	}
	if !ValidName(form.Lastname) {
		err += "L|"
	}
	// Validate Username
	if !ValidUsername(form.Username) {
		err += "U|"
	}
	// Password
	if !ValidPassword(form.Password) {
		err += "P|"
	}
	// Password again
	if form.Password != form.PasswordAgain {
		err += "PM|"
	}
	// Email
	if !ValidEmail(form.Email) {
		err += "E|"
	}
	// Email again
	if form.Email != form.EmailAgain {
		err += "EM|"
	}
	// Zipcode
	if !ValidZipcode(form.Zipcode) {
		err += "Z|"
	}
	if err == "" {
		// Determine if username or email is already registered
		if models.DoesUsernameExist(form.Username) {
			err += "UIN|"
		}
		if models.DoesEmailExist(form.Email) {
			err += "EIN|"
		}
	}
	return err
}

func (form *RegisterForm) CreateQuery() (string) {
	// Construct a query
	query := "&f=" + form.Firstname + "&l=" + form.Lastname + "&u=" + form.Username + "&e=" + form.Email + "&z=" + form.Zipcode + "&g=" + form.Gender
	return query
}

func (form *RegisterForm) IsEmpty() (bool) {
	if form.Firstname == "" {
		return true
	}
	if form.Lastname == "" {
		return true
	}
	if form.Username == "" {
		return true
	}
	if form.Password == "" {
		return true
	}
	if form.PasswordAgain == "" {
		return true
	}
	if form.Email == "" {
		return true
	}
	if form.EmailAgain == "" {
		return true
	}
	if form.Zipcode == "" {
		return true
	}
	if form.Gender != "m" && form.Gender != "f" {
		return true
	}
	if form.Tos != "I_AGREE" {
		return true
	}
	return false
}














/** API Response **/
func (form *LoginForm) ValidateForm() (string, string) {
	err := ""
	idType := ""

	// Ensure the form is not empty
	if !form.IsComplete() {
		return "", "invalid_form"
	}

	var wg sync.WaitGroup

	wg.Add(2)

	go func(idType *string, err *string) {
		defer wg.Done()
		// Ensure the user supplied a username or email
		*idType = form.GetIDType()

		// handle bad id type
		if *idType == "" {
			*err += "U|"
		}
	}(&idType, &err)

	go func(err *string) {
		defer wg.Done()
		// Ensure the user supplied a valid password
		if !ValidPassword(form.Password) {
			*err += "P|"
		}
	}(&err)

	wg.Wait()

	if err != "" {
		err = err[:len(err)-1]
	}

	return idType, err
}

func (form *LoginForm) GetIDType() (string) {
	idType := ""
	if ValidUsername(form.UsernameOrEmail) {
		idType = "username"
	} else {
		if ValidEmail(form.UsernameOrEmail) {
			idType = "email"
		}
	}
	return idType
}

/** Returns true/false if login form is missing any fields **/
func (form *LoginForm) IsComplete() (bool) {
	if form.UsernameOrEmail == "" {
		return false
	}
	if form.Password == "" {
		return false
	}
	return true
}

func (form *LoginForm) TryLogin() (string) {
	idType := form.GetIDType()

	if idType == "username" {
		if !form.ValidUsernameLogin() {
			return "U_invalid_login|"
		}
	} else if idType == "email" {
		if !form.ValidEmailLogin() {
			return "E_invalid_login|"
		}
	}
	return ""
}

/**
	Determines if a given email + password combination exists
	Returns boolean value
*/
func (form *LoginForm) ValidEmailLogin() (bool) {
	count, err := db.Conn.DB("dmvexchange").C("USERS").Find(bson.M {
		"email": form.UsernameOrEmail,
		"pwd": MD5String(form.Password),
	}).Count()
	if err != nil {
		fmt.Println("[-] MongoDB error => ", err)
		return false
	}
	if count > 0 {
		return true
	}
	return false
}

/**
	Determines if a given username + password combination exists
	Returns boolean value
*/
func (form *LoginForm) ValidUsernameLogin() (bool) {
	count, err := db.Conn.DB("dmvexchange").C("USERS").Find(bson.M {
		"username": form.UsernameOrEmail,
		"pwd": MD5String(form.Password),
	}).Count()
	if err != nil {
		fmt.Println("[-] MongoDB error => ", err)
		return false
	}
	if count > 0 {
		return true
	}
	return false
}









// Validate a post description
func ValidDescription(desc string) (bool) {
	if len(desc) < 2 || len(desc) > 2500 {
		return false
	}
	return true
}

// Validate Name
func ValidName(name string) (bool) {
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
func ValidUsername(name string) bool {
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
func ValidPassword(pwd string) bool {
	if len(pwd) < 2 || len(pwd) > 32 {
		return false
	}
	return true
}

// Validate email
func ValidEmail(email string) bool {
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
func ValidZipcode(code string) bool {
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