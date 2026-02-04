package ws

type Hub struct {
	// roomID -> clients
	rooms map[uint]map[*Client]bool

	register   chan *Client
	unregister chan *Client
	broadcast  chan RoomMessage
}

type RoomMessage struct {
	RoomID uint
	Data   []byte
}

func NewHub() *Hub {
	return &Hub{
		rooms:      make(map[uint]map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan RoomMessage, 256),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case c := <-h.register:
			if h.rooms[c.roomID] == nil {
				h.rooms[c.roomID] = make(map[*Client]bool)
			}
			h.rooms[c.roomID][c] = true

		case c := <-h.unregister:
			if h.rooms[c.roomID] != nil {
				if _, ok := h.rooms[c.roomID][c]; ok {
					delete(h.rooms[c.roomID], c)
					close(c.send)
					if len(h.rooms[c.roomID]) == 0 {
						delete(h.rooms, c.roomID)
					}
				}
			}

		case msg := <-h.broadcast:
			if h.rooms[msg.RoomID] == nil {
				continue
			}
			for c := range h.rooms[msg.RoomID] {
				select {
				case c.send <- msg.Data:
				default:
					// client trop lent -> drop
					delete(h.rooms[msg.RoomID], c)
					close(c.send)
				}
			}
		}
	}
}
