package signaling

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
	"log/slog"
	"net/http"
	"sync"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type peer struct {
	Conn     *websocket.Conn
	peerConn *webrtc.PeerConnection
}

func newPeer(conn *websocket.Conn) *peer {
	return &peer{Conn: conn}
}

type Handler struct {
	peers map[uuid.UUID]*peer
	lock  sync.Mutex
	log   *slog.Logger
}

func NewHandler(log *slog.Logger) *Handler {
	return &Handler{
		peers: make(map[uuid.UUID]*peer),
		log:   log,
	}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.log.Error(fmt.Errorf("failed to upgrade connection: %w", err).Error())
		return
	}
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			h.log.Error(fmt.Errorf("failed to close connection: %w", err).Error())
		}
	}(conn)

	id := uuid.New()
	h.log.Info(fmt.Sprintf("new peer connected: %s", id))

	h.lock.Lock()
	h.peers[id] = newPeer(conn)
	h.lock.Unlock()

	return
}
