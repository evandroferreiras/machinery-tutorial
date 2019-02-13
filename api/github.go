package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var gitHubURL = "https://api.github.com/search/repositories?q=language:%s&sort=stars&order=desc"

type repoItems struct {
	FullName        string `json:"full_name,omitempty"`
	Name            string `json:"name,omitempty"`
	HTMLURL         string `json:"html_url,omitempty"`
	StargazersCount int    `json:"stargazers_count,omitempty"`
	WatchersCount   int    `json:"watchers_count,omitempty"`
	Forks           int    `json:"forks,omitempty"`
}

type repo struct {
	Items []repoItems
}

// GetTopRepoByLanguage ...
func GetTopRepoByLanguage(language string) ([]string, error) {
	url := fmt.Sprintf(gitHubURL, language)
	fmt.Println(url)
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	data, _ := ioutil.ReadAll(response.Body)
	var dat repo
	if err := json.Unmarshal(data, &dat); err != nil {
		return nil, err
	}

	result := make([]string, 0)
	if len(dat.Items) > 0 {
		for _, value := range dat.Items[:5] {
			result = append(result, value.Name)
		}
	}

	return result, nil

}
