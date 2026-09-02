package task

import (
	"context"
	"io"
	"log"
	"math"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/moby/moby/api/pkg/stdcopy"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/network"
	"github.com/moby/moby/client"
)

type State int

const (
	Pending State = iota
	Scheduled
	Running
	Completed
	Failed
)

type Task struct {
	ID            uuid.UUID
	ContainerID   string
	Name          string
	State         State
	Image         string
	Cpu           float64
	Memory        int
	Disk          int
	ExposedPorts  network.PortSet
	PortBindings  map[string]string
	RestartPolicy container.RestartPolicyMode
	StartTime     time.Time
	FinishTime    time.Time
}

type Event struct {
	ID        uuid.UUID
	State     State
	Timestamp time.Time
	Task      Task
}

// Config Task configuration
type Config struct {
	Name          string // name of the container which is also the name of the task
	AttachStdin   bool
	AttachStdout  bool
	AttachStdErr  bool
	ExposedPorts  network.PortSet
	Cmd           []string
	Image         string // image of the container will run
	Cpu           float64
	Memory        int64
	Disk          int64
	Env           []string                    // variables to be passed into the container
	RestartPolicy container.RestartPolicyMode // tells docker daemon what to do if the container dies unexpectedly
}

func NewConfig(t *Task) *Config {
	return &Config{
		Name:          t.Name,
		ExposedPorts:  t.ExposedPorts,
		Image:         t.Image,
		Cpu:           t.Cpu,
		Memory:        int64(t.Memory),
		Disk:          int64(t.Disk),
		RestartPolicy: t.RestartPolicy,
	}
}

// Docker - encapsulates everything we need to run a task as a docker container
type Docker struct {
	Client *client.Client // Docker client
	Config Config         // task configuration
}

func NewDocker(c *Config) *Docker {
	dc, _ := client.New(client.FromEnv)
	return &Docker{
		Client: dc,
		Config: *c,
	}
}

type DockerResult struct {
	Error       error
	Action      string
	ContainerID string
	Result      string
}

func (d *Docker) Run() DockerResult {
	ctx := context.Background()

	reader, err := d.Client.ImagePull(ctx, d.Config.Image, client.ImagePullOptions{})
	if err != nil {
		log.Printf("Error pulling image %s: %v\n", d.Config.Image, err)
		return DockerResult{Error: err}
	}
	_, err = io.Copy(os.Stdout, reader)
	if err != nil {
		return DockerResult{
			Error: err,
		}
	}
	rp := container.RestartPolicy{
		Name: d.Config.RestartPolicy,
	}

	r := container.Resources{
		Memory:   d.Config.Memory,
		NanoCPUs: int64(d.Config.Cpu * math.Pow(10, 9)),
	}

	cc := container.Config{
		Image:        d.Config.Image,
		Tty:          false,
		Env:          d.Config.Env,
		ExposedPorts: d.Config.ExposedPorts,
	}

	hc := container.HostConfig{
		RestartPolicy:   rp,
		Resources:       r,
		PublishAllPorts: true,
	}

	ccOptions := client.ContainerCreateOptions{
		Config:           &cc,
		HostConfig:       &hc,
		NetworkingConfig: nil,
		Platform:         nil,
		Name:             d.Config.Name,
	}
	resp, err := d.Client.ContainerCreate(ctx, ccOptions)
	if err != nil {
		log.Printf("Error creating container using image %s: %v\n", d.Config.Image, err)
		return DockerResult{Error: err}
	}

	_, err = d.Client.ContainerStart(ctx, resp.ID, client.ContainerStartOptions{})
	if err != nil {
		log.Printf("Error startign container %s: %v\n", resp.ID, err)
		return DockerResult{Error: err}
	}
	out, err := d.Client.ContainerLogs(
		ctx, resp.ID, client.ContainerLogsOptions{ShowStdout: true, ShowStderr: true},
	)
	if err != nil {
		log.Printf("Error retrieving logs  for container %s: %v\n", resp.ID, err)
		return DockerResult{Error: err}
	}
	stdcopy.StdCopy(os.Stdout, os.Stderr, out)
	return DockerResult{
		ContainerID: resp.ID,
		Action:      "start",
		Result:      "success",
	}
}

func (d *Docker) Stop(id string) DockerResult {
	log.Printf("Attempting to stop container %v", id)
	ctx := context.Background()
	_, err := d.Client.ContainerStop(ctx, id, client.ContainerStopOptions{})

	if err != nil {
		log.Printf("Failed to stop container %s: %v", id, err)
		return DockerResult{Error: err}
	}

	_, err = d.Client.ContainerRemove(ctx, id, client.ContainerRemoveOptions{})
	if err != nil {
		log.Printf("Failed to remove container %s: %v", id, err)
		return DockerResult{Error: err}
	}
	log.Printf("Stopped and removed container %v", id)

	return DockerResult{
		ContainerID: id,
		Action:      "stop",
		Result:      "success",
	}
}
