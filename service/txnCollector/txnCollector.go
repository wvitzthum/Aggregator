package transactioncollector

import (
	"context"
	"log"
	"net/url"

	codecs "aggregator/service/codecs"
	models "aggregator/service/models"

	"github.com/go-redis/redis/v9"
	"github.com/gorilla/websocket"
	"github.com/lovoo/goka"
	"go.uber.org/zap"
)

type TxnCollector struct {
	Brokers     []string
	Logger      *zap.Logger
	WindowTopic *models.Topic
	RedisClient *redis.Client
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
	tc.Logger.Info("successfully subscribed", zap.String("url", u.String()))

	// create a new emitter that's going to take the data from the websocket and
	// put it on kafka for us to play with downstream
	emitterWindow, err := goka.NewEmitter(tc.Brokers, *tc.WindowTopic.Stream, tc.WindowTopic.Codec)
	if err != nil {
		log.Fatalf("error creating emitter: %v", err)
	}
	defer emitterWindow.Finish()

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
			// EMIT message to a redis for retention
			data, err := codecs.TxnEncodeGOB(txn)
			if err != nil {
				tc.Logger.Error("error encoding message: %v", zap.String("error", err.Error()))
			}
			err = tc.RedisClient.Set(ctx, txn.X.Hash, data, 0).Err()
			if err != nil {
				tc.Logger.Error("Error settign redis entry", zap.Error(err))
			}
			// Emit a message for each Input of the transaction so we build windows for all addresses involved
			for _, inp := range txn.X.Inputs {
				err = emitterWindow.EmitSync(inp.PrevOut.Addr, txn.X.Hash)
				if err != nil {
					tc.Logger.Error("error emitting message: %v", zap.String("error", err.Error()))
				}
			}

		}
	}
}
