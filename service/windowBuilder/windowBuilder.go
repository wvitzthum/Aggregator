package windowBuilder

import (
	"aggregator/service/models"
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/lovoo/goka"
	"go.uber.org/zap"
)

type WindowState struct {
	g *goka.GroupGraph
}

type WindowBuilder struct {
	Logger      *zap.Logger
	Brokers     []string
	SourceTopic *models.Topic
	OutTopic    *models.Topic
	Processor   *goka.Processor
	Done        chan bool
}

func (wb *WindowBuilder) Init(options ...goka.ProcessorOption) error {
	if wb.Logger == nil || wb.SourceTopic == nil {
		wb.Logger.Fatal("Could not init window builder")
	}

	var err error
	wb.Processor, err = goka.NewProcessor(wb.Brokers,
		goka.DefineGroup("window",
			goka.Input(*wb.SourceTopic.Stream, wb.SourceTopic.Codec, wb.buildWindow),
			goka.Persist(wb.OutTopic.Codec),
		),
		options...,
	)
	if err != nil {
		wb.Logger.Fatal("Could not Initialize Processor", zap.String("Error", err.Error()))
	}

	return nil

}

func (wb *WindowBuilder) Run(ctx context.Context, wg sync.WaitGroup) {
	defer wg.Done()

	err := wb.Processor.Run(ctx)
	if err != nil {
		wb.Logger.Error(err.Error())
	}
	wb.Logger.Info("window builder shut down nicely")
}

func (wb *WindowBuilder) buildWindow(gctx goka.Context, msg interface{}) {
	var window []string
	var ok bool

	// get the existing window against this key
	t := time.Now()
	windowI := gctx.Value()
	if windowI == nil {
		// make a new window
		window = []string{}
	} else {
		window, ok = windowI.([]string)
		if !ok {
			log.Println(windowI)
			gctx.Fail(fmt.Errorf("didn't receive a window from ctx.Value"))
		}
	}
	//log.Println("get", time.Since(t), len(window))

	// assert the msg is an Event
	txn, ok := msg.(string)
	if !ok {
		gctx.Fail(fmt.Errorf("couldn't assert that the received message was of type Event"))
	}

	//t = time.Now()
	// insert the new event into the history ensuring that order is correct
	newWindow := append(window, txn)

	// emit the new window
	gctx.SetValue(newWindow)

	wb.Logger.Info("update window", zap.String("Key", gctx.Key()), zap.Duration("Build time", time.Since(t)), zap.Int("length", len(newWindow)))
}
