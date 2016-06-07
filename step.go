package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"regexp"
	"strconv"
)

// -----------------------
// --- Constants
// -----------------------

const (
	envPullRequestTitle = "GPRI_PULL_REQUEST_TITLE"
	envOtherPullRequestTitles = "GPRI_OTHER_PULL_REQUEST_TITLES"
	envPullRequestLabels = "GPRI_PULL_REQUEST_LABELS"
	envPullRequestNumber = "GPRI_PULL_REQUEST_NUMBER"
)

// -----------------------
// --- Models
// -----------------------

//GithubEvent from GET /repos/:owner/:repo/events
type GithubEvent struct {
	ID int `json:"id"`
	EventType string `json:"event"`
	CommitID string `json:"commit_id"`
	Issue GithubIssue
}

//GithubIssue inside of GithubEvent
type GithubIssue struct {
	Number int
	Title string
	Labels []GithubLabel
}

//GithubLabel ...
type GithubLabel struct {
	Name string
}

// -----------------------
// --- Functions
// -----------------------

func logFail(format string, v ...interface{}) {
	errorMsg := fmt.Sprintf(format, v...)
	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", errorMsg)
	os.Exit(1)
}

func logWarn(format string, v ...interface{}) {
	errorMsg := fmt.Sprintf(format, v...)
	fmt.Printf("\x1b[33;1m%s\x1b[0m\n", errorMsg)
}

func logInfo(format string, v ...interface{}) {
	fmt.Println()
	errorMsg := fmt.Sprintf(format, v...)
	fmt.Printf("\x1b[34;1m%s\x1b[0m\n", errorMsg)
}

func logDetails(format string, v ...interface{}) {
	errorMsg := fmt.Sprintf(format, v...)
	fmt.Printf("  %s\n", errorMsg)
}

func logDone(format string, v ...interface{}) {
	errorMsg := fmt.Sprintf(format, v...)
	fmt.Printf("  \x1b[32;1m%s\x1b[0m\n", errorMsg)
}

func exportEnvironmentWithEnvman(keyStr, valueStr string) error {
	envman := exec.Command("envman", "add", "--key", keyStr)
	envman.Stdin = strings.NewReader(valueStr)
	envman.Stdout = os.Stdout
	envman.Stderr = os.Stderr
	return envman.Run()
}

// -----------------------
// --- Main
// -----------------------

func main() {
	//
	// Validate options
	githubUsername := os.Getenv("github_username")
	accessToken := os.Getenv("github_access_token")
	commitSHA := os.Getenv("commit_sha")
	isPR := os.Getenv("is_PR") == "true"
        pullRequestID := os.Getenv("pull_request_id")
	repositoryURL := os.Getenv("github_repository_url")

	logInfo("Configs:")
	logDetails("github_username: %s", githubUsername)
	logDetails("github_access_token: ***")
	logDetails("commit_sha: %s", commitSHA)
        logDetails("is_PR: %t", isPR)
	logDetails("pull_request_id: %s", pullRequestID)
	logDetails("repository_url: %s", repositoryURL)


	if githubUsername == "" {
		logFail("No App github_username provided as environment variable. Terminating...")
	}

	if accessToken == "" {
		logFail("No App github_access_token provided as environment variable. Terminating...")
	}

	if commitSHA == "" {
		logFail("commitSHA is empty!")
	}

	if repositoryURL == "" {
		logFail("repositoryURL is empty!")
	}

        //Parse git repo address
        var githubRepoURLRe = regexp.MustCompile(`:(.+)/(.+)\.git$`)
	var githubRepoOwner string
	var githubRepoName string
	
        if githubRepoURLParts := githubRepoURLRe.FindStringSubmatch(repositoryURL); githubRepoURLParts != nil {
		githubRepoOwner = githubRepoURLParts[1]
		githubRepoName = githubRepoURLParts[2]
        } else {
		logFail("This doesn't look like a github repository: %s", repositoryURL)
	}

	logInfo("Github repository is: %s/%s", githubRepoOwner, githubRepoName)

	
	//
	// Get Pull Request ID and recent PR titles
	githubClient := Github{githubUsername, accessToken}
	mergeEvents := githubClient.fetchMergeEvents(githubRepoOwner, githubRepoName)

	var mergeEvent GithubEvent
	otherChanges := ""
	
	for _, event := range mergeEvents {
		if event.CommitID == commitSHA {
			mergeEvent = event
		} else {
			otherChanges += mergeEvent.Issue.Title + "\n"
		}
	}

	prettyPrint(mergeEvent)
	fmt.Println(otherChanges)

	//
	// Labels
	labels := ""
	for index, label := range mergeEvent.Issue.Labels {
		labels += label.Name
		if index + 1 < len(mergeEvent.Issue.Labels) {
			labels += ":"
		}
	}
	
	//
	// Export results
	exportEnvironmentWithEnvman(envPullRequestTitle, mergeEvent.Issue.Title)
	exportEnvironmentWithEnvman(envOtherPullRequestTitles, otherChanges)
	exportEnvironmentWithEnvman(envPullRequestLabels, labels)
	exportEnvironmentWithEnvman(envPullRequestNumber, strconv.Itoa(mergeEvent.Issue.Number))

}

//Github client that stores a username and access token
type Github struct {
	Username string
	AccessToken string
}

func (githubClient *Github) fetchMergeEvents(username string, repository string) []GithubEvent {
	httpClient := http.Client{}
	fetchLimit := 30
	
	req, _ := http.NewRequest("GET", "https://api.github.com/repos/" + username + "/" + repository + "/issues/events?per_page=" + strconv.Itoa(fetchLimit), nil)
	req.SetBasicAuth(githubClient.Username, githubClient.AccessToken)
	req.Header.Set("User-Agent", "Step for Bitrise: Pull Request information https://github.com/philippelatulippe/steps-pull-request-information")
	
	response, requestErr := httpClient.Do(req)
	defer response.Body.Close()

	if requestErr != nil {
		logFail("Failed to fetch github events: %s", requestErr)
	}

	if response.StatusCode != 200 {
		if response.ContentLength == 0 {
			logFail("HTTP error %s", response.Status)
		} else {
			var v map[string]interface{}
			json.NewDecoder(response.Body).Decode(&v)
			logFail("HTTP error %s; error message: %s", response.Status, v["message"])
		}
	}
	
	decoder := json.NewDecoder(response.Body)

	//Read the JSON array open bracket, so that we can stream the array's elements
	token, tokenErr := decoder.Token()
	if tokenErr != nil {
		logFail("Failed to parse github events: %s", tokenErr)
	}
	if delimiter, ok := token.(json.Delim); !ok || delimiter != '[' {
		logFail("Failed to parse github events: %s", "response not an array")
	}


	events := make([]GithubEvent, 0, fetchLimit)
	var event GithubEvent
	for decoder.More() {
		if err := decoder.Decode(&event); err != nil {
			logFail("Failed to parse github events: %s", err)
		}

		if event.EventType == "merged" {
			events = append(events, event)
		}
	}

	return events
}

//TODO:delete
func prettyPrint(something interface{}) {
	b, err := json.MarshalIndent(something, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println(string(b))
}
