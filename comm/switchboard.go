package comm

// SB is the global communication switch board
var SB SwitchBoard

// SwitchBoard is the communication struct for plugins
type SwitchBoard struct {
	ChanResponse  chan *Response
	ChanPersist   chan Persist
	ChanHeartbeat chan string
}
