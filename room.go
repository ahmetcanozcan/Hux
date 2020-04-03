package sergo

// Room :
type Room struct {
	sockets map[*Socket]bool
}

// NewRoom :
func NewRoom() *Room {
	return &Room{
		sockets: make(map[*Socket]bool),
	}
}

//Emit :
func (r Room) Emit(name string, data string) {
	for sck := range r.sockets {
		sck.Emit(name, data)
	}
}
