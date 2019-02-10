package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sort"
)

var stackOverflowURL = "https://api.stackexchange.com/2.2/tags?order=desc&sort=popular&site=stackoverflow"

type tagItems struct {
	Count int    `json:"count"`
	Name  string `json:"name"`
}

type tags struct {
	Items []tagItems
}

func GetTopStackOverFlowTags() ([]string, error) {
	response, err := http.Get(stackOverflowURL)
	if err != nil {
		return nil, err
	}

	data, _ := ioutil.ReadAll(response.Body)
	var dat tags
	if err := json.Unmarshal(data, &dat); err != nil {
		return nil, err
	}

	items := dat.Items
	sort.Slice(items, func(i, j int) bool {
		return items[i].Count > items[j].Count
	})

	result := make([]string, 0)
	for _, value := range dat.Items[:5] {
		result = append(result, value.Name)
	}

	return result, nil
}
