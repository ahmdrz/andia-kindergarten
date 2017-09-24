package engine

const (
	STATE_READY = iota
	STATE_WAITING
	STATE_LISTENING
)

type Engine struct {
	State   int
	Caption string
}
