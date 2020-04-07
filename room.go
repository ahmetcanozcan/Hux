package hux

// Room : a Room is  a socket group
type Room struct {
	sockets map[*Socket]bool
	joinCh  chan *Socket
	leaveCh chan *Socket
}

// NewRoom : Instantiate a new room
func NewRoom() *Room {
	r := &Room{
		sockets: make(map[*Socket]bool),
		joinCh:  make(chan *Socket),
		leaveCh: make(chan *Socket),
	}
	go r.run()
	return r
}

func (r *Room) run() {
	for {
		select {
		case sckt := <-r.joinCh:
			r.sockets[sckt] = true
		case sckt := <-r.leaveCh:
			delete(r.sockets, sckt)
		default:

		}
	}
}

//Add : add a socket to the room
func (r *Room) Add(s *Socket) {
	r.joinCh <- s
}

// Remove : remove a socket from the room
func (r *Room) Remove(s *Socket) {
	r.leaveCh <- s
}

// Emit : emit message to all sockets in the room
func (r *Room) Emit(name string, data interface{}) {
	for sck := range r.sockets {
		sck.Emit(name, data)
	}
}
