package sergo

// Room :
type Room struct {
	sockets map[*Socket]bool
	joinCh  chan *Socket
	leaveCh chan *Socket
}

// NewRoom :
func NewRoom() *Room {
	r := &Room{
		sockets: make(map[*Socket]bool),
	}
	go r.run()
	return r
}

func (r *Room) run() {
	select {
	case sckt := <-r.joinCh:
		r.sockets[sckt] = true
	case sckt := <-r.leaveCh:
		delete(r.sockets, sckt)

	default:

	}
}

//Add :
func (r *Room) Add(s *Socket) {
	r.joinCh <- s
}

// Remove :
func (r *Room) Remove(s *Socket) {
	r.leaveCh <- s
}

//Emit :
func (r Room) Emit(name string, data string) {
	for sck := range r.sockets {
		sck.Emit(name, data)
	}
}
