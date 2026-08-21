package observability

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type Metrics struct {
	httpRequests       atomic.Uint64
	reviewsCreated     atomic.Uint64
	websocketActive    atomic.Int64
	websocketBroadcast atomic.Uint64
}

func (metrics *Metrics) RecordHTTPRequest()        { metrics.httpRequests.Add(1) }
func (metrics *Metrics) RecordReviewCreated()      { metrics.reviewsCreated.Add(1) }
func (metrics *Metrics) RecordWebSocketConnected() { metrics.websocketActive.Add(1) }
func (metrics *Metrics) RecordWebSocketClosed()    { metrics.websocketActive.Add(-1) }
func (metrics *Metrics) RecordWebSocketBroadcast() { metrics.websocketBroadcast.Add(1) }
func (metrics *Metrics) ActiveWebSockets() int64   { return metrics.websocketActive.Load() }

func (metrics *Metrics) ServeHTTP(writer http.ResponseWriter, _ *http.Request) {
	writer.Header().Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")
	_, _ = fmt.Fprintf(writer, "# TYPE viger_http_requests_total counter\nviger_http_requests_total %d\n", metrics.httpRequests.Load())
	_, _ = fmt.Fprintf(writer, "# TYPE viger_reviews_created_total counter\nviger_reviews_created_total %d\n", metrics.reviewsCreated.Load())
	_, _ = fmt.Fprintf(writer, "# TYPE viger_websocket_connections gauge\nviger_websocket_connections %d\n", metrics.websocketActive.Load())
	_, _ = fmt.Fprintf(writer, "# TYPE viger_websocket_broadcasts_total counter\nviger_websocket_broadcasts_total %d\n", metrics.websocketBroadcast.Load())
}
