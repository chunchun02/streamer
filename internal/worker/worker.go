package worker

import (
	"github.com/chunchun02/streamer/pkg/connection"
	"github.com/chunchun02/streamer/pkg/work"
	"github.com/google/uuid"
)

// Workers repents for a do-er, will do the job
type Worker struct {
	id            uuid.UUID            // unique UUID for each worker
	workerChannel chan chan *work.Work // Channel to communicate with server
	channel       chan *work.Work      // Work channel of a job need to be done
	end           chan bool            // End channel
}

func NewWorker(workerChannel chan chan *work.Work) *Worker {
	// Generate the unique UUID
	u, _ := uuid.NewRandom()
	// Create a new Worker
	return &Worker{
		id:            u,
		workerChannel: workerChannel,
		channel:       make(chan *work.Work),
		end:           make(chan bool),
	}
}

// Start starts worker to do the job. Loop the connection list and push the data to connection work channel
func (w *Worker) Start(connector *connection.Connectors) {
	go func() {
		for {
			// Communicate with server
			w.workerChannel <- w.channel
			// Watch to see if there is any work need to be done
			select {
			case wo := <-w.channel:
				// Loop all the current connections and do the job for user connection
				for _, con := range connector.Connections {
					// Multi-threading, so in the case connection is unsubscribed, have to ignore it
					if con != nil {
						con.Channel <- wo
					}
				}
			case <-w.end:
				// The worker received the end sign
				return
			}
		}
	}()
}

// Stop stops the worker by send the done flag to End channel
func (w *Worker) Stop() {
	w.end <- true
}
