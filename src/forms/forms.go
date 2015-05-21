package forms

import (
	"net/http"
	"regexp"
	"log"
)

func ValidRegisterForm(req *http.Request) bool {

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
func ValidName(name string) bool {
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