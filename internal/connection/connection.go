package connection

import (
	"github.com/chunchun02/streamer/internal/work"
	"github.com/google/uuid"
)

type Connection struct {
	ID      uuid.UUID
	Work    *work.Work
	Channel chan *work.Work
}

type Connectors struct {
	Connections []*Connection
}

func (con *Connection) Stream() {
	con.Channel <- con.Work
}
