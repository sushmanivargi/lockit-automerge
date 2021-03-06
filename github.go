package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strings"

	"github.com/coupa/lockit-core/vcs"
)

func processStatusChangeWebhook(payload map[string]interface{}) (err error) {
	if payload["state"] == "success" {
		if prNumber, err := searchPullRequest(payload["sha"].(string)); err == nil {
			log.Printf("[prNumber] %s", prNumber)
			log.Printf("[initiating lockitMerge] %s", prNumber)
			err := lockitMerge(prNumber)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func searchPullRequest(sha string) (prNumber string, err error) {
	var number string
	data, err := httpHandler("GET", "/search/issues?q="+sha+"+type:pr+limit:1+is:open", nil)
	if err != nil {
		log.Printf("[ERROR] %s", err.Error())
		return "", err
	}
	if data["total_count"].(float64) >= 1 {
		if result := data["items"].([]interface{}); len(result) > 0 {
			number := fmt.Sprint(result[0].(map[string]interface{})["number"].(float64))
			log.Printf("[number] %v", number)
			return number, nil
		}
	}
	return number, nil
}

func httpHandler(method string, path string, data io.Reader) (result map[string]interface{}, err error) {
	client := &http.Client{}
	req, _ := http.NewRequest(method, "https://api.github.com"+path, data)
	req.Header.Add("Authorization", "token "+config.GithubToken)
	log.Printf("[INFO] Request to GitHub: %s", path)

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[ERROR] %s", err.Error())
		return nil, errors.New("Unable to interact with Github: " + err.Error())
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	switch resp.StatusCode {
	case 200, 201, 204, 304:
		json.Unmarshal([]byte(body), &result)
	case 401:
		err = errors.New("(401 Unauthorized) Make sure your credentials are set properly")
	case 404:
		err = errors.New("Not found, aborting...")
	default:
		err = errors.New("Github returned the following: " + string(body))
	}

	return result, err
}

func lockitMerge(prNumber string) (err error) {
	if prNumber == "" || prNumber == "0" {
		return errors.New("Invalid PR number")
	}
	var projectAliases []string
	var gh vcs.Github
	gh.Token = config.Lockit.GithubToken
	log.Printf("Initiating gh.GetPullRequest")
	pr, err := gh.GetPullRequest(config.Lockit.Owner, config.Lockit.Repo, prNumber, projectAliases)
	log.Printf("[pr] %v", pr)
	if err != nil {
		return err
	}

	prInProgress := false
	for _, ele := range strings.Split(config.Lockit.Labels, ",") {
		if contains(pr.Labels, ele) {
			prInProgress = true
			break
		}
	}
	if pr.Target == "master" && prInProgress == false {
		log.Printf("Running: lockit-cli merge %s/%s %s", config.Lockit.Owner, config.Lockit.Repo, pr.Number)

		cmd := exec.Command("lockit-cli", "merge", config.Lockit.Owner+"/"+config.Lockit.Repo, pr.Number)
		out, err := cmd.CombinedOutput()
		if err != nil {
			return err
		}
		if out != nil {
			println("Output: " + string(out))
		}
	}
	return nil
}

func contains(arr []string, str string) bool {
	if len(arr) == 0 {
		return false
	}
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}
