package featureCalculator

import (
	"aggregator/service/featureCalculation"
	"aggregator/service/models"
	"context"
	"fmt"
	"sync"

	"github.com/lovoo/goka"
	"go.uber.org/zap"
)

type FeatureCalculator struct {
	Logger 		*zap.Logger
	processor 	*goka.Processor
	Brokers		[]string
	SourceTopic *models.Topic
	OutTopic 	*models.Topic
}

func (f *FeatureCalculator) outboundStatsProcessor(ctx goka.Context, msg interface{}) {
	window, ok := msg.(([]models.Txn))
	if !ok {
		ctx.Fail(fmt.Errorf("couldn't convert value to Window"))
	}

	features := featureCalculation.CalcFeatures(window)

	// emit new statistics without changing the key
	ctx.Emit(*f.OutTopic.Stream, ctx.Key(), features)
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
