package executor

type ExecutorError struct {
	T      string `json:"t,omitempty"`
	Er     string `json:"er,omitempty"`
	Bundle string `json:"bundle,omitempty"`
}

type databaseCh chan interface{}
