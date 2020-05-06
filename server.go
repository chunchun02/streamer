package streamer

import (
	"github.com/chunchun02/streamer/internal/connection"
	"github.com/chunchun02/streamer/internal/status"
	"github.com/chunchun02/streamer/internal/work"
	"github.com/chunchun02/streamer/internal/worker"
	"github.com/google/uuid"
	"sync"
)

// Streamer presents for the support methods by server
type Streamer interface {
	Subscribe() *connection.Connection             // Register a new connection
	Unsubscribe(connection *connection.Connection) // Disconnect a connection
	In(packet interface{})                         // Register a new wor
	Start()                                        // Start the server
	Status() status.Status                         // Return the current server status
	Stop()                                         // Stop the server
}

// WorkerChannel to communicate between server and workers
var WorkerChannel = make(chan chan *work.Work)

// Server presents for a application server
type Server struct {
	mu         sync.Mutex             // Mutex to multi-threading
	works      chan *work.Work        // Work queue to register a new work into
	workers    []*worker.Worker       // List of server workers
	End        chan bool              // End channel to send the stop sign
	Connectors *connection.Connectors // List of current connecting users
	WorkerSize int                    // Worker size
}

// NewsServer to create a new server base on worker size
// Please carefully to choose a suitable worker size, base on your application
// If the application needs lots of CPU resource, reduce the worker size
// If the application doesn't need lots of CPU resource and I/O resource instead, increase the worker size to increase the performance
func NewServer(workerSize int) *Server {
	return &Server{
		works:      make(chan *work.Work),
		End:        make(chan bool),
		Connectors: &connection.Connectors{},
		WorkerSize: workerSize,
	}
}

// Subscribe to register a new connection
func (server *Server) Subscribe() *connection.Connection {
	// Lock the Connectors
	server.Connectors.Lock()
	defer server.Connectors.Unlock()
	// Generate a new unique UUID for the connection
	u, _ := uuid.NewRandom()
	// Create a new connection
	con := &connection.Connection{
		ID:      u,
		Channel: make(chan *work.Work),
	}
	// Register the connection to server
	server.Connectors.Connections = append(server.Connectors.Connections, con)
	return con
}

// Unsubscribe to disconnect the connection
func (server *Server) Unsubscribe(connection *connection.Connection) {
	// Lock the Connectors
	server.Connectors.Lock()
	defer server.Connectors.Unlock()
	// Loop all the current connections and remove
	// The maximum time complexity is O(n)
	// And the logic doesn't keep the original order because of the order is doesn't matter
	for i, con := range server.Connectors.Connections {
		if con.ID == connection.ID {
			server.Connectors.Connections[i] = server.Connectors.Connections[len(server.Connectors.Connections)-1]
			server.Connectors.Connections[len(server.Connectors.Connections)-1] = nil
			server.Connectors.Connections = server.Connectors.Connections[:len(server.Connectors.Connections)-1]
			break
		}
	}
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
	for i < server.WorkerSize {
		i++
		// Generate the unique UUID
		u, _ := uuid.NewRandom()
		// Create a new Worker
		w := &worker.Worker{
			ID:            u,
			WorkerChannel: WorkerChannel,
			Channel:       make(chan *work.Work),
			End:           make(chan bool),
		}
		// Start the worker after creating
		w.Start(server.Connectors)
		// Register the worker to server
		server.workers = append(server.workers, w)
	}

	// Start handle the jobs
	go func() {
		for {
			select {
			case <-server.End:
				// If the server received the stop sign, loop ad stop all workers
				for _, w := range server.workers {
					w.Stop()
					return
				}
			case wo := <-server.works:
				// Watch to see if there is any available worker and let the worker do the job
				wkr := <-WorkerChannel
				wkr <- wo
			}
		}
	}()
}

// Status to return te current status of server
func (server *Server) Status() status.Status {
	return status.Status{
		Connections: len(server.Connectors.Connections),
		Workers:     server.WorkerSize,
		Packets:     len(server.works),
	}
}

// Stop stops the server by sending the stop sign to End channel
func (server *Server) Stop() {
	server.End <- true
}
