package machinery

import (
	"fmt"
	"time"

	machinery "github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
	"github.com/RichardKnop/machinery/v1/log"
	"github.com/RichardKnop/machinery/v1/tasks"
)

var configPath = "config/config.yml"

// builder ...
type builder struct {
	server *machinery.Server
}

// NewBuilder ...
func NewBuilder() *builder {
	b := new(builder)
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
	tasks := map[string]interface{}{}

	err = server.RegisterTasks(tasks)
	if err != nil {
		return nil, err
	}

	return server, nil
}

func (b builder) startWorker(errorsChan chan error) {
	fmt.Println("startWorker")
	consumerTag := "worker"
	worker := b.server.NewWorker(consumerTag, 0)
	worker.LaunchAsync(errorsChan)
}

func (b builder) processTasks() error {
	fmt.Println("processTasks")
	// var taskGithub = tasks.Signature{
	// 	Name: "getTopGitHubRepoByLanguage",
	// 	Args: []tasks.Arg{
	// 		{
	// 			Type:  "string",
	// 			Value: "javascript",
	// 		},
	// 	},
	// }

	var taskStackoverflow = tasks.Signature{
		Name: "getTopStackOverFlowTags",
	}

	asyncResult, err := b.server.SendTask(&taskStackoverflow)
	if err != nil {
		return fmt.Errorf("Could not send task: %s", err.Error())
	}

	results, err := asyncResult.Get(time.Duration(time.Millisecond * 5))
	if err != nil {
		return fmt.Errorf("Getting task result failed with error: %s", err.Error())
	}
	log.INFO.Printf(" = %v\n", tasks.HumanReadableResults(results))

	return nil
}

// Do ...
func (b builder) Do() error {
	errorsChan := make(chan error)
	b.startWorker(errorsChan)

	err := b.processTasks()
	if err != nil {
		return err
	}
	return <-errorsChan
}
