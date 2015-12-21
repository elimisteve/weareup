// Steven Phillips / elimisteve
// 2015.12.21

package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

const (
	tokensFilename = "./tokens.json"
)

var (
	listenAddr = ":3333"
)

func PostStatus(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading POST: "+err.Error(),
			http.StatusInternalServerError)
		return
	}

	var status Status
	err = json.Unmarshal(b, &status)
	if err != nil && len(b) != 0 {
		http.Error(w, "Error unmarshaling JSON: "+err.Error(),
			http.StatusBadRequest)
		return
	}

	// Set defaults for sender-specifiable fields
	if status.Status == "" {
		status.Status = I_AM_UP
	}
	if status.SentAt.IsZero() {
		status.SentAt = Now()
	}

	// Set server-specified fields
	status.Token = GetToken(r)
	status.ReceivedAt = Now()

	jsonData, _ := json.Marshal(&status)
	log.Printf("New Status: `%s`\n", jsonData)

	// TODO: Store new status using boltDB
	w.Header().Set("Context-Type", "application/json; charset=utf-8")
	w.Write(jsonData)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", PostStatus).Methods("POST")
	http.Handle("/", MiddlewareAuth(r))
	log.Printf("Listening on %v\n", listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}

func Now() time.Time {
	return time.Now().UTC()
}
