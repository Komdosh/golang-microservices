package domain

type CalculatedBoxFitness struct {
	Box     Box
	Fitness int
}

type BestBoxResult struct {
	Box   *Box
	Error *Error
}

type Error struct {
	Message string
}
