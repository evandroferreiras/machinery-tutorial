package tasks

import "github.com/evandroferreiras/machinery-tutorial/api"

// GetTopStackOverFlowTags ...
func GetTopStackOverFlowTags() ([]string, error) {
	return api.GetTopStackOverFlowTags()
}

// GetTopGitHubRepoByLanguage ...
func GetTopGitHubRepoByLanguage(language string) ([]string, error) {
	return api.GetTopRepoByLanguage(language)
}
