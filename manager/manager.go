package manager

import (
	"cube-orchestrator/task"
	"fmt"

	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
)

// Manager will need to keep track of the workers in the cluster
type Manager struct {
	Pending        queue.Queue
	TaskDb         map[string][]*task.Task      // this store the tasks
	EventDB        map[string][]*task.TaskEvent // stores events
	Workers        []string
	WorkersTaskMap map[string][]uuid.UUID // WorkersTaskMap maps worker IDs to the list of task UUIDs assigned to them.
	TaskWorkerMap  map[uuid.UUID]string   // TaskWorkerMap maps task UUIDs to the IDs of the workers responsible for executing them.
}

func (m *Manager) SelectWorker() {
	fmt.Println("Select an appropriate worker")
}

func (m *Manager) UpdateTasks() {
	fmt.Println("Update the tasks")
}
func (m *Manager) SendWork() {
	fmt.Println("Send Work")
}
