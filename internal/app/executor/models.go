package executor

type ExecutorError struct {
	t      string
	er     error
	bundle string
}

type databaseCh chan interface{}
type errorCh chan ExecutorError
