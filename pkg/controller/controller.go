package controller

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"tinygoexercise/pkg/logger"

	"go.uber.org/zap"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type Interface interface {
	New(port string) *Controller
	ControllContainer(c *Controller)
}

type Controller struct {
	statusStateList *[]types.Container
	sugarLogger     *zap.SugaredLogger
	apiEndpoint     string
	logger          *zap.SugaredLogger
}

var docker *client.Client

const (
	stoppedContainerState string = "exited"
)

// New creates a new Controller with the given port number
func New(port string) *Controller {
	c := Controller{
		logger:      logger.New(),
		apiEndpoint: "http://localhost:" + port + "/watchContainer",
	}

	defer func() {
		err := c.logger.Sync()
		if err != nil {
			c.logger.Error(err)
		}
	}()

	docker = newDockerClient(&c)

	c.statusStateList = getContainerList(&c)
	return &c
}

// httpGet returns the response from the api endpoint
func httpGet(c *Controller, url string) (*http.Response, error) {
	c.logger.Infof("Fetching response from %s...", c.apiEndpoint)
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	c.logger.Infof("Fetching response from %s... Done", c.apiEndpoint)
	return response, nil
}

// decodeResponse decodes JSON response to a list of containers
func decodeResponse(c *Controller, r *http.Response) ([]types.Container, error) {
	c.logger.Info("Decoding response to type Container...")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	var result []types.Container
	if err := json.Unmarshal(body, &result); err != nil {
		c.logger.Error(err)
	}

	c.logger.Info("Decoding response to type Container... Done")
	return result, nil
}

// getContainerList returns the pointer to the list of containers from api endpoint
func getContainerList(c *Controller) *[]types.Container {
	response, err := httpGet(c, c.apiEndpoint)
	if err != nil {
		c.logger.Error(err)
	}
	defer response.Body.Close()

	result, err := decodeResponse(c, response)
	if err != nil {
		c.logger.Error(err)
	}
	return &result
}

// ControllContainer restarts stopped containers in c.statusStateList
// stops and removes unknown containers
func ControllContainer(c *Controller) {
	clist := *getContainerList(c)
	for _, container := range clist {
		if find(c.statusStateList, &container) {
			if container.State == stoppedContainerState {
				restartContainer(c, &container)
			}
		} else {
			c.logger.Infof("Found an unknown container %s", container.Names[0])
			stopContainer(c, &container)
			deleteContainer(c, &container)
		}
	}
}

// TODO: refactor the following functions
// restartContainer restarts a given container
func restartContainer(c *Controller, container *types.Container) {
	c.logger.Infof("Restarting container %s...", container.Names[0])
	if err := docker.ContainerRestart(context.Background(), container.ID, nil); err != nil {
		c.logger.Error(err)
	}
	c.logger.Infof("Container %s restarted.", container.Names[0])
}

// stopContainer stops a given container
func stopContainer(c *Controller, container *types.Container) {
	c.logger.Infof("Stopping container %s...", container.Names[0])
	if err := docker.ContainerStop(context.Background(), container.ID, nil); err != nil {
		c.logger.Error(err)
	}
	c.logger.Infof("Container %s stopped.", container.Names[0])
}

// deleteContainer removes a given container
func deleteContainer(c *Controller, container *types.Container) {
	c.logger.Infof("Deleting container %s...", container.Names[0])
	options := types.ContainerRemoveOptions{
		RemoveVolumes: true,
		RemoveLinks:   false,
	}
	if err := docker.ContainerRemove(context.Background(), container.ID, options); err != nil {
		c.logger.Error(err)
	}
	c.logger.Infof("Container %s deleted.", container.Names[0])
}

// find checks the given item is in the given list
func find(slice *[]types.Container, sliceItem *types.Container) bool {
	for _, item := range *slice {
		if item.ID == sliceItem.ID {
			return true
		}
	}
	return false
}

// newDockerClient initilizes a Docker Client
func newDockerClient(c *Controller) *client.Client {
	c.logger.Info("Initializating docker client...")
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		c.logger.Error(err)
	}
	c.logger.Info("Docker client initializated.")
	return cli
}
