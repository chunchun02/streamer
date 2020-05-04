package worker

import (
	"github.com/chunchun02/streamer/internal/connection"
	"github.com/chunchun02/streamer/internal/work"
	"github.com/google/uuid"
)

type Worker struct {
	ID            uuid.UUID
	WorkerChannel chan chan *work.Work
	Channel       chan *work.Work
	End           chan bool
}

func (w *Worker) Start(connector *connection.Connectors) {
	go func() {
		for {
			w.WorkerChannel <- w.Channel
			select {
			case wo := <-w.Channel:
				for _, con := range connector.Connections {
					if con != nil {
						con.Work = wo
						con.Stream()
					}
				}
			case <-w.End:
				return
			}
		}
	}()
}

func (w *Worker) Stop() {
	w.End <- true
}
