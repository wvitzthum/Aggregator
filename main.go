package main

import (
	"context"
	"fmt"
	"sync"

	codecs "aggregator/service/codecs"
	featureCalculator "aggregator/service/featureCalculator"
	models "aggregator/service/models"
	txnCollector "aggregator/service/txnCollector"
	views "aggregator/service/views"
	windowBuilder "aggregator/service/windowBuilder"

	"github.com/lovoo/goka"
	"go.uber.org/zap"
)

var (
	brokers = []string{"localhost:9092"}
	logger, _ = zap.NewProduction()
	txnTopic goka.Stream = "btc"
	windowTopic goka.Stream = "window-table"
	featureTopic goka.Stream = "features"
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

	verifyTopic(string(txnTopic))
	verifyTopic(string(featureTopic))
	txnStream := &models.Topic{Stream: &txnTopic, Codec: new(codecs.TxnCodec)}
	windowStream := &models.Topic{Stream: &windowTopic, Codec: new(codecs.ArrayCodec)}
	featureStream := &models.Topic{Stream: &featureTopic, Codec: new(codecs.FeaturesCodec)}

	

	// RUN TXN collector
	btcCollector := txnCollector.TxnCollector{Brokers: brokers, Topic: txnStream}
	go btcCollector.RunBTCCollector(ctx)

	wg := sync.WaitGroup{}

	wb := &windowBuilder.WindowBuilder{Logger: logger, SourceTopic: txnStream, OutTopic: windowStream, Brokers: brokers}
	err := wb.Init()
	if err != nil {
		logger.Fatal("Error Initializing Window Builder", zap.String("Error", err.Error()))
	}
	
	fc := &featureCalculator.FeatureCalculator{Logger: logger, SourceTopic: windowStream, OutTopic: featureStream, Brokers: brokers}
	err = fc.Init()
	if err != nil {
		logger.Fatal("Error Feature Calc", zap.String("Error", err.Error()))
	}
	wg.Add(1)
	go wb.Run(ctx, wg)
	wg.Add(1)
	go fc.Run(ctx, wg)
	views.RunViews(brokers)

	wg.Wait()
}
