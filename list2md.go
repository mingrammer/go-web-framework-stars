package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
	"unicode"
)

// Repo describes a Github repository
type Repo struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Stars       int       `json:"stargazers_count"`
	Forks       int       `json:"forks_count"`
	Created     time.Time `json:"created_at"`
	Updated     time.Time `json:"updated_at"`
	URL         string    `json:"html_url"`
}

const (
	head = `# Top Go Web Frameworks
A list of popular github projects related to Go web framework (ranked by stars automatically)
Please update **list.txt** (via Pull Request)

| Project Name | Stars | Forks | Description |
| ------------ | ----- | ----- | ----------- |
`
	tail = "\n*Last Automatic Update: %v*"
)

var result []Repo

func main() {
	accessToken := getAccessToken()

	byteContents, err := ioutil.ReadFile("list.txt")
	if err != nil {
		log.Fatal(err)
	}

	lines := strings.Split(string(byteContents), "\n")
	for _, url := range lines {
		if strings.HasPrefix(url, "https://github.com/") {
			var repo Repo

			apiURL := fmt.Sprintf("https://api.github.com/repos/%s?access_token=%s", strings.TrimFunc(url[19:], trimSpaceAndSlash), accessToken)
			fmt.Println(apiURL)

			resp, err := http.Get(apiURL)
			if err != nil {
				log.Fatal(err)
			}

			decoder := json.NewDecoder(resp.Body)
			if err = decoder.Decode(&repo); err != nil {
				log.Fatal(err)
			}

			if repo.Name != "" {
				result = append(result, repo)
			}
			fmt.Printf("%v\n", repo)
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Stars > result[j].Stars
	})
	saveRanking(result)
}

func trimSpaceAndSlash(r rune) bool {
	return unicode.IsSpace(r) || (r == rune('/'))
}

func getAccessToken() string {
	tokenBytes, err := ioutil.ReadFile("access_token.txt")
	if err != nil {
		log.Fatal("Error occurs when getting access token")
	}
	return string(tokenBytes)
}

func saveRanking(result []Repo) {
	readme, err := os.OpenFile("README.md", os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatal(err)
	}
	readme.WriteString(head)
	for _, repo := range result {
		readme.WriteString(fmt.Sprintf("| [%s](%s) | %d | %d | %s |\n", repo.Name, repo.URL, repo.Stars, repo.Forks, repo.Description))
	}
	readme.WriteString(fmt.Sprintf(tail, time.Now().Format("2006-01-02 15:04:05")))
}
