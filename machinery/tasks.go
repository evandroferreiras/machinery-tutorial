package machinery

import (
	"encoding/json"

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
