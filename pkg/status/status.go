package status

// Status presents for the status of server
type Status struct {
	Connections   int // Number of current connections
	WorkerSize    int // Number of current workers
	WorkQueueSize int // Work Queue Size
	Works         int // Number of current works in the queue
}
