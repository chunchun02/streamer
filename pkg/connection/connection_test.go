package connection

import (
	"log"
	"testing"
)

func TestRegisterConnection_OnlyOneConnection(t *testing.T) {
	connectors := &Connectors{}
	connectors.RegisterConnection()
	if len(connectors.Connections) != 1 {
		log.Fatal("Can not register a connection")
	}
}

func TestRegisterConnection_TwoConnection(t *testing.T) {
	connectors := &Connectors{}
	connectors.RegisterConnection()
	connectors.RegisterConnection()
	if len(connectors.Connections) != 2 {
		log.Fatal("Can not register a connection")
	}
}

func TestUnregisterConnection(t *testing.T) {
	connectors := &Connectors{}
	con1 := connectors.RegisterConnection()
	con2 := connectors.RegisterConnection()
	if len(connectors.Connections) != 2 {
		log.Fatal("Can not register a connection")
	}

	connectors.UnregisterConnection(con1)
	if len(connectors.Connections) != 1 {
		log.Fatal("Can not unregister a connection")
	}

	connectors.UnregisterConnection(con2)
	if len(connectors.Connections) != 0 {
		log.Fatal("Can not unregister a connection")
	}
}
