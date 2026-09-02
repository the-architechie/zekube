package worker

import (
	"errors"
	"fmt"
	"log"
	"time"
	"zekube/task"

	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
)

type Worker struct {
	Name      string
	Queue     queue.Queue
	Db        map[uuid.UUID]*task.Task // in-memory data store TODO: should evolve to a full fledged db
	TaskCount int
	Stats     *Stats
}

func (w *Worker) AddTask(t task.Task) {
	w.Queue.Enqueue(t)
}

func (w *Worker) RunTask() task.DockerResult {
	// pop a task off the worker's queue
	t := w.Queue.Dequeue()
	// if there no tasks in the queue, return an error
	if t == nil {
		log.Println("No tasks in the queue")
		return task.DockerResult{Error: nil}
	}

	//COnvert the item in the queue to a Task type
	taskQueued := t.(task.Task)
	// Retrieve the same task from the db
	taskPersisted := w.Db[taskQueued.ID]
	// if the task is not persisted
	// persist it
	if taskPersisted == nil {
		// persist the task
		taskPersisted = &taskQueued
		w.Db[taskQueued.ID] = &taskQueued

	}
	var result task.DockerResult
	// check the state machine
	if task.CanTransitTo(taskPersisted.State, taskQueued.State) {
		switch taskQueued.State {
		case task.Scheduled:
			result = w.StartTask(taskQueued)
		case task.Completed:
			result = w.StopTask(taskQueued)
		default:
			result.Error = errors.New("illegal state transition")
		}

	} else {
		err := fmt.Errorf("invalid state transition from %v to %v", taskPersisted.State, taskQueued.State)
		result.Error = err
	}

	return result
}

func (w *Worker) StartTask(t task.Task) task.DockerResult {
	t.StartTime = time.Now().UTC()
	config := task.NewConfig(&t)
	d := task.NewDocker(config)
	result := d.Run()
	if result.Error != nil {
		log.Printf("Err running task %v: %v\n", t.ID, result.Error)
		t.State = task.Failed
		w.Db[t.ID] = &t
		return result
	}
	t.ContainerID = result.ContainerID
	t.State = task.Running
	w.Db[t.ID] = &t
	return result
}

func (w *Worker) StopTask(t task.Task) task.DockerResult {
	config := task.NewConfig(&t)
	d := task.NewDocker(config)
	result := d.Stop(t.ContainerID)
	if result.Error != nil {
		log.Printf("Error removing container %v for task %v\n", t.ContainerID, t.ID)
	}
	return result
}

func (w *Worker) GetTasks() []*task.Task {
	var tasks []*task.Task
	for _, t := range w.Db {
		tasks = append(tasks, t)
	}
	return tasks
}

func (w *Worker) CollectStats() {
	for {
		log.Println("Collecting stats")
		w.Stats = GetStats()
		w.Stats.TaskCount = w.TaskCount
		time.Sleep(15 * time.Second)
	}
}
