package worker

import (
	"github.com/chunchun02/streamer/internal/connection"
	"github.com/chunchun02/streamer/internal/work"
	"github.com/google/uuid"
)

// Workers repents for a do-er, will do the job
type Worker struct {
	ID            uuid.UUID            // unique UUID for each worker
	WorkerChannel chan chan *work.Work // Channel to communicate with server
	Channel       chan *work.Work      // Work channel of a job need to be done
	End           chan bool            // End channel
}

// Start starts worker to do the job. Loop the connection list and push the data to connection work channel
func (w *Worker) Start(connector *connection.Connectors) {
	go func() {
		for {
			// Communicate with server
			w.WorkerChannel <- w.Channel
			// Watch to see if there is ay work need to be done
			select {
			case wo := <-w.Channel:
				// Loop all the current connections and do the job for user connection
				for _, con := range connector.Connections {
					// Multi-threading, so in the case connection is unsubscribed, have to ignore
					if con != nil {
						con.Channel <- wo
					}
				}
			case <-w.End:
				// The worker receive the end sign
				return
			}
		}
	}()
}

// Stop stops the worker by send the done flag to End channel
func (w *Worker) Stop() {
	w.End <- true
}
