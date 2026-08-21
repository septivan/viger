package realtime

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/coder/websocket"
	review "github.com/septivan/viger/backend/internal/core/review/domain"
	"github.com/septivan/viger/backend/internal/platform/observability"
)

var ErrConnectionLimit = errors.New("WebSocket connection limit reached")

// Hub broadcasts review notifications; REST remains the authoritative data source.
type Hub struct {
	mutex          sync.RWMutex
	clients        map[*client]struct{}
	allowedOrigins map[string]struct{}
	maximumClients int
	metrics        *observability.Metrics
}

type client struct {
	connection *websocket.Conn
	writes     sync.Mutex
}

type ReviewCreatedEvent struct {
	Type   string        `json:"type"`
	GameID string        `json:"gameId"`
	Review ReviewPayload `json:"review"`
}

type ReviewPayload struct {
	ID           string    `json:"id"`
	GameID       string    `json:"gameId"`
	ReviewerName string    `json:"reviewerName"`
	Rating       int       `json:"rating"`
	Text         string    `json:"text"`
	CreatedAt    time.Time `json:"createdAt"`
}

// NewHub creates an origin-restricted, connection-bounded WebSocket hub.
func NewHub(allowedOrigins []string, maximumClients int, metrics *observability.Metrics) *Hub {
	origins := make(map[string]struct{}, len(allowedOrigins))
	for _, origin := range allowedOrigins {
		origins[strings.TrimSuffix(origin, "/")] = struct{}{}
	}
	return &Hub{clients: make(map[*client]struct{}), allowedOrigins: origins, maximumClients: maximumClients, metrics: metrics}
}

func (hub *Hub) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if !hub.originAllowed(request.Header.Get("Origin")) {
		http.Error(writer, "WebSocket origin is not allowed", http.StatusForbidden)
		return
	}
	if hub.connectionCount() >= hub.maximumClients {
		http.Error(writer, ErrConnectionLimit.Error(), http.StatusServiceUnavailable)
		return
	}
	connection, err := websocket.Accept(writer, request, &websocket.AcceptOptions{InsecureSkipVerify: true})
	if err != nil {
		return
	}
	connection.SetReadLimit(1024)
	connected := &client{connection: connection}
	hub.add(connected)
	defer func() {
		hub.remove(connected)
		_ = connection.Close(websocket.StatusNormalClosure, "connection closed")
	}()

	for {
		if _, _, err = connection.Read(request.Context()); err != nil {
			return
		}
	}
}

func (hub *Hub) PublishReviewCreated(_ context.Context, review review.Review) error {
	event := ReviewCreatedEvent{
		Type: "review.created", GameID: review.GameID,
		Review: ReviewPayload{ID: review.ID, GameID: review.GameID, ReviewerName: review.ReviewerName, Rating: review.Rating, Text: review.Text, CreatedAt: review.CreatedAt},
	}
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}
	hub.metrics.RecordWebSocketBroadcast()
	for _, connected := range hub.snapshot() {
		writeContext, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		connected.writes.Lock()
		err = connected.connection.Write(writeContext, websocket.MessageText, payload)
		connected.writes.Unlock()
		cancel()
		if err != nil {
			hub.remove(connected)
			_ = connected.connection.CloseNow()
		}
	}
	return nil
}

func (hub *Hub) originAllowed(origin string) bool {
	_, allowed := hub.allowedOrigins[strings.TrimSuffix(origin, "/")]
	return allowed
}

func (hub *Hub) connectionCount() int {
	hub.mutex.RLock()
	defer hub.mutex.RUnlock()
	return len(hub.clients)
}

func (hub *Hub) add(connected *client) {
	hub.mutex.Lock()
	hub.clients[connected] = struct{}{}
	hub.mutex.Unlock()
	hub.metrics.RecordWebSocketConnected()
}

func (hub *Hub) remove(connected *client) {
	hub.mutex.Lock()
	if _, found := hub.clients[connected]; found {
		delete(hub.clients, connected)
		hub.metrics.RecordWebSocketClosed()
	}
	hub.mutex.Unlock()
}

func (hub *Hub) snapshot() []*client {
	hub.mutex.RLock()
	defer hub.mutex.RUnlock()
	result := make([]*client, 0, len(hub.clients))
	for connected := range hub.clients {
		result = append(result, connected)
	}
	return result
}
