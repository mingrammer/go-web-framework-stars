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
	deprecatedRepos = [2]string{"https://github.com/go-martini/martini", "https://github.com/pilu/traffic"}
	repos           []Repo
)


func main() {
	accessToken := getAccessToken()
    //fmt.Println("accessToken", accessToken)
    //return
    //var theList map[string]interface{} = map[string]interface{}{}
    theList := map[string]interface{}{}
	byteContents, err := ioutil.ReadFile("list.txt")
	if err != nil {
		log.Fatal(err)
    }

    //lines, err := getTodoLines()
	lines := strings.Split(string(byteContents), "\n")
    fmt.Println("got", len(lines), err)
	for _, url := range lines {
		if strings.HasPrefix(url, "https://github.com/") {
			var repo Repo
			var commit HeadCommit
			//time.Sleep(2 * time.Second)
            //repoAPI := fmt.Sprintf("https://api.github.com/repos/%s?access_token=%s", strings.TrimFunc(url[19:], trimSpaceAndSlash), accessToken)
			repoAPI := fmt.Sprintf("https://api.github.com/repos/%s", strings.TrimFunc(url[19:], trimSpaceAndSlash))
			fmt.Println(repoAPI)

			fmt.Println(fmt.Sprintf("token %s", accessToken))
            //return
			//resp, err := http.Get(repoAPI)
			//resp.Header.Add("Authorization",fmt.Sprintf("token %s", accessToken))
            //resp.Header.Add("User-Agent","request")
            client := &http.Client{}
            req , err := http.NewRequest("GET", repoAPI, nil)
			if err != nil {
				log.Fatal(err)
			}
            req.Header.Set("Authorization", fmt.Sprintf("token %s", accessToken))
            //req.Header.Set("User-Agent","request")

            resp, err := client.Do(req)

			if err != nil {
				log.Fatal(err)
			}
			if resp.StatusCode != 200 {
				log.Fatal(resp.Status)
			}

			decoder := json.NewDecoder(resp.Body)
			if err = decoder.Decode(&repo); err != nil {
				log.Fatal(err)
			}

			//commitAPI := fmt.Sprintf("https://api.github.com/repos/%s/commits/%s?access_token=%s", strings.TrimFunc(url[19:], trimSpaceAndSlash), repo.DefaultBranch, accessToken)
			commitAPI := fmt.Sprintf("https://api.github.com/repos/%s/commits/%s", strings.TrimFunc(url[19:], trimSpaceAndSlash), repo.DefaultBranch)
			fmt.Println(commitAPI)

			//resp, err = http.Get(commitAPI)
            req , err = http.NewRequest("GET", commitAPI, nil)
			if err != nil {
				log.Fatal(err)
			}
            req.Header.Set("Authorization", fmt.Sprintf("token %s", accessToken))
			//resp.Header.Add("Authorization",fmt.Sprintf("token %s", accessToken))
			//resp.Header.Add("User-Agent","request")
            resp, err = client.Do(req)
			if err != nil {
				log.Fatal(err)
			}
			if resp.StatusCode != 200 {
				log.Fatal(resp.Status)
			}

			decoder = json.NewDecoder(resp.Body)
			if err = decoder.Decode(&commit); err != nil {
				log.Fatal(err)
			}

			repo.LastCommitDate = commit.Commit.Committer.Date
			repos = append(repos, repo)

			fmt.Printf("Repository: %v\n", repo)
			fmt.Printf("Head Commit: %v\n", commit)

            theList[url] = repo
		}
	}

	sort.Slice(repos, func(i, j int) bool {
		return repos[i].Stars > repos[j].Stars
	})
	saveRanking(repos)
}

func getTodoLines() (todoLines []string, err error) {
    defer func(){
        if err := recover(); err != nil {
            fmt.Println("Recovered in f", err)
        }
    }()
    currentTime := time.Now()
    fmt.Println("\n######################################\n")
    fmt.Println(currentTime.Format("2006-01-02 15:04:05"))

	byteContents, err := ioutil.ReadFile("list.txt")
	if err != nil {
		log.Fatal(err)
	}

	lines := strings.Split(string(byteContents), "\n")
    startingHour := 7
    maxLines := 1//5//15
    currentHour := currentTime.Hour()
    fmt.Println(fmt.Sprintf("currentHour %d",currentHour))
    if(currentHour<startingHour){
        //panic(nil)
        fmt.Println("Bye")
        return
    }
    offsetHour := currentHour - startingHour
    fmt.Println(offsetHour)
    //return lines//[100:1]
    begin := 0
    fmt.Println("offsetHour", offsetHour)
    begin = offsetHour * maxLines
    if begin > len(lines){
        begin = len(lines) - 1
    }
    end := 0
    end = begin + maxLines
    if end > len(lines){
        end = len(lines) - 1
    }
    fmt.Println("begin", begin, "end", end)

    if end == begin {
        //end file
    }
    todoLines = lines[begin:end]
    /*fileName := "userfile.json"
    content, err := ioutil.ReadFile(fileName)
    
    var theList map[string]interface{}
    //theList := []interface{}{}

    if err != nil {
		//log.Fatal(err)
        file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0755)
        if err != nil {
		    log.Fatal(err)
	    }
        file.Close()
    } else {
        json.Unmarshal(content, &theList)
        theList["good"] = map[string]interface{}{
            "name":"Alfred",
            "age":"40",
        }

    }
    fmt.Println("content", content)

    file2, _ := json.Marshal(theList) //, "", " ")
    err = ioutil.WriteFile("test.json", file2, 0644)
    if err != nil {
		log.Fatal(err)
	}*/
    return
}
func trimSpaceAndSlash(r rune) bool {
	return unicode.IsSpace(r) || (r == rune('/'))
}

func getAccessToken() string {
	tokenBytes, err := ioutil.ReadFile("access_token.txt")
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
