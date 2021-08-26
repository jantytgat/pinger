package pinger

type Host struct {
	Address         string
	PacketsSent     int
	PacketsReceived int
	Available       bool
	PacketLoss      float64
	Comment         string
}
