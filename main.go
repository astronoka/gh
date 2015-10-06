package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/url"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

var (
	githubHost  = ""
	accessToken = ""
	owner       = ""
	repository  = ""
	sha         = ""
	context     = ""
	targeturl   = ""
	description = ""
	status      = ""

	display = true
)

func init() {
	flag.StringVar(&githubHost, "host", "", "Github host. (api endpoint)")
	flag.StringVar(&accessToken, "token", "", "Required. Access token")
	flag.StringVar(&owner, "owner", "", "Required. Owner. (or Organaization)")
	flag.StringVar(&repository, "repo", "", "Required. Repository name.")
	flag.StringVar(&sha, "sha", "", "Required. Commit SHA.")
	flag.StringVar(&context, "context", "", "A string label to differentiate this status from the status of other systems.")
	flag.StringVar(&targeturl, "targeturl", "", "The target URL to associate with this status.")
	flag.StringVar(&description, "desc", "", "A short description of the status.")
	flag.StringVar(&status, "status", "success", `Required. The state of the status. Can be one of pending, success, error, or failure.`)

	flag.BoolVar(&display, "display", true, `Display commit status json.`)
}

func main() {

	flag.Parse()

	if display {
		doDisplay()
		return
	}
	doUpdate()
}

func doDisplay() {
	if accessToken == "" || owner == "" || repository == "" || sha == "" {
		flag.Usage()
		return
	}

	client := buildClient()

	result, _, err := GitCommitStatus(client, owner, repository, sha)
	if err != nil {
		fmt.Fprintf(os.Stderr, "gh: %v\n", err)
		os.Exit(1)
	}
	printjson(result)
}

func doUpdate() {

	if accessToken == "" || owner == "" ||
		repository == "" || sha == "" || status == "" {
		flag.Usage()
		return
	}

	client := buildClient()

	err := updateBuildStatus(client)
	if err != nil {
		fmt.Fprintf(os.Stderr, "gh: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("gh: successfully updated")
}

func buildClient() *github.Client {

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	client := github.NewClient(tc)
	if githubHost != "" {
		var err error
		client.BaseURL, err = url.Parse(githubHost)
		if err != nil {
			fmt.Fprintf(os.Stderr, "gh: %v\n", err)
			os.Exit(1)
		}
	}

	return client
}

type CommitStatus struct {
	State       string `json:"state,omitempty"`
	TargetURL   string `json:"target_url,omitempty"`
	Description string `json:"description,omitempty"`
	Context     string `json:"context,omitempty"`
}

func updateBuildStatus(client *github.Client) error {
	s := &CommitStatus{
		State:       status,
		TargetURL:   targeturl,
		Description: description,
		Context:     context,
	}
	_, _, err := PutGitCommitStatus(client, owner, repository, sha, s)
	if err != nil {
		return err
	}
	return nil
}

func PutGitCommitStatus(client *github.Client, owner, repo, sha string, body *CommitStatus) (interface{}, *github.Response, error) {

	u := fmt.Sprintf("repos/%v/%v/statuses/%v", owner, repo, sha)

	if body == nil {
		body = &CommitStatus{}
	}

	req, err := client.NewRequest("POST", u, body)
	if err != nil {
		return nil, nil, err
	}

	var c interface{}
	resp, err := client.Do(req, c)
	if err != nil {
		return nil, resp, err
	}

	return c, resp, err
}

func GitCommitStatus(client *github.Client, owner, repo, sha string) (interface{}, *github.Response, error) {

	u := fmt.Sprintf("repos/%v/%v/commits/%v/status", owner, repo, sha)
	req, err := client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var c interface{}
	resp, err := client.Do(req, &c)
	if err != nil {
		return nil, resp, err
	}
	return c, resp, err
}

func GitCommitStatuses(client *github.Client, owner, repo, sha string) (interface{}, *github.Response, error) {

	u := fmt.Sprintf("repos/%v/%v/commits/%v/statuses", owner, repo, sha)
	req, err := client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var c interface{}
	resp, err := client.Do(req, &c)
	if err != nil {
		return nil, resp, err
	}
	return c, resp, err
}

func printjson(v interface{}) {
	j, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(j))
}
