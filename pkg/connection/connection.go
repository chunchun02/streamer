package connection

import (
	"github.com/chunchun02/streamer/pkg/work"
	"github.com/google/uuid"
)

// Connection presents for a connection between user and server
type Connection struct {
	ID      uuid.UUID       // Unique UUID for each connection
	Channel chan *work.Work // Channel to send data to user
}

// Connectors present for a list of all current connecting users
type Connectors struct {
	Connections []*Connection // List of all current connections
}

// Register a new connection to Connectors
func (connectors *Connectors) RegisterConnection() *Connection {
	// Generate a new unique UUID for the connection
	u, _ := uuid.NewRandom()
	// Create a new connection
	con := &Connection{
		ID:      u,
		Channel: make(chan *work.Work),
	}
	// Register the connection to server
	connectors.Connections = append(connectors.Connections, con)
	return con
}

// Unregister a connection from Connectors
func (connectors *Connectors) UnregisterConnection(connection *Connection) {
	// Loop all the current connections and remove
	// The maximum time complexity is O(n)
	// And the logic doesn't keep the original order because of the order is doesn't matter
	for i, con := range connectors.Connections {
		if con.ID == connection.ID {
			connectors.Connections[i] = connectors.Connections[len(connectors.Connections)-1]
			connectors.Connections[len(connectors.Connections)-1] = nil
			connectors.Connections = connectors.Connections[:len(connectors.Connections)-1]
			break
		}
	}
}
