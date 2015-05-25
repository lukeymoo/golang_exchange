package session
/**
	Functions for routes
*/

import (
	redis "github.com/alphazero/Go-Redis"
	"log"
	"strings"
	"time"
	"strconv"
	"fmt"
	"sync"
)

var (
	rc redis.AsyncClient
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
	client, err := redis.NewAsynchClientWithSpec(spec)
	if err != nil {
		log.Fatal(err)
	}
	rc = client
	return
}

func CreateSessionObj(t *TemplateContext) {
	// if there is a session check for username, userid etc..
	if IsSession() {
		var wg sync.WaitGroup

		wg.Add(4)

		go func() {
			defer wg.Done()
			t.LOGGED_IN = GetVariable("LOGGED_IN")
		}()
		go func() {
			defer wg.Done()
			t.USERNAME 	= GetVariable("USERNAME")
		}()
		go func() {
			defer wg.Done()
			t.USER_ID 	= GetVariable("USER_ID")
		}()
		go func() {
			defer wg.Done()
			t.EMAIL 	= GetVariable("EMAIL")
		}()

		wg.Wait()
	} else { // If not logged in fill in defaults
		t.LOGGED_IN = "false"
		t.USERNAME = ""
		t.USER_ID = ""
		t.EMAIL = ""
	}
	return
}

// Sets session, for login/registration
func SetSession(username string, email string, user_id string) (bool) {
	username_formatted := strings.ToLower(username)
	email_formatted := strings.ToLower(email)
	
	var wg sync.WaitGroup

	wg.Add(5)

	go func() {
		defer wg.Done()
		rc.Set("LOGGED_IN", []byte("true"))
	}()
	go func() {
		defer wg.Done()
		rc.Set("USERNAME", []byte(username_formatted))
	}()
	go func() {
		defer wg.Done()
		rc.Set("USER_ID", []byte(user_id))
	}()
	go func() {
		defer wg.Done()
		rc.Set("EMAIL", []byte(email_formatted))
	}()
	go func() {
		defer wg.Done()
		UpdateLastActivity()
	}()

	wg.Wait()
	return true
}

// Clears session
func ClearSession() {
	var wg sync.WaitGroup
	wg.Add(5)
	go func() {
		defer wg.Done()
		rc.Set("LOGGED_IN", []byte("false"))
	}()
	go func() {
		defer wg.Done()
		rc.Set("LAST_ACTIVITY", []byte(""))
	}()
	go func() {
		defer wg.Done()
		rc.Set("USERNAME", []byte(""))
	}()
	go func() {
		defer wg.Done()
		rc.Set("USER_ID", []byte(""))
	}()
	go func() {
		defer wg.Done()
		rc.Set("EMAIL", []byte(""))
	}()
	wg.Wait()
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

// Gets a specified session variable
func GetVariable(variable string) string {
	futureResult, err := rc.Get(variable)
	if err != nil {
		log.Fatal(err)
	}
	result, err := futureResult.Get()
	if err != nil {
		log.Fatal(err)
	}
	if result == nil {
		return ""
	}
	return string(result)
}

// Checks if LAST_ACTIVITY is set and if not sets it to unix string literal of unix epoch
func IsTimedOut() bool {

	last_activity := GetVariable("LAST_ACTIVITY")

	// If the last activity was never set, set one
	if last_activity == "" {
		UpdateLastActivity()
		return false
	} else { // If it was set, determine if timed out ( 3600 seconds )
		current_ts, err := strconv.Atoi(strconv.FormatInt(time.Now().Unix(), 10))
		if err != nil {
			fmt.Println("Session handler error determining if timed out => ", err)
			return true
		}
		last_activity_i, err := strconv.Atoi(GetVariable("LAST_ACTIVITY"))
		if err != nil {
			fmt.Println("Session handler error determing if timed out ( part 2 ) => ", err)
			return true
		}
		if current_ts - last_activity_i > 3600 {
			return true
		} else {
			UpdateLastActivity()
			return false
		}
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
	// Queue
	val, err := rc.Get("LOGGED_IN")
	if err != nil {
		log.Fatal(err)
	}
	// Check response
	resp, err := val.Get()
	if err != nil {
		log.Fatal(err)
	}
	if resp == nil {
		rc.Set("LOGGED_IN", []byte("false"))
		return false
	}
	if string(resp) == "true" {
		return true
	}

	if string(resp) == "false" {
		return false
	}
	return false
}