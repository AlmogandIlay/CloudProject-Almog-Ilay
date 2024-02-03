package helper

import (
	"CloudDrive/authentication"
	"encoding/json"
)

// Encodes a json user request to user struct
func GetEncodedUser(userJson json.RawMessage) authentication.User {
	var user authentication.User
	json.Unmarshal([]byte(userJson), &user) // Json encoding
	return user
}
