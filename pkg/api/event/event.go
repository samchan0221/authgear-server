package event

import "fmt"

type Type string

type Payload interface {
	UserID() string
	IsAdminAPI() bool
	FillContext(ctx *Context)
}

type BlockingPayload interface {
	Payload
	BlockingEventType() Type
}

type NonBlockingPayload interface {
	Payload
	NonBlockingEventType() Type
	ForWebHook() bool
	ForAudit() bool
}

type Event struct {
	ID            string  `json:"id"`
	Seq           int64   `json:"seq"`
	Type          Type    `json:"type"`
	Payload       Payload `json:"payload"`
	Context       Context `json:"context"`
	IsNonBlocking bool    `json:"-"`
}

func newEvent(seq int64, payload Payload, context Context) *Event {
	e := &Event{
		Payload: payload,
		Context: context,
	}
	e.Seq = seq
	e.ID = fmt.Sprintf("%016x", seq)
	return e
}

func NewBlockingEvent(seqNo int64, payload BlockingPayload, context Context) *Event {
	event := newEvent(seqNo, payload, context)
	event.Type = payload.BlockingEventType()
	return event
}

func NewNonBlockingEvent(seqNo int64, payload NonBlockingPayload, context Context) *Event {
	event := newEvent(seqNo, payload, context)
	event.Type = payload.NonBlockingEventType()
	event.IsNonBlocking = true
	return event
}
