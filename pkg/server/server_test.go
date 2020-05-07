package server

import (
	"log"
	"testing"
)

func TestServer_In(t *testing.T) {
	sv := NewServer(10, 10)
	sv.In(1)
	if len(sv.works) != 1 {
		log.Fatal("Can not push data into the server")
	}

	sv.In(2)
	if len(sv.works) != 2 {
		log.Fatal("Can not push data into the server")
	}
}

func TestServer_Subscribe(t *testing.T) {
	sv := NewServer(10, 10)
	sv.Subscribe()

	if len(sv.connectors.Connections) != 1 {
		log.Fatal("Can not subscribe a new connection")
	}
}

func TestServer_Unsubscribe(t *testing.T) {
	sv := NewServer(10, 10)
	con := sv.Subscribe()
	if len(sv.connectors.Connections) != 1 {
		log.Fatal("Can not subscribe a new connection")
	}

	sv.Unsubscribe(con)
	if len(sv.connectors.Connections) != 0 {
		log.Fatal("Can not unsubscribe a new connection")
	}
}

func TestServer_Status(t *testing.T) {
	sv := NewServer(10, 10)
	stt := sv.Status()

	if stt.WorkerSize != 10 {
		log.Fatal("Worker size is not correct")
	}

	if stt.Connections != 0 {
		log.Fatal("The number of connection is not correct")
	}

	if stt.WorkQueueSize != 10 {
		log.Fatal("Work queue size is not correct")
	}

	if stt.Works != 0 {
		log.Fatal("Current queuing work is not correct")
	}
}

func TestServer_Stop(t *testing.T) {
	sv := NewServer(10, 10)
	sv.Start()
	sv.Stop()
	if len(sv.end) != 0 {
		log.Fatal("Can not stop server")
	}
}

func TestServer_Start(t *testing.T) {
	sv := NewServer(10, 10)
	sv.Start()

	if len(sv.workers) != 10 {
		log.Fatal("Can not create workers")
	}
}
