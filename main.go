package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"

	"github.com/coupa/lockit-automerge/auth"
	"github.com/coupa/lockit-automerge/github"
	"github.com/coupa/lockit-automerge/lockit"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron"
)

const lockitFile = ".lockit-cli"

func init() {
	required := []string{
		"LOCKIT_GITHUB_TOKEN",
		"LOCKIT_API_KEY",
		"LOCKIT_HOST",
		"GITHUB_TOKEN",
		"GITHUB_HOOK_SECRET",
		"GITHUB_WEBHOOK_ENABLED",
	}

	for _, field := range required {
		value := os.Getenv(field)
		if value == "0" || value == "" {
			log.Fatalf("[ERROR]: The following *required* env variable is not set: %s\n", field)
		}
	}

	//Enable Jira Integration
	user, err := user.Current()
	if err != nil {
		log.Fatalf("[ERROR]: Fetching current user: %s", err.Error())
	}
	filepath := user.HomeDir + "/" + lockitFile
	data := []byte("{\"jira_enabled\":false}")
	err = ioutil.WriteFile(filepath, data, 0644)
	if err != nil {
		log.Fatalf("[ERROR]: Writing file %s to enable jira integration: %s", filepath, err.Error())
	}

}

func main() {
	scheduleCron()

	router := gin.Default()
	buildRoutes(router)
	http.Handle("/", router)
	router.Run()
}

func buildRoutes(router *gin.Engine) {
	// Webhook routes
	hooks := router.Group("/hooks")
	{
		hooks.Use(auth.GithubMiddleware())
		hooks.POST("/github", github.WebhookHandler)
	}
}

func scheduleCron() {
	//Schedule a cron job to retry lockit auto merge
	c := cron.New()
	c.AddFunc("@every 30m", func() { lockit.RetryMerge() })
	c.Start()
}
