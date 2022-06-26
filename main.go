package main

import (
	"context"
	"fmt"

	codecs "aggregator/codecs"
	models "aggregator/models"
	txnCollector "aggregator/txnCollector"
	windowBuilder "aggregator/windowBuilder"

	"github.com/lovoo/goka"
	"go.uber.org/zap"
)

var (
	brokers = []string{"localhost:19092"}
	logger, _ = zap.NewProduction()
	aggTopic = new(codecs.ArrayCodec)
	srcStream goka.Stream = "btc"
)

func verifyTopic(tmgr goka.TopicManager, topic string) {
	err := tmgr.EnsureStreamExists(topic, 10)
	if err != nil {
		logger.Fatal(fmt.Sprintf("Error creating topic: %s", topic), zap.String("Error", err.Error()))
	}
}

func main() {
	ctx := context.Background()

	tmgr, err := goka.NewTopicManager(brokers, goka.DefaultConfig(), goka.NewTopicManagerConfig())
	if err != nil {
		logger.Fatal("error creating topic manager", zap.String("Error", err.Error()))
	}

	verifyTopic(tmgr, "btc")
	srcTopic := &models.Topic{Stream: srcStream, Codec: new(codecs.TxnCodec)}

	// RUN TXN collector
	btcCollector := txnCollector.TxnCollector{Brokers: brokers, Topic: srcTopic}
	go btcCollector.RunBTCCollector(ctx)


	verifyTopic(tmgr, "events")

	done := make(chan bool)

	wb := &windowBuilder.WindowBuilder{Logger: logger, SourceTopic: srcTopic, AggTopic: aggTopic, Done: done}
	
	go wb.Run(ctx, brokers)

	runView()

	<-done
}
