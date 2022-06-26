package windowBuilder

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lovoo/goka"
	"github.com/lovoo/goka/tester"
	"github.com/stretchr/testify/assert"
)

func AssertEqualWindow(t *testing.T, expected, actual []Event) {

	assert.True(t, len(expected) == len(actual))

	for i, e := range expected {
		assert.Equal(t, e.Value, actual[i].Value)
		// comparing times is a pain in the arse
		assert.True(t, e.T.Equal(actual[i].T))
	}

}

func TestWindowBuilderOneMessage(t *testing.T) {

	// instantiate the tester
	tester := tester.New(t)
	wb := new(WindowBuilder)

	// build the test processor
	p, err := wb.Init(nil, goka.WithTester(tester))
	if err != nil {
		t.Fatal(err)

	}

	// run the test processor
	go p.Run(context.Background())

	// build a single message
	key := uuid.New().String()
	msg := Event{
		T:     time.Now(),
		Value: rand.NormFloat64(),
	}

	// get the input (there's only one) and the table rather than hardcoding them
	input := p.Graph().InputStreams().Topics()[0]
	table := goka.GroupTable(p.Graph().Group())

	// form the expected windoww
	expected := []Event{msg}

	// generate the message on the input
	tester.Consume(input, key, msg)

	// get the result from the table
	actualI := tester.TableValue(table, key)
	actual, ok := actualI.([]Event)
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
	p, err := wb.Init(nil, goka.WithTester(tester))
	if err != nil {
		t.Fatal(err)
	}
	input := p.Graph().InputStreams().Topics()[0]
	table := goka.GroupTable(p.Graph().Group())

	// run the test processor
	go p.Run(context.Background())

	key := uuid.New().String()
	expected := make([]Event, 0)
	N := 20

	// build up a window message by message, checking validity each step
	for i := 0; i < N; i++ {
		msg := Event{
			T:     time.Now(),
			Value: rand.NormFloat64(),
		}
		expected = append(expected, msg)
		tester.Consume(input, key, msg)
		actualI := tester.TableValue(table, key)
		actual, ok := actualI.([]Event)
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
	p, err := wb.Init(nil, goka.WithTester(tester))
	if err != nil {
		t.Fatal(err)
	}
	input := p.Graph().InputStreams().Topics()[0]
	table := goka.GroupTable(p.Graph().Group())

	// run the test processor
	go p.Run(context.Background())

	N := 2000
	M := 20

	// make M keys
	expected := make(map[string][]Event)
	keys := make([]string, M)
	for i := 0; i < M; i++ {
		key := uuid.New().String()
		keys[i] = key
		w := make([]Event, 0)
		expected[key] = w
	}

	for i := 0; i < N; i++ {
		// pick a random key
		key := keys[rand.Intn(M)]
		msg := Event{
			T:     time.Now(),
			Value: rand.NormFloat64(),
		}
		newWindow := append(expected[key], msg)
		expected[key] = newWindow
		tester.Consume(input, key, msg)
		// assert as we go
		actualI := tester.TableValue(table, key)
		actual, ok := actualI.([]Event)
		if !ok {
			t.Fatal("could not assert actual was []Event")
		}
		AssertEqualWindow(t, expected[key], actual)
	}

}
