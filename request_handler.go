package main

import (
  "crypto/subtle"
  "log"
  "net/http"
  "regexp"
  "time"
)

// Parses and routes the available endpoints
func router(w http.ResponseWriter, r *http.Request) {
  healthPath := regexp.MustCompile(`^/health`)
  lockitMergePath := regexp.MustCompile(`^/merge`)

  switch {
  case healthPath.MatchString(r.URL.Path):
    handleGetHealth(w, r)
  case lockitMergePath.MatchString(r.URL.Path):
    handleGetlockitMerge(w, r)
  default:
    w.WriteHeader(http.StatusNotFound)
  }
}

// Function handler that provides HTTP authentication via Basic Auth
func basicAuth(handler http.HandlerFunc) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    user, pass, ok := r.BasicAuth()
    if !ok ||
      subtle.ConstantTimeCompare([]byte(user), []byte(config.LockitAutomerge.Username)) != 1 ||
      subtle.ConstantTimeCompare([]byte(pass), []byte(config.LockitAutomerge.Password)) != 1 {
      w.WriteHeader(http.StatusUnauthorized)
      return
    }
    handler(w, r)
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

// Endpoint for a health check
func handleGetlockitMerge(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(http.StatusOK)
  w.Write([]byte("{\"status\": \"MERGED\", \"timestamp\": \"" + time.Now().Format(time.RFC822) + "\"}\n"))
  return
}
