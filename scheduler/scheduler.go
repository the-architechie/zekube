package scheduler

type Schecuker interface {
	SelectCandidateNodes()
	Score()
	Pick()
}
