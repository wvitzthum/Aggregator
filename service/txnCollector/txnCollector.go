package transactioncollector

import (
	"context"
	"log"
	"net/url"

	models "aggregator/service/models"

	"github.com/gorilla/websocket"
	"github.com/lovoo/goka"
)

type TxnCollector struct {
	Brokers	[]string
	WindowTopic	*models.Topic
	TxnTopic	*models.Topic
}

// runBTCCollector reads from the blockchain.info websocket and writes to the
// local Kafka's BTC topic
func (tc *TxnCollector) RunBTCCollector(ctx context.Context) {

	// pick up transactions from blockchain.info
	u, _ := url.Parse("wss://ws.blockchain.info/inv")
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal(err)
	}

	// send a subscription message to start the websocket's data
	hi := []byte("{\"op\":\"unconfirmed_sub\"}")
	defer c.Close()
	err = c.WriteMessage(websocket.TextMessage, hi)
	if err != nil {
		log.Println(err)
	}
	log.Println("successfully subscribed to", u.String())

	// create a new emitter that's going to take the data from the websocket and
	// put it on kafka for us to play with downstream
	emitterWindow, err := goka.NewEmitter(tc.Brokers, *tc.WindowTopic.Stream, tc.WindowTopic.Codec)
	if err != nil {
		log.Fatalf("error creating emitter: %v", err)
	}
	defer emitterWindow.Finish()

	emitterTXN, err := goka.NewEmitter(tc.Brokers, *tc.TxnTopic.Stream, tc.TxnTopic.Codec)
	if err != nil {
		log.Fatalf("error creating emitter: %v", err)
	}
	defer emitterTXN.Finish()

	// txnChan is where we're going to put the parsed JSON message from the
	// websocket
	txnChan := make(chan *models.Txn)
	var txn *models.Txn

	// this for loop runs for the life of the service
	for {

		// this bit goes and waits for a message on the websocket
		// this is in a little go function so its not blocking when the cancel
		// signal shows up
		go func() {
			msg := new(models.Txn)
			err := c.ReadJSON(msg)
			if err != nil {
				log.Fatal(err)
			}
			txnChan <- msg
		}()

		// this select either waits for the txn to show up on txnChan or for the
		// cancel signal to show up via the ctx.Done channel.
		select {
		case <-ctx.Done():
			log.Println("shutting down cleanly")
			return
		case txn = <-txnChan:
			// EMIT message to a topic for retention
			emitterTXN.Emit(txn.X.Hash, txn)
			// Emit a message for each Input of the transaction so we build windows for all addresses involved
			for _, inp := range(txn.X.Inputs) {
				key := inp.PrevOut.Addr
				err = emitterWindow.EmitSync(key, txn)
				if err != nil {
					log.Fatalf("error emitting message: %v", err)
				}
			}

		}
	}
}