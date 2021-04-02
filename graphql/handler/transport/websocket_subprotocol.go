package transport

import (
	"encoding/json"
	"errors"
)

const (
	initMessageType messageType = iota
	connectionAckMessageType
	keepAliveMessageType
	connectionErrorMessageType
	connectionCloseMessageType
	startMessageType
	stopMessageType
	dataMessageType
	completeMessageType
	errorMessageType
)

var (
	supportedSubprotocols = []string{
		graphqlwsSubprotocol,
	}

	errWsConnClosed = errors.New("websocket connection closed")
	errInvalidMsg   = errors.New("invalid message received")
)

type (
	messageType int
	message     struct {
		payload json.RawMessage
		id      string
		t       messageType
	}
	messageExchanger interface {
		NextMessage() (message, error)
		Send(m *message) error
	}
)

func (t messageType) String() string {
	var text string
	switch t {
	default:
		text = "unknown"
	case initMessageType:
		text = "init"
	case connectionAckMessageType:
		text = "connection ack"
	case keepAliveMessageType:
		text = "keep alive"
	case connectionErrorMessageType:
		text = "connection error"
	case connectionCloseMessageType:
		text = "connection close"
	case startMessageType:
		text = "start"
	case stopMessageType:
		text = "stop subscription"
	case dataMessageType:
		text = "data"
	case completeMessageType:
		text = "complete"
	case errorMessageType:
		text = "error"
	}
	return text
}

func contains(list []string, elem string) bool {
	for _, e := range list {
		if e == elem {
			return true
		}
	}

	return false
}

func (t *Websocket) injectGraphQLWSSubprotocols() {
	// the list of subprotocols is specified by the consumer of the Websocket struct,
	// in order to preserve backward compatibility, we inject the graphql specific subprotocols
	// at runtime
	if !t.didInjectSubprotocols {
		defer func() {
			t.didInjectSubprotocols = true
		}()

		for _, subprotocol := range supportedSubprotocols {
			if !contains(t.Upgrader.Subprotocols, subprotocol) {
				t.Upgrader.Subprotocols = append(t.Upgrader.Subprotocols, subprotocol)
			}
		}
	}
}
