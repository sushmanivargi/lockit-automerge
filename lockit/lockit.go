package lockit

import (
	"errors"
	"log"
	"os/exec"
	"strings"

	"github.com/coupa/lockit-automerge/config"
	"github.com/coupa/lockit-automerge/queue"
	"github.com/coupa/lockit-automerge/util"
	"github.com/coupa/lockit-core/vcs"
)

func Merge(prNumber string) (err error) {
	log.Printf("[INFO] Initiating merge for #%s", prNumber)
	config := config.Setup()
	var gh vcs.Github
	var projectAliases []string
	gh.Token = config.Github.Token
	pr, err := gh.GetPullRequest(config.Lockit.Owner, config.Lockit.Repo, prNumber, projectAliases)
	if util.IsError(err) {
		return
	}

	prBlockedForMerge := false
	for _, ele := range strings.Split(config.Lockit.Labels, ",") {
		if util.ArrayContainsString(pr.Labels, ele) {
			prBlockedForMerge = true
			break
		}
	}
	if pr.Target != "master" || prBlockedForMerge == true {
		return errors.New("Not in merge-able state.")
	}

	cmd := exec.Command("lockit-cli", "merge", config.Lockit.Owner+"/"+config.Lockit.Repo, pr.Number)
	_, err = cmd.CombinedOutput()
	// Only to debug the output of lockit merge command
	// if out != nil {
	//  log.Printf("Output: %v", string(out))
	// }
	if err != nil {
		log.Printf("[INFO] Pushing PR #%s to retry queue.", pr.Number)
		err = queue.Enqueue(pr.Number)
		if util.IsError(err) {
			return
		}
	} else {
		log.Printf("[INFO] Merge complete for PR #%s", prNumber)
	}
	return
}

func RetryMerge() {
	prNumber, err := queue.Dequeue()
	if err == nil {
		if err = Merge(prNumber); err != nil {
			log.Printf("[ERROR] PR merge: %s", err.Error())
		}
		RetryMerge()
	}
}
