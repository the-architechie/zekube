package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
	"zekube/task"
	"zekube/worker"

	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
)

func main() {

	host := os.Getenv("ZEKUBE_HOST")
	port, _ := strconv.Atoi(os.Getenv("ZEKUBE_PORT"))

	fmt.Println("Starting the Zekube worker")
	db := make(map[uuid.UUID]*task.Task)

	// create a worker
	w := worker.Worker{
		Queue: *queue.New(),
		Db:    db,
	}
	api := worker.Api{Address: host, Port: port, Worker: &w}
	fmt.Println("Starting a task")
	go runTasks(&w)
	go w.CollectStats()
	api.Start()
}

func runTasks(w *worker.Worker) {
	for {
		if w.Queue.Len() != 0 {
			result := w.RunTask()
			if result.Error != nil {
				log.Printf("Error running task: %v\n", result.Error)
			}
		} else {
			log.Printf("No tasks to process currently \n")
		}
		log.Printf("Sleeping for 10 seconds... \n")
		time.Sleep(10 * time.Second)
	}
}
