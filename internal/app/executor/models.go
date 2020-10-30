package executor

type ExecutorError struct {
	T      string `json:"t,omitempty"`
	Er     error  `json:"er,omitempty"`
	Bundle string `json:"bundle,omitempty"`
}

type databaseCh chan interface{}
type errorCh chan ExecutorError
