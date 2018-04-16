package main

import (
	"os"
	"strings"
)

// Config is a struct to maintain config data
type Config struct {
	Server      ServerConfig
	Lockit      LockitConfig
	GithubToken string
}

// ServerConfig holds config data for Server
type ServerConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	UseSSL   string
}

type LockitConfig struct {
	Owner       string
	Repo        string
	GithubToken string
	ApiKey      string
	Host        string
}

// Setup converts the env variables into struct values
func (c *Config) Setup() {
	c.GithubToken = os.Getenv("API_GITHUB_TOKEN")

	c.Server = ServerConfig{
		Host:     os.Getenv("LOCKIT_AUTOMERGE_HOST"),
		Port:     os.Getenv("LOCKIT_AUTOMERGE_PORT"),
		Username: os.Getenv("LOCKIT_AUTOMERGE_USERNAME"),
		Password: os.Getenv("LOCKIT_AUTOMERGE_PASSWORD"),
		UseSSL:   os.Getenv("LOCKIT_AUTOMERGE_USE_SSL"),
	}

	c.Lockit = LockitConfig{
		Owner:       os.Getenv("LOCKIT_AUTOMERGE_OWNER"),
		Repo:        os.Getenv("LOCKIT_AUTOMERGE_REPO"),
		GithubToken: os.Getenv("LOCKIT_GITHUB_TOKEN"),
		ApiKey:      os.Getenv("LOCKIT_API_KEY"),
		Host:        os.Getenv("LOCKIT_HOST"),
	}

	if config.Lockit.Owner == "" || config.Lockit.Owner == "0" {
		config.Lockit.Owner = "coupa"
	}

	if config.Lockit.Repo == "" || config.Lockit.Repo == "0" {
		config.Lockit.Repo = "coupa_development"
	}

	if config.Server.Port == "" || config.Server.Port == "0" {
		config.Server.Port = "8080"
	}
}

// ServerURL returns a canonical string representation of the server location
func (c *Config) ServerURL() string {
	parts := []string{
		"http",
		"://",
		"",
		c.Server.Host,
		"",
	}
	if c.Server.UseSSL == "1" || c.Server.UseSSL == "true" {
		parts[0] = "https"
	}
	if c.Server.Username != "0" && c.Server.Username != "" &&
		c.Server.Password != "0" && c.Server.Password != "" {
		parts[2] = c.Server.Username + ":" + c.Server.Password + "@"
	}
	if c.Server.Port != "" && c.Server.Port != "0" {
		parts[4] = ":" + c.Server.Port
	}

	return strings.Join(parts, "")
}
