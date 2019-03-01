package machinery

import (
	"fmt"
	"time"

	machinery "github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
	"github.com/RichardKnop/machinery/v1/tasks"
)

var configPath = "config/config.yml"

// Builder ...
type Builder struct {
	server *machinery.Server
}

// NewBuilder ...
func NewBuilder() *Builder {
	b := new(Builder)
	b.server, _ = startServer()

	return b
}

func loadConfig() (*config.Config, error) {
	return config.NewFromYaml(configPath, true)
}

func startServer() (*machinery.Server, error) {
	cnf, err := loadConfig()
	if err != nil {
		return nil, err
	}

	// Create server instance
	server, err := machinery.NewServer(cnf)
	if err != nil {
		return nil, err
	}

	//Register tasks
	tasks := map[string]interface{}{
		"getTopStackOverFlowTags":    GetTopStackOverFlowTags,
		"getTopGitHubRepoByLanguage": GetTopGitHubRepoByLanguage,
		"printAllResults":            PrintAllResults,
	}

	err = server.RegisterTasks(tasks)
	if err != nil {
		return nil, err
	}

	return server, nil
}

func (b Builder) startWorker(errorsChan chan error) {
	consumerTag := "worker"
	worker := b.server.NewWorker(consumerTag, 10)
	worker.LaunchAsync(errorsChan)
}

func (b Builder) executeTaskToGetStackOverflowResults() ([]string, error) {
	var taskStackoverflow = tasks.Signature{
		Name: "getTopStackOverFlowTags",
	}

	asyncResult, err := b.server.SendTask(&taskStackoverflow)
	if err != nil {
		return nil, fmt.Errorf("Could not send task: %s", err.Error())
	}

	results, err := asyncResult.Get(time.Duration(time.Millisecond * 5))
	if err != nil {
		return nil, fmt.Errorf("Getting task result failed with error: %s", err.Error())
	}

	return results[0].Interface().([]string), nil
}

func (b Builder) processTasks() error {
	tags, err := b.executeTaskToGetStackOverflowResults()
	if err != nil {
		return err
	}
	var githubTasks = make([]*tasks.Signature, 0)
	for _, val := range tags {
		var taskGithub = tasks.Signature{
			Name: "getTopGitHubRepoByLanguage",
			Args: []tasks.Arg{
				{
					Type:  "string",
					Value: fmt.Sprintf("%v", val),
				},
			},
		}
		githubTasks = append(githubTasks, &taskGithub)
	}
	group, err := tasks.NewGroup(githubTasks...)
	if err != nil {
		return fmt.Errorf("Could not define new group: %s", err.Error())
	}
	var taskPrintAll = tasks.Signature{
		Name: "printAllResults",
	}
	chord, err := tasks.NewChord(group, &taskPrintAll)
	_, err = b.server.SendChord(chord, 0)
	if err != nil {
		return fmt.Errorf("Could not send chord: %s", err.Error())
	}
	return nil
}

// Do ...
func (b Builder) Do() error {
	errorsChan := make(chan error)
	b.startWorker(errorsChan)
	err := b.processTasks()
	if err != nil {
		return err
	}
	return <-errorsChan
}
