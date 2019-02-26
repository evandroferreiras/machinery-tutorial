package machinery

import (
	"context"
	"fmt"
	"time"

	machinery "github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
	"github.com/RichardKnop/machinery/v1/log"
	"github.com/RichardKnop/machinery/v1/tasks"
	"github.com/google/uuid"
	opentracing "github.com/opentracing/opentracing-go"
	opentracing_log "github.com/opentracing/opentracing-go/log"
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
	tasks := map[string]interface{}{
		"getTopStackOverFlowTags":    GetTopStackOverFlowTags,
		"getTopGitHubRepoByLanguage": GetTopGitHubRepoByLanguage,
	}

	err = server.RegisterTasks(tasks)
	if err != nil {
		return nil, err
	}

	return server, nil
}

func (b builder) startWorker(errorsChan chan error) {
	consumerTag := "worker"
	worker := b.server.NewWorker(consumerTag, 10)
	worker.LaunchAsync(errorsChan)
}

func (b builder) processTasks() error {
	var taskStackoverflow = tasks.Signature{
		Name: "getTopStackOverFlowTags",
	}

	span, ctx := opentracing.StartSpanFromContext(context.Background(), "send")
	defer span.Finish()

	batchID := uuid.New().String()
	span.SetBaggageItem("batch.id", batchID)
	span.LogFields(opentracing_log.String("batch.id", batchID))

	log.INFO.Println("Starting batch:", batchID)

	asyncResult, err := b.server.SendTaskWithContext(ctx, &taskStackoverflow)
	if err != nil {
		return fmt.Errorf("Could not send task: %s", err.Error())
	}

	results, err := asyncResult.Get(time.Duration(time.Millisecond * 5))
	if err != nil {
		return fmt.Errorf("Getting task result failed with error: %s", err.Error())
	}

	for _, r := range results {
		log.INFO.Println(r.Interface())
		i := r.Interface()
		a := i.([]string)

		var t = make([]*tasks.Signature, 0)

		for _, val := range a {
			var taskGithub = tasks.Signature{
				Name: "getTopGitHubRepoByLanguage",
				Args: []tasks.Arg{
					{
						Type:  "string",
						Value: fmt.Sprintf("%v", val),
					},
				},
			}

			t = append(t, &taskGithub)
		}
		group, _ := tasks.NewGroup(t...)
		asyncResults, err := b.server.SendGroup(group, 0)
		if err != nil {
			return fmt.Errorf("Getting group result failed with error: %s", err.Error())
		}
		for _, asyncResult := range asyncResults {
			results, err := asyncResult.Get(time.Duration(time.Millisecond * 5))
			if err != nil {
				return fmt.Errorf("Getting group AsyncResult failed with error: %s", err.Error())
			}
			for _, result := range results {
				fmt.Println(result.Interface())
			}
		}
	}

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
