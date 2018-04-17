package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/coupa/lockit-automerge/config"
	"github.com/coupa/lockit-automerge/github/webhook"
	"github.com/coupa/lockit-automerge/lockit"
	"github.com/coupa/lockit-automerge/queue"
	"github.com/coupa/lockit-automerge/util"
)

var (
	Hook *webhook.Hook
)

func WebhookHandler(gc *gin.Context) {
	var payload map[string]interface{}
	if err := json.Unmarshal(Hook.Payload, &payload); err != nil {
		log.Printf("[ERROR] json.Unmarshal: %s", err.Error())
		gc.AbortWithStatus(http.StatusInternalServerError)
	}
	switch Hook.Event {
	case "status":
		defer processStatusChangeWebhook(gc, payload)
	case "issue_comment":
		defer processIssueCommentWebhook(gc, payload)
	}
	gc.Writer.WriteHeader(http.StatusOK)
}

func processStatusChangeWebhook(gc *gin.Context, payload map[string]interface{}) {
	if payload["state"] == "success" {
		prNumber, err := searchPullRequestNumber(gc, payload["sha"].(string))
		if err != nil {
			log.Printf("[ERROR] Searching PR number: %s", err.Error())
			return
		}
		if !queue.IsEmpty() {
			if err = queue.DeleteItem(prNumber); err != nil { //Remove from retry queue, if already present
				log.Printf("[ERROR] Deleting from retry queue: %s", err.Error())
			}
		}
		if err := lockit.Merge(prNumber); err != nil {
			log.Printf("[ERROR] Merging PR: %s", err.Error())
		}
	}
}

func processIssueCommentWebhook(gc *gin.Context, payload map[string]interface{}) {
	issue := payload["issue"].(map[string]interface{})
	comment := payload["comment"].(map[string]interface{})
	prNumber := strconv.FormatFloat(issue["number"].(float64), 'f', 0, 64)

	if payload["action"] != "created" || issue["state"] != "open" {
		return
	}
	if strings.ToLower(comment["body"].(string)) == "lockit merge" {
		if !queue.IsEmpty() {
			if err := queue.DeleteItem(prNumber); err != nil { //Remove from retry queue, if already present
				log.Printf("[ERROR] Deleting from retry queue: %s", err.Error())
			}
		}
		if err := lockit.Merge(prNumber); err != nil {
			log.Printf("[ERROR] Merging PR: %s", err.Error())
		}
	}
}

func searchPullRequestNumber(gc *gin.Context, sha string) (prNumber string, err error) {
	data, err := httpHandler(gc, "GET", "/search/issues?q="+sha+"+type:pr+limit:1+is:open", nil)
	if util.IsError(err) {
		return
	}
	if data["total_count"].(float64) >= 1 {
		if result := data["items"].([]interface{}); len(result) > 0 {
			prNumber = fmt.Sprint(result[0].(map[string]interface{})["number"].(float64))
			return
		}
	}
	return
}

func httpHandler(gc *gin.Context, method string, path string, data io.Reader) (result map[string]interface{}, err error) {
	config := config.Setup()
	var errMessage string
	var body []byte

	client := &http.Client{}
	req, _ := http.NewRequest(method, "https://api.github.com"+path, data)
	req.Header.Add("Authorization", "token "+config.Github.Token)
	log.Printf("[INFO] Request sent: %s", path)

	resp, err := client.Do(req)
	if err != nil {
		errMessage = "Unable to send request: " + err.Error()
		goto FAIL
	}

	defer resp.Body.Close()
	body, _ = ioutil.ReadAll(resp.Body)

	if !util.ArrayContainsInt([]int{200, 201, 204, 304}, resp.StatusCode) {
		errMessage = "Received unexpected status code: " + strconv.Itoa(resp.StatusCode)
		goto FAIL
	}
	json.Unmarshal([]byte(body), &result)

FAIL:
	if errMessage != "" {
		return nil, errors.New(errMessage)
	}

	return
}
