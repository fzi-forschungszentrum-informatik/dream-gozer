package webapi

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

import (
	"github.com/fzi-forschungszentrum-informatik/gozer/model"
	"github.com/fzi-forschungszentrum-informatik/gozer/storage"
	"github.com/google/uuid"
)

// authorizationHandler encapsulates a Web request handler that requires user authentication.
func authorizationHandler(handler func(http.ResponseWriter, *http.Request, *model.User), st *storage.Storage) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		startTime := time.Now()

		guid, secret, ok := r.BasicAuth()

		if !ok {
			log.Printf("Autorization failue. Request seems to miss HTTP BasicAuth information.")
			http.Error(w, "Authorization Error", http.StatusUnauthorized)
			return
		}

		_, err := uuid.Parse(guid)

		if err != nil {
			log.Printf("Autorization failue. GUID '%s' seems to be malformed. Could not be parsed.", guid)
			http.Error(w, "Authorization Error", http.StatusUnauthorized)
			return
		}

		user, err := st.UserByGUID(guid)

		if err != nil {
			handleInternalError(w, "Internal database error while reading user credentials.", err)
			return
		}

		if user == nil {
			log.Printf("Autorization failue. User with GUID '%s' does not exist.", guid)
			http.Error(w, "Authorization Error", http.StatusUnauthorized)
			return
		}

		err = user.Authorize(secret)

		if err != nil {
			log.Printf("Authorization failure. %s", err)
			http.Error(w, "Authorization Error", http.StatusUnauthorized)
			return
		}

		log.Printf("Processing authorized %s-request on '%s' with GUID '%s'.", r.Method, r.URL, user.GUID)

		handler(w, r, user)

		log.Printf("Total response time: %v", time.Since(startTime))
	}
}

// authorizationHandler encapsulates a Web request handler that requires no user authentication.
func defaultHandler(handler func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Processing %s-request on '%s'.", r.Method, r.URL)
		startTime := time.Now()
		handler(w, r)
		log.Printf("Total response time: %v", time.Since(startTime))
	}
}

// handleBadRequest writes a standard response to the client in the case the request has unexpected or missing parameters.
func handleBadRequest(w http.ResponseWriter, msg string) {
	log.Print(msg)
	http.Error(w, "Not Implemented Error", http.StatusBadRequest)
}

// handleBadRequest writes a standard response to the client in the case a handler fails to complete.
func handleInternalError(w http.ResponseWriter, msg string, err error) {
	log.Printf("%s %s", msg, err)
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}

// readRequest reads any type of a JSON datastructure from the body of a HTTP request and unmarshals it to specified request data structure.
func readRequest(w http.ResponseWriter, r *http.Request, request interface{}) (err error) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Could not read message body. %s", err)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, request)
	if err != nil {
		log.Printf("Could not unmarshal JSON request from HTTP body. %s", err)
		return
	}

	return nil
}

// writeResponse marshals any type of data structure to JSON and writes it to the body of a HTTP response.
func writeResponse(w http.ResponseWriter, response interface{}) error {

	w.Header().Set("Content-Type", "application/json")

	jData, err := json.Marshal(response)
	if err != nil {
		log.Printf("JSON marshalling for HTTP response failed. %s", err)
		return err
	}

	if _, err = w.Write(jData); err != nil {
		log.Printf("Could not write response to HTTP ResponseWriter. %s", err)
		return err
	}

	return nil
}
