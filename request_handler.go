package main

import (
	"crypto/subtle"
	"log"
	"net/http"
	"time"
  "io/ioutil"
  "encoding/json"
)

// Function handler that provides HTTP authentication via Basic Auth
func basicAuth(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok ||
			subtle.ConstantTimeCompare([]byte(user), []byte(config.Server.Username)) != 1 ||
			subtle.ConstantTimeCompare([]byte(pass), []byte(config.Server.Password)) != 1 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		handler(w, r)
	}
}

func GetOnly(h http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			h(w, r)
			return
		}
		http.Error(w, "Invalid request method.", http.StatusMethodNotAllowed)
	}
}

func PostOnly(h http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			h(w, r)
			return
		}
		http.Error(w, "Invalid request method.", http.StatusMethodNotAllowed)
	}
}

// Function handler to log incoming requests
func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

// Endpoint for a health check
func handleGetHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{\"status\": \"OK\", \"timestamp\": \"" + time.Now().Format(time.RFC822) + "\"}\n"))
	return
}

func githubHookHandler(w http.ResponseWriter, r *http.Request) {
  r.ParseForm()
  body, err := ioutil.ReadAll(r.Body)
  if err != nil {
    log.Printf("[ERROR] %s", err.Error())
    w.Write([]byte("{\"status\": \"Internal Server Error\"}"))
    w.WriteHeader(http.StatusInternalServerError)
  }
  event := r.Header.Get("X-Github-Event")
  var payload map[string]interface{}
  if err = json.Unmarshal([]byte(body), &payload); err != nil {
    log.Printf("[ERROR] Problem with json.Unmarshal: %s", err.Error())
    w.Write([]byte("{\"status\": \"Internal Server Error\"}"))
    w.WriteHeader(http.StatusInternalServerError)
  }

  if event == "status"{
    err := processStatusChangeWebhook(payload)
    if err != nil{
      w.Write([]byte("{\"status\": \"Internal Server Error\"}"))
      w.WriteHeader(http.StatusInternalServerError)
    }
  }
  w.WriteHeader(http.StatusOK)
}