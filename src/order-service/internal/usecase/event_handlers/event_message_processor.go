package event_handlers

import (
	"context"
	"fmt"
	"log"
)

type EventMsgProcessor struct {
	handlers map[string]EventHandler
}

func NewEventMsgProcessor(hands []EventHandler) *EventMsgProcessor {
	var handlers = make(map[string]EventHandler)
	for _, h := range hands {
		handlers[h.EventType()] = h
	}

	return &EventMsgProcessor{
		handlers: handlers,
	}
}

func (p *EventMsgProcessor) Process(ctx context.Context, payload []byte, eventType string) error {
	handler := p.handlers[eventType]
	if handler == nil {
		return fmt.Errorf("event handler for event %s not found", eventType)
	}

	err := handler.Handle(ctx, payload)
	if err != nil {
		log.Printf("Error handling event %s: %v", eventType, err)
	}

	return nil
}
