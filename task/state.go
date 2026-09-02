package task

var stateMachineMap = map[State][]State{
	Pending:   []State{Running, Scheduled, Failed},
	Scheduled: []State{Scheduled, Running, Failed},
	Running:   []State{Running, Completed, Failed},
	Completed: []State{},
	Failed:    []State{},
}

func Contains(states []State, state State) bool {
	for _, s := range states {
		if s == state {
			return true
		}
	}
	return false
}

func CanTransitTo(src State, dst State) bool {
	return Contains(stateMachineMap[src], dst)
}
