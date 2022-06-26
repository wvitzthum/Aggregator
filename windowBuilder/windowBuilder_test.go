package windowBuilder

import (
	"aggregator/models"
	"context"
	"math/rand"
	"testing"

	"github.com/google/uuid"
	"github.com/lovoo/goka"
	"github.com/lovoo/goka/tester"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

var (
	windowBuilder WindowBuilder
)

func AssertEqualWindow(t *testing.T, expected, actual []models.Txn) {

	assert.True(t, len(expected) == len(actual))

	for i, e := range expected {
		assert.Equal(t, e.X.Hash, actual[i].X.Hash)
		// comparing times is a pain in the arse
		//assert.True(t, e.T.Equal(actual[i].T))
	}

}

func TestMain(m *testing.M) {
	logger, _ := zap.NewProduction()
	windowBuilder = WindowBuilder{Logger: logger}
}

func TestWindowBuilderOneMessage(t *testing.T) {

	// instantiate the tester
	tester := tester.New(t)
	wb := new(WindowBuilder)


	// run the test processor
	go wb.Processor.Run(context.Background())

	// build a single message
	key := uuid.New().String()
	msg := models.Txn{}

	// get the input (there's only one) and the table rather than hardcoding them
	input := wb.Processor.Graph().InputStreams().Topics()[0]
	table := goka.GroupTable(wb.Processor.Graph().Group())

	// form the expected windoww
	expected := []models.Txn{msg}

	// generate the message on the input
	tester.Consume(input, key, msg)

	// get the result from the table
	actualI := tester.TableValue(table, key)
	actual, ok := actualI.([]models.Txn)
	if !ok {
		t.Fatal("could not assert actual was []Event")
	}

	// assert the result
	AssertEqualWindow(t, expected, actual)
}

func TestIncrementalWindowBuild(t *testing.T) {
	// instantiate the tester
	tester := tester.New(t)

	wb := new(WindowBuilder)

	// build the test processor
	err := wb.Init(nil, goka.WithTester(tester))
	if err != nil {
		t.Fatal(err)
	}
	input := wb.Processor.Graph().InputStreams().Topics()[0]
	table := goka.GroupTable(wb.Processor.Graph().Group())

	// run the test processor
	go wb.Processor.Run(context.Background())

	key := uuid.New().String()
	expected := make([]models.Txn, 0)
	N := 20

	// build up a window message by message, checking validity each step
	for i := 0; i < N; i++ {
		msg := models.Txn{}
		expected = append(expected, msg)
		tester.Consume(input, key, msg)
		actualI := tester.TableValue(table, key)
		actual, ok := actualI.([]models.Txn)
		if !ok {
			t.Fatal("could not assert actual was []Event")
		}
		AssertEqualWindow(t, expected, actual)
	}
}

func TestMultipleWindowBuild(t *testing.T) {
	// instantiate the tester
	tester := tester.New(t)
	wb := new(WindowBuilder)

	// build the test processor
	err := wb.Init(nil, goka.WithTester(tester))
	if err != nil {
		t.Fatal(err)
	}
	input := wb.Processor.Graph().InputStreams().Topics()[0]
	table := goka.GroupTable(wb.Processor.Graph().Group())

	// run the test processor
	go wb.Processor.Run(context.Background())

	N := 2000
	M := 20

	// make M keys
	expected := make(map[string][]models.Txn)
	keys := make([]string, M)
	for i := 0; i < M; i++ {
		key := uuid.New().String()
		keys[i] = key
		w := make([]models.Txn, 0)
		expected[key] = w
	}

	for i := 0; i < N; i++ {
		// pick a random key
		key := keys[rand.Intn(M)]
		msg := models.Txn{}
		newWindow := append(expected[key], msg)
		expected[key] = newWindow
		tester.Consume(input, key, msg)
		// assert as we go
		actualI := tester.TableValue(table, key)
		actual, ok := actualI.([]models.Txn)
		if !ok {
			t.Fatal("could not assert actual was []Event")
		}
		AssertEqualWindow(t, expected[key], actual)
	}

}
