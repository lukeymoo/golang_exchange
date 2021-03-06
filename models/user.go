/**
	Contains methods to manipulate 
*/
package models

import (
	db "../database"
	"gopkg.in/mgo.v2/bson"
	"strings"
	"fmt"
	"crypto/md5"
	"encoding/hex"
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
	formattedUser.Password 	= MD5String(user.Password)
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
	return MD5String(username + email)
}
/**
	Encodes a given string into an MD5 sum and returns it in string format
*/
func MD5String(text string) (string) {
	hashed := md5.New()
	hashed.Write([]byte(text))
	return hex.EncodeToString(hashed.Sum(nil))
}