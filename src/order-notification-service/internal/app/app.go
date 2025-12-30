package app

import (
	"context"
	"log"
	"net/http"
	"order-notification-service/internal/infrastructure/msg_broker/proto/mappers"
	"time"

	"order-notification-service/internal/config"
	"order-notification-service/internal/infrastructure/msg_broker/kafka"
	"order-notification-service/internal/infrastructure/worker"
	"order-notification-service/internal/transport/websocket"
	"order-notification-service/internal/usecase/event_handling"
	"order-notification-service/internal/usecase/event_handling/order_status_changed"
	"order-notification-service/internal/usecase/notify_order_status_changed"

	"github.com/google/uuid"
)

type App struct {
	hub           *websocket.Hub
	handlerFn     http.HandlerFunc
	consumeWorker worker.BackgroundWorker
	conf          *config.Config
}

func Build(conf *config.Config) *App {
	kafkaConfig := kafka.NewKafkaConfig(&conf.Broker)
	ordEvtReader := kafka.NewKafkaReader(kafkaConfig, kafkaConfig.OrderNotificationEventTopic)

	stChMapper := proto_mappers.NewOrderStatusChangedPayloadMapper()
	cancelMapper := proto_mappers.NewOrderCancelledPayloadMapper()
	_ = cancelMapper

	hub := websocket.NewHub()
	hubJsonWrapper := websocket.NewNotifyHubJSONWrapper(hub)

	statusChangedNotifyUC := notify_order_status_changed.NewUsecase(hubJsonWrapper)

	statusChangedEventHandler := order_status_changed.NewEventHandler(
		statusChangedNotifyUC, stChMapper, kafkaConfig.OrderStatusChangedEventType)

	eventHands := []event_handling.EventHandler{
		statusChangedEventHandler,
	}

	msgProcessor := event_handling.NewEventMsgProcessor(eventHands)

	kafkaConsumerWorker := kafka.NewKafkaConsumerWorker(ordEvtReader, msgProcessor, conf.Broker.AllowedEventTypes)

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

	// to test
	// TODO: delete
	go func() {
		time.Sleep(30 * time.Second)

		testOrderID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")

		a.hub.NotifyOrder(testOrderID, []byte(`{"hello":"world"}`))
	}()

	http.HandleFunc("/ws", a.handlerFn)
	err := http.ListenAndServe(a.conf.Net.Addr, nil)
	if err != nil {
		log.Printf("err while serving http: %v", err)
	}
}
