// Steven Phillips / elimisteve
// 2015.12.21

package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/elimisteve/fun"
	"github.com/gorilla/context"
)

const (
	authHeader = "token"
)

var validTokens []string

func init() {
	body, err := ioutil.ReadFile(tokensFilename)
	if err != nil {
		log.Fatalf("Error reading %s: %v\n", tokensFilename, err)
	}

	var t []string

	err = json.Unmarshal(body, &t)
	if err != nil {
		log.Fatalf("Error reading tokens file: %v\n", err)
	}

	// Sets global var
	validTokens = t
}

func MiddlewareAuth(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check header
		token := r.Header.Get(authHeader)
		if !fun.SliceContains(validTokens, token) {
			// Check GET param
			r.ParseForm()
			token = r.Form.Get(authHeader)
			if !fun.SliceContains(validTokens, token) {
				http.Error(w, "Invalid token '"+token+"'", http.StatusUnauthorized)
				return
			}
		}

		log.Printf("%v %v from %v\n", r.URL.Path, r.Method, token)
		SetToken(r, token)
		h.ServeHTTP(w, r)
	})
}

func GetToken(r *http.Request) string {
	if t := context.Get(r, authHeader); t != nil {
		return t.(string)
	}
	return ""
}

func SetToken(r *http.Request, token string) {
	context.Set(r, authHeader, token)
}
