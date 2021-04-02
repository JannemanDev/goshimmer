package client

import (
	"net"
	"sync"

	"github.com/iotaledger/goshimmer/packages/ledgerstate"
	"github.com/iotaledger/goshimmer/packages/txstream"
	"github.com/iotaledger/hive.go/events"
	"github.com/iotaledger/hive.go/logger"
)

// Client represents the client-side connection to a txstream server
type Client struct {
	clientID      string
	log           *logger.Logger
	chSend        chan txstream.Message
	chSubscribe   chan ledgerstate.Address
	chUnsubscribe chan ledgerstate.Address
	shutdown      chan bool
	Events        Events
	wgConnected   sync.WaitGroup
}

// Events contains all events emitted by the Client
type Events struct {
	// TransactionReceived is triggered when a transaction message is received from the server
	TransactionReceived *events.Event
	// InclusionStateReceived is triggered when an inclusion state message is received from the server
	InclusionStateReceived *events.Event
	// Connected is triggered when the client connects successfully to the server
	Connected *events.Event
}

// DialFunc is a function that performs the TCP connection to the server
type DialFunc func() (addr string, conn net.Conn, err error)

func handleTransactionReceived(handler interface{}, params ...interface{}) {
	handler.(func(*txstream.MsgTransaction))(params[0].(*txstream.MsgTransaction))
}

func handleInclusionStateReceived(handler interface{}, params ...interface{}) {
	handler.(func(*txstream.MsgTxInclusionState))(params[0].(*txstream.MsgTxInclusionState))
}

func handleConnected(handler interface{}, params ...interface{}) {
	handler.(func())()
}

// New creates a new client
func New(clientID string, log *logger.Logger, dial DialFunc) *Client {
	n := &Client{
		clientID:      clientID,
		log:           log,
		chSend:        make(chan txstream.Message),
		chSubscribe:   make(chan ledgerstate.Address),
		chUnsubscribe: make(chan ledgerstate.Address),
		shutdown:      make(chan bool),
		Events: Events{
			TransactionReceived:    events.NewEvent(handleTransactionReceived),
			InclusionStateReceived: events.NewEvent(handleInclusionStateReceived),
			Connected:              events.NewEvent(handleConnected),
		},
	}

	go n.subscriptionsLoop()
	go n.connectLoop(dial)

	return n
}

// Close shuts down the client
func (n *Client) Close() {
	close(n.shutdown)
	n.Events.TransactionReceived.DetachAll()
	n.Events.InclusionStateReceived.DetachAll()
	n.Events.Connected.DetachAll()
}