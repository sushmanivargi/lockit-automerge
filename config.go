package main

import (
  "os"
  "strings"
)

// Config is a struct to maintain config data
type Config struct {
  LockitAutomerge LockitAutomergeConfig
}


// LockitAutomergeConfig holds config data for LockitAutomerge
type LockitAutomergeConfig struct {
  Host      string
  Port      string
  Username  string
  Password  string
  UseSSL    string
}

// Setup converts the env variables into struct values
func (c *Config) Setup() {
  c.LockitAutomerge = LockitAutomergeConfig{
    Host:      os.Getenv("LOCKITAUTOMERGE_HOST"),
    Port:      os.Getenv("LOCKITAUTOMERGE_PORT"),
    Username:  os.Getenv("LOCKITAUTOMERGE_USERNAME"),
    Password:  os.Getenv("LOCKITAUTOMERGE_PASSWORD"),
    UseSSL:    os.Getenv("LOCKITAUTOMERGE_USE_SSL"),
  }

  if config.LockitAutomerge.Port == "" || config.LockitAutomerge.Port == "0" {
    config.LockitAutomerge.Port = "8080"
  }
}

// ServerURL returns a canonical string representation of the server location
func (c *Config) ServerURL() string {
  parts := []string{
    "http",
    "://",
    "",
    c.LockitAutomerge.Host,
    "",
  }
  if c.LockitAutomerge.UseSSL == "1" || c.LockitAutomerge.UseSSL == "true" {
    parts[0] = "https"
  }
  if c.LockitAutomerge.Username != "0" && c.LockitAutomerge.Username != "" &&
    c.LockitAutomerge.Password != "0" && c.LockitAutomerge.Password != "" {
    parts[2] = c.LockitAutomerge.Username + ":" + c.LockitAutomerge.Password + "@"
  }
  if c.LockitAutomerge.Port != "" && c.LockitAutomerge.Port != "0" {
    parts[4] = ":" + c.LockitAutomerge.Port
  }

  return strings.Join(parts, "")
}