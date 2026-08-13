package service

import (
	"net/url"
)

// func IsPrimitiveIDValid(ID string) bool {

// 	var unknownValue interface{} = ID

// 	_, ok := unknownValue.(primitive.ObjectID)

// 	fmt.Printf("primitive id status %v", ok)
// 	return ok

// }
func IsValidPIN(pin string) bool {

	ok := false

	if len(pin) == 4 {
		ok = true
	}
	return ok
}
func IsValidURL(link string) bool {

	var ok bool
	u, err := url.ParseRequestURI(link)

	if err == nil && u.Scheme != "" {
		ok = true
	}
	return ok

}
