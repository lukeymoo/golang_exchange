package forms

import (
	"net/http"
	"regexp"
	"log"
	"../models"
)

type RegisterForm struct {
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

type LoginForm struct {
	UsernameOrEmail 	string
	Password 			string
}

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
func (form *LoginForm) ValidateForm() (string) {
	err := ""
	// determine if the user supplied a username or an email
	idType := ""
	if ValidUsername(form.UsernameOrEmail) {
		idType = "username"
	} else {
		if ValidEmail(form.UsernameOrEmail) {
			idType = "email"
		}
	}

	if idType == "" {
		err += "U|"
	}

	// validate password
	if !ValidPassword(form.Password) {
		err += "P|"
	}

	return err
}

func (form *LoginForm) IsEmpty() (bool) {
	if form.UsernameOrEmail == "" {
		return true
	}
	if form.Password == "" {
		return true
	}
	return false
}




func ValidLoginForm(req *http.Request) (bool) {
	if req.Form.Get("u") == "" {
		return false
	}
	if req.Form.Get("p") == "" {
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