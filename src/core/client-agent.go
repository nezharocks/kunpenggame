package core

// ClientAgent is
type ClientAgent struct {
	CenterID   int
	CenterName string
	ServerIP   string
	ServerPort int
	StopCh     chan struct{}
}

// NewClientAgent is
func NewClientAgent(id int, name string, serverIP string, serverPort int) *ClientAgent {
	return &ClientAgent{
		CenterID:   id,
		CenterName: name,
		ServerIP:   serverIP,
		ServerPort: serverPort,
		StopCh:     make(chan struct{}, 1),
	}
}
