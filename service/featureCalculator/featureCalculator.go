package featureCalculator

import (
	"aggregator/service/models"
	"context"
	"fmt"
	"sync"
	"time"

	codecs "aggregator/service/codecs"

	"github.com/go-redis/redis/v9"
	"github.com/lovoo/goka"
	"go.uber.org/zap"
)

type FeatureCalculator struct {
	Logger      *zap.Logger
	RedisClient *redis.Client
	processor   *goka.Processor
	Brokers     []string
	SourceTopic *models.Topic
	OutTopic    *models.Topic
}

func (f *FeatureCalculator) outboundStatsProcessor(gctx goka.Context, msg interface{}) {
	window, ok := msg.([]string)
	if !ok {
		gctx.Fail(fmt.Errorf("couldn't convert value to Window"))
	}
	txns := []*models.Txn{}
	if len(window) == 0 {
		return
	}
	for _, v := range window {
		// TODO: refactor this into its own function
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		val, err := f.RedisClient.Get(ctx, v).Result()
		if err == redis.Nil || err != nil {
			f.Logger.Error("key does not exist", zap.String("key", v))
			continue
		}
		txn, err := codecs.TxnDecodeGOB([]byte(val))
		if err != nil {
			f.Logger.Error("Unable to decode", zap.String("key", v))
			continue
		}
		txns = append(txns, txn)
	}

	features := CalcFeatures(txns)

	// emit new statistics without changing the key
	gctx.Emit(*f.OutTopic.Stream, gctx.Key(), features)
}

func (f *FeatureCalculator) Init() error {
	var err error

	f.processor, err = goka.NewProcessor(f.Brokers, goka.DefineGroup("outboundStats",
		goka.Input(*f.SourceTopic.Stream, f.SourceTopic.Codec, f.outboundStatsProcessor),
		goka.Output(*f.OutTopic.Stream, f.OutTopic.Codec)),
	)
	if err != nil {
		f.Logger.Error(err.Error())
	}
	return err
}

func (f *FeatureCalculator) Run(ctx context.Context, wg sync.WaitGroup) {
	defer wg.Done()

	err := f.processor.Run(ctx)
	if err != nil {
		f.Logger.Error(err.Error())
	}
	f.Logger.Info("window builder shut down nicely")
}
