package status

// Status presents for the status of server
type Status struct {
	Connections int // Number of current connections
	Workers     int // Number of current workers
	Packets     int // Number of current data packets
}
