package streamer

import (
	"github.com/chunchun02/streamer/internal/connection"
	"github.com/chunchun02/streamer/internal/status"
	"github.com/chunchun02/streamer/internal/work"
	"github.com/chunchun02/streamer/internal/worker"
	"github.com/google/uuid"
	"sync"
)

type Streamer interface {
	Subscribe() *connection.Connection
	Unsubscribe(connection *connection.Connection)
	In(packet interface{})
	Start()
	Status() status.Status
	Stop()
}

var WorkerChannel = make(chan chan *work.Work)

type Server struct {
	mu         sync.Mutex
	works      chan *work.Work
	workers    []*worker.Worker
	End        chan bool
	Connectors *connection.Connectors
	WorkerSize int
}

func NewServer(workerSize int) *Server {
	return &Server{
		works:      make(chan *work.Work),
		End:        make(chan bool),
		Connectors: &connection.Connectors{},
		WorkerSize: workerSize,
	}
}

func (server *Server) Subscribe() *connection.Connection {
	server.mu.Lock()
	defer server.mu.Unlock()
	u, _ := uuid.NewRandom()
	con := &connection.Connection{
		ID:      u,
		Channel: make(chan *work.Work),
	}
	server.Connectors.Connections = append(server.Connectors.Connections, con)
	return con
}

func (server *Server) Unsubscribe(connection *connection.Connection) {
	server.mu.Lock()
	defer server.mu.Unlock()
	for i, con := range server.Connectors.Connections {
		if con.ID == connection.ID {
			server.Connectors.Connections[i] = server.Connectors.Connections[len(server.Connectors.Connections)-1]
			server.Connectors.Connections[len(server.Connectors.Connections)-1] = nil
			server.Connectors.Connections = server.Connectors.Connections[:len(server.Connectors.Connections)-1]
			break
		}
	}
}

func (server *Server) In(packet interface{}) {
	server.works <- &work.Work{ID: 1, Data: packet}
}

func (server *Server) Start() {
	var i int
	for i < server.WorkerSize {
		i++
		u, _ := uuid.NewRandom()
		w := &worker.Worker{
			ID:            u,
			WorkerChannel: WorkerChannel,
			Channel:       make(chan *work.Work),
			End:           make(chan bool),
		}
		w.Start(server.Connectors)
		server.workers = append(server.workers, w)
	}

	go func() {
		for {
			select {
			case <-server.End:
				for _, w := range server.workers {
					w.Stop()
					return
				}
			case wo := <-server.works:
				wkr := <-WorkerChannel
				wkr <- wo
			}
		}
	}()
}

func (server *Server) Status() status.Status {
	return status.Status{
		Connections: len(server.Connectors.Connections),
		Workers:     server.WorkerSize,
		Packets:     len(server.works),
	}
}

func (server *Server) Stop() {
	server.End <- true
}
