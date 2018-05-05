package config

import (
	"os"

	//"github.com/coupa/lockit-automerge/github"
)

// Config is a struct to maintain config data
type Config struct {
	Lockit      LockitConfig
	Github      GithubConfig
}

type GithubConfig struct {
	Token          string
	HookSecret     string
	WebhookEnabled string
}

type LockitConfig struct {
	Owner       string
	Repo        string
	Labels      string
	ApiKey      string
	Host        string
}

// Setup converts the env variables into struct values
func Setup() Config{
  c := Config{}
  
	c.Github = GithubConfig{
		Token:           os.Getenv("GITHUB_TOKEN"),
		HookSecret:      os.Getenv("GITHUB_HOOK_SECRET"),
		WebhookEnabled:  os.Getenv("GITHUB_WEBHOOK_ENABLED"),
	}

	c.Lockit = LockitConfig{
		Owner:       os.Getenv("GITHUB_OWNER"),
		Repo:        os.Getenv("GITHUB_REPO"),
		Labels:      os.Getenv("LOCKIT_IGNORE_LABELS"),
		ApiKey:      os.Getenv("LOCKIT_API_KEY"),
		Host:        os.Getenv("LOCKIT_HOST"),
	}

	if c.Lockit.Owner == "" || c.Lockit.Owner == "0" {
		c.Lockit.Owner = "coupa"
	}

	if c.Lockit.Repo == "" || c.Lockit.Repo == "0" {
		c.Lockit.Repo = "coupa_development"
	}

	if c.Lockit.Labels == "" || c.Lockit.Labels == "0" {
		c.Lockit.Labels = "wip"
	}
	return c
}
