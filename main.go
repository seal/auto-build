package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"
)

type Config struct {
	GithubToken   string `json:"githubToken"`
	RepoOwner     string `json:"repoOwner"`
	Reponame      string `json:"repoName"`
	CodeDirectory string `json:"codeDirectory"`
	CodeRun       string `json:"codeRun"`
}

func getConfig() (Config, error) {
	configFile, err := os.Open("config.json")
	if err != nil {
		return Config{}, err
	}
	defer configFile.Close()
	configBody, err := io.ReadAll(configFile)
	if err != nil {
		return Config{}, err
	}
	var c Config
	err = json.Unmarshal(configBody, &c)
	if err != nil {
		return Config{}, err
	}
	return c, nil
}
func GetRecentCommits(c *Config) ([]Commit, error) {
	config, err := getConfig()
	if err != nil {
		return []Commit{}, err
	}
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/commits", c.RepoOwner, c.Reponame)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []Commit{}, err
	}
	req.Header.Set("Authorization", "token "+config.GithubToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return []Commit{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []Commit{}, err
	}

	var commit []Commit
	err = json.Unmarshal(body, &commit)
	if err != nil {
		log.Println(string(body))
		return []Commit{}, err
	}
	return commit, nil

}

var MostRecentSha string

var commands [][]string

func main() {
	config, err := getConfig()
	if err != nil {
		panic(err)
	}
	commands = make([][]string, 4)
	commands[0] = []string{"tmux", "kill-session", "-t", "codeSession"}
	commands[1] = []string{"tmux", "new-session", "-d", "-s", "codeSession", "-n", "tmuxWindow"}
	commands[2] = []string{"tmux", "send-keys", "-t", "codeSession:tmuxWindow", fmt.Sprintf("cd %s", config.CodeDirectory), "Enter"}
	commands[3] = []string{"tmux", "send-keys", "-t", "codeSession:tmuxWindow", config.CodeRun, "Enter"}
	commits, err := GetRecentCommits(&config)
	if err != nil {
		panic(err)
	}
	log.Printf("got %d commits", len(commits))
	if len(commits) == 0 {
		panic(errors.New("no commits"))
	}
	MostRecentSha = commits[0].Sha
	for {
		commits, err := GetRecentCommits(&config)
		if err != nil {
			panic(err)
		}
		if commits[0].Sha == MostRecentSha {
			log.Println("wow new sha")
			for k, v := range commands {
				cmd := exec.Command(v[0], v[1:]...)
				stdout, err := cmd.Output()
				if err != nil {
					if k == 0 { // Kill session doesn't work ( none exists)
						continue
					}
					log.Println(err)
					return
				}
				log.Println("STdout", stdout)
			}
		}
		time.Sleep(10 * time.Second)
	}
}

type Commit struct {
	URL         string `json:"url"`
	Sha         string `json:"sha"`
	NodeID      string `json:"node_id"`
	HTMLURL     string `json:"html_url"`
	CommentsURL string `json:"comments_url"`
	Commit      struct {
		URL    string `json:"url"`
		Author struct {
			Name  string    `json:"name"`
			Email string    `json:"email"`
			Date  time.Time `json:"date"`
		} `json:"author"`
		Committer struct {
			Name  string    `json:"name"`
			Email string    `json:"email"`
			Date  time.Time `json:"date"`
		} `json:"committer"`
		Message string `json:"message"`
		Tree    struct {
			URL string `json:"url"`
			Sha string `json:"sha"`
		} `json:"tree"`
		CommentCount int `json:"comment_count"`
		Verification struct {
			Verified  bool   `json:"verified"`
			Reason    string `json:"reason"`
			Signature any    `json:"signature"`
			Payload   any    `json:"payload"`
		} `json:"verification"`
	} `json:"commit"`
	Author struct {
		Login             string `json:"login"`
		ID                int    `json:"id"`
		NodeID            string `json:"node_id"`
		AvatarURL         string `json:"avatar_url"`
		GravatarID        string `json:"gravatar_id"`
		URL               string `json:"url"`
		HTMLURL           string `json:"html_url"`
		FollowersURL      string `json:"followers_url"`
		FollowingURL      string `json:"following_url"`
		GistsURL          string `json:"gists_url"`
		StarredURL        string `json:"starred_url"`
		SubscriptionsURL  string `json:"subscriptions_url"`
		OrganizationsURL  string `json:"organizations_url"`
		ReposURL          string `json:"repos_url"`
		EventsURL         string `json:"events_url"`
		ReceivedEventsURL string `json:"received_events_url"`
		Type              string `json:"type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"author"`
	Committer struct {
		Login             string `json:"login"`
		ID                int    `json:"id"`
		NodeID            string `json:"node_id"`
		AvatarURL         string `json:"avatar_url"`
		GravatarID        string `json:"gravatar_id"`
		URL               string `json:"url"`
		HTMLURL           string `json:"html_url"`
		FollowersURL      string `json:"followers_url"`
		FollowingURL      string `json:"following_url"`
		GistsURL          string `json:"gists_url"`
		StarredURL        string `json:"starred_url"`
		SubscriptionsURL  string `json:"subscriptions_url"`
		OrganizationsURL  string `json:"organizations_url"`
		ReposURL          string `json:"repos_url"`
		EventsURL         string `json:"events_url"`
		ReceivedEventsURL string `json:"received_events_url"`
		Type              string `json:"type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"committer"`
	Parents []struct {
		URL string `json:"url"`
		Sha string `json:"sha"`
	} `json:"parents"`
}
