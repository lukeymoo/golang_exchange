package helper

import (
	"encoding/hex"
	"crypto/md5"
	"encoding/json"
	"net/http"
	"fmt"
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