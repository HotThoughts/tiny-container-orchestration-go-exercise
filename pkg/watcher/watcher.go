package watcher

import (
	"context"
	"encoding/json"
	"net/http"
	"tinygoexercise/pkg/logger"

	"go.uber.org/zap"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type Interface interface {
	InitWatcher(port string)
}

type Watcher struct {
	logger *zap.SugaredLogger
}

var (
	docker *client.Client
	w      Watcher
)

// InitWatcher initilizes a Watcher
func InitWatcher(port string) {
	w = Watcher{logger: logger.New()}
	defer func() {
		err := w.logger.Sync()
		if err != nil {
			w.logger.Error(err)
		}
	}()

	docker = newDockerClient()

	http.HandleFunc("/", watchContainer)

	w.logger.Infof("Container Watcher server stated at localhost:%s", port)
	w.logger.Fatal(http.ListenAndServe(":"+port, nil))
}

func newDockerClient() *client.Client {
	w.logger.Info("Initializating docker client...")
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		w.logger.Error(err)
	}
	w.logger.Info("Docker client initializated.")
	return cli
}

// watchContainer creates a api endpoint for list of containers running on the local docekr host
func watchContainer(rw http.ResponseWriter, _ *http.Request) {
	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Type", "application/json")

	w.logger.Info("Getting containers info... ")
	containers, err := docker.ContainerList(context.Background(), types.ContainerListOptions{
		All: true,
	})
	if err != nil {
		w.logger.Panic(err)
	}
	w.logger.Info("Getting containers info... Done")

	containersJson, err := json.Marshal(containers)
	if err != nil {
		w.logger.Panic(err)
	}
	w.logger.Info("Writing response... ")

	_, err = rw.Write(containersJson)
	if err != nil {
		w.logger.Panic(err)
	}
	w.logger.Info("Writing response... Done.")
}
