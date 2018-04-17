# lockit-automerge
## Auto merge PR with Lockit

### **Features**

1) Auto merges PR with lockit cli merge command if all the checks on the PR are green
  * Only PRs targeting master will be auto-merged.
  * If the PR contains any of the following labels it will not be auto merged: “wip”, “pending”, “DO NOT MERGE”, “on hold”, “rework request”
  * Auto merge will work only if the required fields below are updated in the jira ticket
    - Documentation Impact
    - For CD bugs, build failure, or JZ tickets:
      - Root cause
      - Regression
      - Root cause summary
2) Jira integration enabled by default to resolve jira tickets.
3) Supports lockit merge with comment "lockit merge" on a PR.
4) If master is not in merge-able state, pushes PR to the retry queue.
5) Retries merging PR from the queue every 30 minutes through a cron job.