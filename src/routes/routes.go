package routes
/**
	Functions for routes
*/

import (
	"fmt"
	"net/http"
	"github.com/julienschmidt/httprouter"
	redis "github.com/alphazero/Go-Redis"
	"html/template"
	"regexp"
	"log"
	"time"
	"encoding/binary"
	"strconv"
);

var (
	rc redis.Client
);

/**
	Initializes redis session client
*/
func InitRedis() {
	spec := redis.DefaultSpec().Db(0).Password("9b3af6edcf71b34520a7d16412ad9325OMGOMG");
	client, err := redis.NewSynchClientWithSpec(spec);
	if err != nil {
		log.Fatal(err);
	}
	rc = client;
	return;
}

/**
	GET
	/
	Display home page
*/
func IndexPage(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	tmp, err := template.ParseFiles("template/layout/master.html", "template/index.html");
	if err != nil {
		fmt.Println("Error =>", err);
		fmt.Fprint(res, "Server error occurred");
		return;
	}
	var context = map[string] string {
		"TITLE": "Home",
		"PAGE": "HOME",
		"LOGGED_IN": "false",
	};
	err = tmp.Execute(res, context);
	if err != nil {
		fmt.Println("Error =>", err);
		return;
	}
	return;
}

/**
	GET
	/register
	Display register page
*/
func RegisterPage(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	tmp, err := template.ParseFiles("template/layout/master.html", "template/index.html");
	if err != nil {
		fmt.Println("Error =>", err);
		fmt.Fprint(res, "Server error occurred");
		return;
	}
	var context = map[string] string {
		"TITLE": "Register",
		"PAGE": "REGISTER",
		"LOGGED_IN": "false",
	};
	err = tmp.Execute(res, context);
	if err != nil {
		fmt.Println("Server error occurred");
		return;
	}
	return;
}

/**
	POST
	/register/process
	Processes register form & creates an account
*/
func ProcessRegister(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	// Ensure not logged in
	if isSession() {
		http.Redirect(res, req, "/register?err=logged_in", 301);
		return;
	}
	// Ensure we received all register form elements
	formError := req.ParseForm()
	if formError != nil {
		fmt.Println("Error =>", formError);
		fmt.Fprint(res, "Server error");
		return;
	}

	// Ensure we recieved all form fields
	if !validRegisterForm(req) {
		http.Redirect(res, req, "/register?err=invalid_form", 301);
		return;
	}

	err := ""
	query := "&f=" + req.Form.Get("f") + "&l=" + req.Form.Get("l") +
			"&u=" + req.Form.Get("u") + "&e=" + req.Form.Get("e") + "&z=" + req.Form.Get("z") +
			"&g=" + req.Form.Get("g");
	
	// Validate names
	if !validName(req.Form.Get("f")) {
		err += "F|";
	}
	if !validName(req.Form.Get("l")) {
		err += "L|";
	}
	// Validate Username
	if !validUsername(req.Form.Get("u")) {
		err += "U|";
	}
	// Password
	if !validPassword(req.Form.Get("p")) {
		err += "P|";
	}
	// Password again
	if req.Form.Get("p") != req.Form.Get("pa") {
		err += "PM|";
	}
	// Email
	if !validEmail(req.Form.Get("e")) {
		err += "E|";
	}
	// Email again
	if req.Form.Get("e") != req.Form.Get("ea") {
		err += "EM|";
	}
	// Zipcode
	if !validZipcode(req.Form.Get("z")) {
		err += "Z|";
	}
	// If errors redirect with error codes
	if err != "" {
		err = err[:len(err)-1];
		http.Redirect(res, req, "/register?err=" + err + query, 301);
		return;
	}

	/** Determine if the username and/or email is in use **/
	if isSession() {
		if isTimedOut() {
			fmt.Fprint(res, "Session has timed out");
		} else {
			fmt.Fprint(res, "Logged in");
		}
	} else {
		fmt.Fprint(res, "Not logged in");
	}
	return;
}



/**
	GET
	/login/logout
	Logs the user out
*/
func Logout(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	rc.Del("LOGGED_IN");
	rc.Del("LAST_ACTIVITY");
	rc.Del("USERNAME");
	rc.Del("USER_ID");
	rc.Del("EMAIL");
	http.Redirect(res, req, "/", 301);
	return
}




func FakeLogin(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	rc.Set("LOGGED_IN", []byte("true"));
	timestamp := []byte(strconv.FormatInt(time.Now().Unix(), 10));
	rc.Set("LAST_ACTIVITY", timestamp);
	fmt.Fprint(res, "Fake session values set!");
	return
}

func Debug(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
	logged_in, err := rc.Get("LOGGED_IN");
	if err != nil {
		log.Fatal(err);
	}
	logged_in_s := string(logged_in);
	fmt.Println("Logged in =>", logged_in_s);

	last_activity, err := rc.Get("LAST_ACTIVITY");
	if err != nil {
		log.Fatal(err);
	}
	last_activity_s := string(last_activity);
	fmt.Println("Last activity => ", last_activity_s);
	fmt.Fprint(res, "Debug outputted to console!");
	return;
}









/** Helpers **/

func isSession() bool {
	// Check if logged in
	if isLoggedIn() {
		// Check if timed out
		if isTimedOut() {
			return false;
		} else {
			return true;
		}
	}
	return false;
}

// Checks if LAST_ACTIVITY is set and if not sets it to time.Now().Unix()
func isTimedOut() bool {
	last_activity, err := rc.Get("LAST_ACTIVITY");
	if err != nil {
		log.Fatal(err);
	}
	// If the last activity was never set, set one
	if last_activity == nil {
		// Update last activity
		updateLastActivity();
		ts, err := rc.Get("LAST_ACTIVITY");
		if err != nil {
			log.Fatal(err);
		}
		fmt.Printf("Timestamp => %s", ts);
		return false;
	}
	return false;
}

func updateLastActivity() {
	ts_int := time.Now().Unix(); // signed int64 unix epoch
	var buf = make([]byte, 8); // create space for conversion
	binary.PutVarint(buf, ts_int); // read into buf
	rc.Set("LAST_ACTIVITY", buf); // update session variable
	return;
}

// Checks if LOGGED_IN is set && if not sets it to false
func isLoggedIn() bool {
	val, err := rc.Get("LOGGED_IN");
	if err != nil {
		log.Fatal(err);
	}
	if val == nil {
		rc.Set("LOGGED_IN", []byte("false"));
		return false;
	}
	if string(val) == "true" {
		return true;
	}

	if string(val) == "false" {
		return false;
	}
	return false;
}


func validRegisterForm(req *http.Request) bool {

	/** Ensure fields are set and not empty **/
	// First name
	if req.Form.Get("f") == "" {
		return false;
	}
	// Last name
	if req.Form.Get("l") == "" {
		return false;
	}
	// Username
	if req.Form.Get("u") == "" {
		return false;
	}
	// Password
	if req.Form.Get("p") == "" {
		return false;
	}
	// Password again
	if req.Form.Get("pa") == "" {
		return false;
	}
	// Email
	if req.Form.Get("e") == "" {
		return false;
	}
	// Email again
	if req.Form.Get("ea") == "" {
		return false;
	}
	// Zipcode
	if req.Form.Get("z") == "" {
		return false;
	}
	// Gender must be set
	if req.Form.Get("g") == "" {
		return false;
	}
	// Gender must be 'm' or 'f'
	if req.Form.Get("g") != "m" && req.Form.Get("g") != "f" {
		return false;
	}
	// Terms and conditions must be I_AGREE
	if req.Form.Get("tos") != "I_AGREE" {
		return false;
	}

	return true;
}

// Validate Name
func validName(name string) bool {
	// Regex
	match, err := regexp.MatchString(`^[A-Za-z]+(([\'-])?[A-Za-z]+$)`, name);
	if err != nil {
		log.Fatal(err);
	}
	if !match {
		return false;
	}

	// Test length
	if len(name) < 2 || len(name) > 32 {
		return false;
	}


	return true;
}

// Validate Username
func validUsername(name string) bool {
	// Regex
	match, err := regexp.MatchString(`^[A-Za-z0-9_]+$`, name);
	if err != nil {
		log.Fatal(err);
	}
	if !match {
		return false;
	}
	// Test length
	if len(name) < 2 || len(name) > 16 {
		return false;
	}
	return true;
}

// Validate Password
func validPassword(pwd string) bool {
	if len(pwd) < 2 || len(pwd) > 32 {
		return false;
	}
	return true;
}

// Validate email
func validEmail(email string) bool {
	match, err := regexp.MatchString(`^([a-zA-Z0-9_.+-])+\@(([a-zA-Z0-9-])+\.)+([a-zA-Z0-9]{2,4})+$`, email);
	if err != nil {
		log.Fatal(err);
	}
	if !match {
		return false;
	}
	if len(email) > 64 {
		return false;
	}
	return true;
}

// Validate zipcode
func validZipcode(code string) bool {
	match, err := regexp.MatchString(`^[0-9]+$`, code);
	if err != nil {
		log.Fatal(err);
	}
	if !match {
		return false;
	}
	if len(code) != 5 {
		return false;
	}
	return true;
}