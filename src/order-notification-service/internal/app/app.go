package app

import (
	"context"
	"log"
	"net/http"

	"order-notification-service/internal/config"
	"order-notification-service/internal/infrastructure/messaging/kafka"
	"order-notification-service/internal/infrastructure/worker"
	protostatuschanged "order-notification-service/internal/transport/proto/order_status_changed"
	"order-notification-service/internal/transport/websocket"
	"order-notification-service/internal/usecase/event_handling"
	"order-notification-service/internal/usecase/event_handling/order_status_changed"
	"order-notification-service/internal/usecase/notify_order_status_changed"
)

type App struct {
	hub           *websocket.Hub
	handlerFn     http.HandlerFunc
	consumeWorker worker.BackgroundWorker
	conf          *config.Config
}

func Build(conf *config.Config) *App {
	kafkaConfig := kafka.NewKafkaConfig(&conf.Broker)

	ordNotifyReader := kafka.NewKafkaReader(kafkaConfig, kafkaConfig.OrderNotificationEventTopic)

	stChMapper := protostatuschanged.NewPayloadMapper()

	hub := websocket.NewHub()

	statusChangedNotifyUC := notify_order_status_changed.NewUsecase(hub)

	statusChangedEventHandler := order_status_changed.NewEventHandler(
		statusChangedNotifyUC, stChMapper, kafkaConfig.OrderStatusChangedEventType)

	eventHands := []event_handling.EventHandler{
		statusChangedEventHandler,
	}

	msgProcessor := event_handling.NewEventMsgProcessor(eventHands)

	kafkaConsumerWorker := kafka.NewKafkaConsumerWorker(ordNotifyReader, msgProcessor)

	return &App{
		conf:          conf,
		hub:           hub,
		handlerFn:     websocket.NewHandler(hub),
		consumeWorker: kafkaConsumerWorker,
	}
}

func (a *App) Run() {
	go a.hub.Run()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// launching broker consume worker
	go func(worker worker.BackgroundWorker) {
		err := worker.Run(ctx)
		if err != nil {
			log.Printf("consume worker stopped with error: %v", err)
			cancel()
		}
	}(a.consumeWorker)

	http.HandleFunc("/ws", websocket.NewHandler(a.hub))
	err := http.ListenAndServe(a.conf.Net.Addr, nil)
	if err != nil {
		log.Printf("err while serving http: %v", err)
	}
}
