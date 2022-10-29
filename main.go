package main

import (
	"context"
	"fmt"

	codecs "aggregator/service/codecs"
	models "aggregator/service/models"
	views "aggregator/service/views"
	txnCollector "aggregator/service/txnCollector"
	windowBuilder "aggregator/service/windowBuilder"

	"github.com/lovoo/goka"
	"go.uber.org/zap"
)

var (
	brokers = []string{"localhost:9092"}
	logger, _ = zap.NewProduction()
	aggTopic = new(codecs.ArrayCodec)
	srcStream goka.Stream = "btc"
)

func verifyTopic(topic string) {
	
	tmgr, err := goka.NewTopicManager(brokers, goka.DefaultConfig(), goka.NewTopicManagerConfig())
	if err != nil {
		logger.Fatal("error creating topic manager", zap.String("Error", err.Error()))
	}

	err = tmgr.EnsureStreamExists(topic, 10)
	if err != nil {
		logger.Fatal(fmt.Sprintf("Error creating topic: %s", topic), zap.String("Error", err.Error()))
	}
}

func main() {
	ctx := context.Background()

	verifyTopic("btc")
	srcTopic := &models.Topic{Stream: srcStream, Codec: new(codecs.TxnCodec)}

	// RUN TXN collector
	btcCollector := txnCollector.TxnCollector{Brokers: brokers, Topic: srcTopic}
	go btcCollector.RunBTCCollector(ctx)

	done := make(chan bool)

	wb := &windowBuilder.WindowBuilder{Logger: logger, SourceTopic: srcTopic, AggTopic: aggTopic, Done: done}
	err := wb.Init(brokers)

	go wb.Run(ctx, brokers)
	if err != nil {
		logger.Fatal("Error Initializing Window Builder", zap.String("Error", err.Error()))
	}
	views.RunViews(brokers)

	<-done
}
