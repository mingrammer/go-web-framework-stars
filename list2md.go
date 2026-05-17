package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
	"unicode"
)

// Repo describes a Github repository with additional field, last commit date
type Repo struct {
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	DefaultBranch  string    `json:"default_branch"`
	Stars          int       `json:"stargazers_count"`
	Forks          int       `json:"forks_count"`
	Issues         int       `json:"open_issues_count"`
	Created        time.Time `json:"created_at"`
	Updated        time.Time `json:"updated_at"`
	URL            string    `json:"html_url"`
	LastCommitDate time.Time `json:"-"`
}

// HeadCommit describes a head commit of default branch
type HeadCommit struct {
	Sha    string `json:"sha"`
	Commit struct {
		Committer struct {
			Name  string    `json:"name"`
			Email string    `json:"email"`
			Date  time.Time `json:"date"`
		} `json:"committer"`
	} `json:"commit"`
}

const (
	head = `# Top Go Web Frameworks
A list of popular github projects related to Go web framework (ranked by stars automatically)
Please update **list.txt** (via Pull Request)

| Project Name | Stars | Forks | Open Issues | Description | Last Commit |
| ------------ | ----- | ----- | ----------- | ----------- | ----------- |
`
	tail = "\n*Last Automatic Update: %v*"

	warning = "⚠️ No longer maintained ⚠️  "
)

var (
	deprecatedRepos = [3]string{"https://github.com/go-martini/martini", "https://github.com/pilu/traffic"}
	repos           []Repo
)

func main() {
	accessToken := getAccessToken()

	byteContents, err := os.ReadFile("list.txt")
	if err != nil {
		log.Fatal(err)
	}

	lines := strings.Split(string(byteContents), "\n")
	for _, url := range lines {
		if strings.HasPrefix(url, "https://github.com/") {
			var repo Repo
			var commit HeadCommit

			repoAPI := fmt.Sprintf(
				"https://api.github.com/repos/%s",
				strings.TrimFunc(url[19:], trimSpaceAndSlash),
			)
			fmt.Println(repoAPI)

			statusCode, err := fetchJSON(accessToken, repoAPI, &repo)
			if err != nil {
				if statusCode == http.StatusNotFound {
					log.Printf("Skipping missing repository %s (%s): %v", url, repoAPI, err)
					continue
				}
				log.Fatalf("Failed to fetch repository %s (%s): %v", url, repoAPI, err)
			}

			commitAPI := fmt.Sprintf(
				"https://api.github.com/repos/%s/commits/%s",
				strings.TrimFunc(url[19:], trimSpaceAndSlash),
				repo.DefaultBranch,
			)
			fmt.Println(commitAPI)

			statusCode, err = fetchJSON(accessToken, commitAPI, &commit)
			if err != nil {
				if statusCode == http.StatusNotFound {
					log.Printf("Skipping repository with missing head commit %s (%s): %v", url, commitAPI, err)
					continue
				}
				log.Fatalf("Failed to fetch head commit for %s (%s): %v", url, commitAPI, err)
			}

			repo.LastCommitDate = commit.Commit.Committer.Date
			repos = append(repos, repo)

			fmt.Printf("Repository: %v\n", repo)
			fmt.Printf("Head Commit: %v\n", commit)

			time.Sleep(3 * time.Second)
		}
	}

	sort.Slice(repos, func(i, j int) bool {
		return repos[i].Stars > repos[j].Stars
	})
	saveRanking(repos)
}

func trimSpaceAndSlash(r rune) bool {
	return unicode.IsSpace(r) || (r == rune('/'))
}

func fetchJSON(accessToken, url string, target any) (int, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("authorization", fmt.Sprintf("Bearer %s", accessToken))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		bodyMessage := strings.TrimSpace(string(body))
		if bodyMessage == "" {
			return resp.StatusCode, fmt.Errorf("status %s", resp.Status)
		}
		return resp.StatusCode, fmt.Errorf("status %s: %s", resp.Status, bodyMessage)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	if err = decoder.Decode(target); err != nil {
		return resp.StatusCode, err
	}
	return resp.StatusCode, nil
}

func getAccessToken() string {
	tokenBytes, err := os.ReadFile("access_token.txt")
	if err != nil {
		log.Fatal("Error occurs when getting access token")
	}
	return strings.TrimSpace(string(tokenBytes))
}

func saveRanking(repos []Repo) {
	readme, err := os.OpenFile("README.md", os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer readme.Close()
	readme.WriteString(head)
	for _, repo := range repos {
		if isDeprecated(repo.URL) {
			repo.Description = warning + repo.Description
		}
		readme.WriteString(fmt.Sprintf("| [%s](%s) | %d | %d | %d | %s | %v |\n", repo.Name, repo.URL, repo.Stars, repo.Forks, repo.Issues, repo.Description, repo.LastCommitDate.Format("2006-01-02 15:04:05")))
	}
	readme.WriteString(fmt.Sprintf(tail, time.Now().Format(time.RFC3339)))
}

func isDeprecated(repoURL string) bool {
	for _, deprecatedRepo := range deprecatedRepos {
		if repoURL == deprecatedRepo {
			return true
		}
	}
	return false
}
