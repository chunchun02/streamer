package server

import (
	"github.com/chunchun02/streamer/internal/worker"
	"github.com/chunchun02/streamer/pkg/connection"
	"github.com/chunchun02/streamer/pkg/status"
	"github.com/chunchun02/streamer/pkg/work"
	"github.com/google/uuid"
	"sync"
)

// Streamer presents for the support methods by server
type Streamer interface {
	Subscribe() *connection.Connection             // Register a new connection
	Unsubscribe(connection *connection.Connection) // Disconnect a connection
	In(packet interface{})                         // Register a new worker
	Start()                                        // Start the server
	Status() status.Status                         // Return the current server status
	Stop()                                         // Stop the server
}

// WorkerChannel to communicate between server and workers
var workerChannel = make(chan chan *work.Work)

// Server presents for a application server
type Server struct {
	mu         sync.Mutex             // Mutex to multi-threading
	works      chan *work.Work        // Work queue to register a new work into
	workers    []*worker.Worker       // List of server workers
	connectors *connection.Connectors // List of current connecting users
	workerSize int                    // Worker size
	end        chan bool              // End channel to send the stop sign
}

// NewsServer to create a new server base on worker size and work queue size
// Please carefully to choose a suitable worker size, base on your application
// If the application needs lots of CPU resource, reduce the worker size
// If the application doesn't need lots of CPU resource and I/O resource instead, increase the worker size to increase the performance
func NewServer(workerSize int, workQueueSize int) *Server {
	return &Server{
		works:      make(chan *work.Work, workQueueSize),
		end:        make(chan bool),
		connectors: &connection.Connectors{},
		workerSize: workerSize,
	}
}

// Subscribe to register a new connection
func (server *Server) Subscribe() *connection.Connection {
	// Lock the server
	server.mu.Lock()
	defer server.mu.Unlock()
	return server.connectors.RegisterConnection()
}

// Unsubscribe to disconnect the connection
func (server *Server) Unsubscribe(connection *connection.Connection) {
	// Lock the server
	server.mu.Lock()
	defer server.mu.Unlock()
	server.connectors.UnregisterConnection(connection)
}

// In to register a job
func (server *Server) In(packet interface{}) {
	// Generate a new unique UUID
	u, _ := uuid.NewRandom()
	// Create a new Work for a new job and push it to the work queue
	server.works <- &work.Work{ID: u, Data: packet}
}

// Start starts the server
func (server *Server) Start() {
	var i int
	// Loop and create all initial workers
	for i < server.workerSize {
		i++
		// Create a new worker
		w := worker.NewWorker(workerChannel)
		// Start the worker after creating
		w.Start(server.connectors)
		// Register the worker to server
		server.workers = append(server.workers, w)
	}

	// Start handling the jobs
	go func() {
		for {
			select {
			case <-server.end:
				// If the server received the stop sign, loop ad stop all workers
				for _, w := range server.workers {
					w.Stop()
				}
				return
			case wo := <-server.works:
				// Watch to see if there is any available worker and let the worker do the job
				wkr := <-workerChannel
				wkr <- wo
			}
		}
	}()
}

// Status to return te current status of server
func (server *Server) Status() status.Status {
	return status.Status{
		Connections:   len(server.connectors.Connections),
		WorkerSize:    server.workerSize,
		WorkQueueSize: cap(server.works),
		Works:         len(server.works),
	}
}

// Stop stops the server by sending the stop sign to End channel
func (server *Server) Stop() {
	server.end <- true
}
