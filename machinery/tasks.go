package machinery

import (
	"encoding/json"
	"fmt"

	"github.com/evandroferreiras/machinery-tutorial/api"
)

// GitHubResponse ...
type GitHubResponse struct {
	Language     string   `json:"language,omitempty"`
	Repositories []string `json:"repositories,omitempty"`
}

// GetTopStackOverFlowTags ...
func GetTopStackOverFlowTags() ([]string, error) {
	return api.GetTopStackOverFlowTags()
}

// GetTopGitHubRepoByLanguage ...
func GetTopGitHubRepoByLanguage(language string) (string, error) {
	repositories, err := api.GetTopRepoByLanguage(language, 10)
	if err != nil {
		return "", err
	}
	marshalled, err := json.Marshal(GitHubResponse{language, repositories})
	return string(marshalled), err
}

// PrintAllResults ...
func PrintAllResults(args ... string) error {
	fmt.Println("-RELATORIO--------------------------")
	
	for _, r := range args {
		fmt.Println(r)
	}
	fmt.Println("-----------------------------------")
	return nil
}
