package stats

import (
	"github.com/EthereumHD/Scan/src/statistics/model"
	"errors"
	"log"
	"sync"
	"time"
)

// ErrNodeNotAuthorized is returned when a report is received for a node that has
// not been authorized yet.
var ErrNodeNotAuthorized = errors.New("node has not been authorized")

// ErrAuthFailed is returned when a node fails to authorize.
var ErrAuthFailed = errors.New("authorization failed")

// ErrInvalidReport is returned when the collector receives an invalid type.
var ErrInvalidReport = errors.New("invalid report")

// Node contains all the stats metadata about an Ethereum node.
type Node struct {
	ID       model.ID      `json:"id"`
	Info     model.Info    `json:"info"`
	Latency  model.Latency `json:"latency"`
	Block    model.Block   `json:"block"`
	Pending  model.Pending `json:"pending"`
	Status   model.Status  `json:"status"`
	LastSeen time.Time     `json:"last_seen"`
}

type collector struct {
	mu    sync.RWMutex
	nodes map[model.ID]Node
	auth  func(model.AuthReport) bool
}

func (col *collector) Collect(report model.Report) error {
	col.mu.Lock()
	defer col.mu.Unlock()

	// TODO: Uncollect on disconnect? Or sweep based on last seen?
	if col.nodes == nil {
		(*col).nodes = map[model.ID]Node{}
	}

	if authReport, ok := report.(model.AuthReport); ok {
		if col.auth != nil && !col.auth(authReport) {
			return ErrAuthFailed
		}
		col.nodes[authReport.ID] = Node{
			ID:       authReport.ID,
			Info:     authReport.Info,
			LastSeen: time.Now(),
		}
		log.Printf("collected node: %s", authReport.ID)
		return nil
	}

	node, ok := col.nodes[report.NodeID()]
	if !ok {
		return ErrNodeNotAuthorized
	}
	node.LastSeen = time.Now()

	switch report := report.(type) {
	case model.LatencyReport:
		node.Latency = report.Latency
	case model.BlockReport:
		node.Block = report.Block
	case model.PendingReport:
		node.Pending = report.Pending
	case model.StatusReport:
		node.Status = report.Status
	case model.DisconnectReport:
		delete(col.nodes, report.NodeID())
	}

	return nil
}

// Get returns a Node with the given ID, if it has been collected.
func (col *collector) Get(ID model.ID) (Node, bool) {
	col.mu.RLock()
	defer col.mu.RUnlock()

	node, ok := col.nodes[ID]
	return node, ok
}

// List returns a slice of IDs that are being collected.
func (col *collector) List() []model.ID {
	col.mu.RLock()
	defer col.mu.RUnlock()

	ids := make([]model.ID, 0, len(col.nodes))
	for id := range col.nodes {
		ids = append(ids, id)
	}
	return ids
}
