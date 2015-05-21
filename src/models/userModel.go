/**
	Contains methods to manipulate 
*/


package models

import (
	db "../database"
	"../helper"
	"gopkg.in/mgo.v2/bson"
	"strings"
	"fmt"
)

type User struct {
	Id 				bson.ObjectId 	`bson:"_id"` /** Automatically generated **/
	/** Manually provided **/
	Firstname 		string			`bson:"fname"`
	Lastname		string			`bson:"lname"`
	Username		string			`bson:"username"`
	Email			string			`bson:"email"`
	Password		string			`bson:"pwd"`
	Zipcode			string 			`bson:"zip"`
	Token			string			`bson:"token"` /** Automatically generated **/
}

/**
	Creates a database entry for user
*/
func SaveUser(user User) (bool) {
	var formattedUser User // Formatted user fields
	formattedUser.Id 		= bson.NewObjectId()
	formattedUser.Firstname = strings.ToLower(user.Firstname)
	formattedUser.Lastname 	= strings.ToLower(user.Lastname)
	formattedUser.Username 	= strings.ToLower(user.Username)
	formattedUser.Email 	= strings.ToLower(user.Email)
	formattedUser.Password 	= helper.MD5String(user.Password)
	formattedUser.Token 	= GenerateActivationToken(formattedUser.Username, formattedUser.Email)
	formattedUser.Zipcode  	= user.Zipcode
	err := db.Conn.DB("dmvexchange").C("USERS").Insert(formattedUser)
	if err != nil {
		fmt.Println("[-] MongoDB error inserting user", err)
		return false
	}
	return true
}

/**
	Determines if a given email + password combination exists
	Returns boolean value
*/
func ValidEmailLogin(email string, password string) (bool) {
	email_formatted := strings.ToLower(email)
	password_formatted := helper.MD5String(password)
	count, err := db.Conn.DB("dmvexchange").C("USERS").Find(bson.M {
		"email": email_formatted,
		"password": password_formatted,
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
func ValidUsernameLogin(username string, password string) (bool) {
	username_formatted := strings.ToLower(username)
	password_formatted := helper.MD5String(password)
	count, err := db.Conn.DB("dmvexchange").C("USERS").Find(bson.M {
		"username": username_formatted,
		"password": password_formatted,
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
	Determines if a username has been registered
	Returns boolean value
*/
func DoesUsernameExist(username string) (bool) {
	formatted := strings.ToLower(username)
	count, err := db.Conn.DB("dmvexchange").C("USERS").Find(bson.M{"username":formatted}).Count()
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
	Returns the document of a user found by username
	Returns User Struct ( Firstname, Lastname, Username etc... )
*/
func FindUserByUsername(username string) (User) {
	formatted := strings.ToLower(username)
	var temp User // Placeholder for user document
	err := db.Conn.DB("dmvexchange").C("USERS").Find(bson.M{"username":formatted}).One(&temp)
	if err != nil {
		fmt.Println("[-] MongoDB error => ", err)
	}
	return temp
}

/**
	Determines if an email address has been registered
	Returns boolean
*/
func DoesEmailExist(email string) (bool) {
	formatted := strings.ToLower(email)
	count, err := db.Conn.DB("dmvexchange").C("USERS").Find(bson.M{"email":formatted}).Count()
	if err != nil {
		fmt.Println("[-] MongoDB error => ", err)
	}
	if count > 0 {
		return true
	}
	return false
}

/**
	Returns the document of a user found by Email
	Returns User struct ( Firstname, Lastname, Username etc... )
*/
func FindUserByEmail(email string) (User) {
	formatted := strings.ToLower(email)
	var temp User // Placeholder for user document
	err := db.Conn.DB("dmvexchange").C("USERS").Find(bson.M{"email":formatted}).One(&temp)
	if err != nil {
		fmt.Println("[-] MongoDB error => ", err)
	}
	return temp
}

/**
	Generates an activation token by hashing a concatenated string
	made of a username, email and string formatted time since unix epoch in seconds
*/
func GenerateActivationToken(username string, email string) (string) {
	return helper.MD5String(username + email)
}