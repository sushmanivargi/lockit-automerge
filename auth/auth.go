package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/coupa/lockit-automerge/config"
	"github.com/coupa/lockit-automerge/github"
	"github.com/coupa/lockit-automerge/github/webhook"
	"github.com/coupa/lockit-automerge/util"
)

// GithubMiddleware is a handler function for authorizing webhook events from Github
func GithubMiddleware() gin.HandlerFunc {
	config := config.Setup()
	return func(gc *gin.Context) {
		if config.Github.WebhookEnabled == "false" {
			gc.AbortWithStatus(http.StatusLocked)
		}
		hook, err := webhook.Parse([]byte(config.Github.HookSecret), gc.Request)
		if err != nil {
			gc.AbortWithStatus(http.StatusForbidden)
		}
		if !util.ArrayContainsString([]string{"status", "issue_comment"}, hook.Event) {
			gc.AbortWithStatus(http.StatusNotAcceptable)
		}
		github.Hook = hook
		gc.Next()
	}
}
