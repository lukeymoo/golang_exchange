package session
/**
	Functions for routes
*/

import (
	"fmt"
	redis "github.com/alphazero/Go-Redis"
	"log"
	"time"
	"strconv"
)

var (
	rc redis.Client
)

/** TITLE, PAGE are set manually **/
type TemplateContext struct {
	TITLE string
	PAGE string
	LOGGED_IN string
	USERNAME string
	USER_ID string
	EMAIL string
}

/**
	Initializes redis session client
*/
func InitRedis() {
	spec := redis.DefaultSpec().Db(0).Password("9b3af6edcf71b34520a7d16412ad9325OMGOMG")
	client, err := redis.NewSynchClientWithSpec(spec)
	if err != nil {
		log.Fatal(err)
	}
	rc = client
	return
}

func CreateSessionObj(t *TemplateContext) {
	// if there is a session check for username, userid etc..
	if IsSession() {
		// Set logged in
		t.LOGGED_IN = "true"
		// Check username
		username, err := rc.Get("USERNAME");
		if err != nil {
			log.Fatal(err)
		}
		if username == nil {
			t.USERNAME = ""
		} else {
			if len(string(username)) > 0 {
				t.USERNAME = string(username)
			}
		}
		// Check user id
		userid, err := rc.Get("USER_ID")
		if err != nil {
			log.Fatal(err)
		}
		if userid == nil {
			t.USER_ID = ""
		} else {
			if len(string(userid)) > 0 {
				t.USER_ID = string(userid)
			}
		}
		// check email
		email, err := rc.Get("EMAIL")
		if err != nil {
			log.Fatal(err)
		}
		if email == nil {
			t.EMAIL = ""
		} else {
			if len(string(email)) > 0 {
				t.EMAIL = string(email)
			}
		}
	} else { // If not logged in fill in defaults
		t.LOGGED_IN = "false"
		t.USERNAME = ""
		t.USER_ID = ""
		t.EMAIL = ""
	}
	return
}

// Clears session
func ClearSession() {
	rc.Del("LOGGED_IN")
	rc.Del("LAST_ACTIVITY")
	rc.Del("USERNAME")
	rc.Del("USER_ID")
	rc.Del("EMAIL")
	return
}

// Evaluates LAST_ACTIVITY to determine if still logged in
func IsSession() bool {
	// Check if logged in
	if IsLoggedIn() {
		// Check if timed out
		if IsTimedOut() {
			return false
		} else {
			return true
		}
	}
	return false
}

// Checks if LAST_ACTIVITY is set and if not sets it to unix string literal of unix epoch
func IsTimedOut() bool {
	last_activity, err := rc.Get("LAST_ACTIVITY")
	if err != nil {
		log.Fatal(err)
	}
	// If the last activity was never set, set one
	if last_activity == nil {
		// Update last activity
		UpdateLastActivity()
		ts, err := rc.Get("LAST_ACTIVITY")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Timestamp => %s", ts)
		return false
	}
	return false
}

// Sets LAST_ACTIVITY to string literal of unix epoch
func UpdateLastActivity() {
	rc.Set("LAST_ACTIVITY", []byte(strconv.FormatInt(time.Now().Unix(), 10)));
	return
}

// Checks if LOGGED_IN is set && if not sets it to false
func IsLoggedIn() bool {
	val, err := rc.Get("LOGGED_IN")
	if err != nil {
		log.Fatal(err)
	}
	if val == nil {
		rc.Set("LOGGED_IN", []byte("false"))
		return false
	}
	if string(val) == "true" {
		return true
	}

	if string(val) == "false" {
		return false
	}
	return false
}