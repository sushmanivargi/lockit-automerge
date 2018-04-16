package main

import (
  "log"
  "net/http"
  "os"
)

var (
  config = &Config{}
)

func init() {
  required := []string{
    "LOCKIT_GITHUB_TOKEN",
    "LOCKIT_API_KEY",
    "LOCKIT_HOST",
  }
  for _, field := range required {
    value := os.Getenv(field)
    if value == "0" || value == "" {
      log.Fatalf("[ERROR]: The following *required* env variable is not set: %s\n", field)
    }
  }
  config.Setup()
}

func main() {
  http.HandleFunc("/", basicAuth(router))
  log.Print(http.ListenAndServe(":"+config.LockitAutomerge.Port, logRequest(http.DefaultServeMux)))
}