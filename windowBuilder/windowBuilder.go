package windowBuilder

import (
	"aggregator/codecs"
	"aggregator/models"
	"context"
	"fmt"
	"log"

	"github.com/lovoo/goka"
	"go.uber.org/zap"
)

type WindowState struct {
	g *goka.GroupGraph
}

type WindowBuilder struct {
	Logger			*zap.Logger
	SourceTopic		*models.Topic
	AggTopic		*codecs.ArrayCodec
	Processor		*goka.Processor
	Done			chan bool
}

func (wb *WindowBuilder) Init(brokers []string, options ...goka.ProcessorOption) error {
	if wb.Logger == nil || wb.SourceTopic == nil {
		wb.Logger.Fatal("Could not init windowbuilder")
	}
	
	var err error
	wb.Processor, err = goka.NewProcessor(brokers,
		goka.DefineGroup("window",
			goka.Input(wb.SourceTopic.Stream, wb.SourceTopic.Codec, wb.buildWindow),
			goka.Persist(wb.AggTopic),
		),
		options...,
	)
	if err != nil {
		wb.Logger.Fatal("Could not Initialize Processor", zap.String("Error", err.Error()))
	}

	return nil

}

func (wb *WindowBuilder) Run(ctx context.Context, brokers []string) {
	defer close(wb.Done)

	err := wb.Processor.Run(ctx)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("window builder shut down nicely")
}

func (wb *WindowBuilder) buildWindow(gctx goka.Context, msg interface{}) {
	var window []models.Txn
	var ok bool

	// get the existing window against this key
	//t := time.Now()
	windowI := gctx.Value()
	if windowI == nil {
		// make a new window
		window = make([]models.Txn, 0)
	} else {
		window, ok = windowI.([]models.Txn)
		if !ok {
			log.Println(windowI)
			gctx.Fail(fmt.Errorf("didn't receive a window from ctx.Value"))
		}
	}
	//log.Println("get", time.Since(t), len(window))

	// assert the msg is an Event
	txn, ok := msg.(models.Txn)
	if !ok {
		gctx.Fail(fmt.Errorf("couldn't assert that the received message was of type Event"))
	}

	//t = time.Now()
	// insert the new event into the history ensuring that order is correct
	newWindow := append(window, txn)

	// emit the new window
	gctx.SetValue(newWindow)
	//log.Println("set", time.Since(t), len(newWindow))

}
