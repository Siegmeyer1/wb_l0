package utils

import (
	"log"

	"github.com/nats-io/stan.go"
)

func ConnectStan(clientID string) stan.Conn {
	clusterID := "test-cluster"    // nats cluster id
	url := "nats://127.0.0.1:4222" // nats url

	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL(url),
		stan.Pings(1, 3), //set how exacly to ping server. In this case we ping once a second, and if 3 pings fail we consider connection lost
		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) { //what to do if connection is lost
			log.Fatalf("Connection lost, reason: %v", reason)
		}))
	if err != nil {
		log.Fatalf("Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s", err, url)
	}
	log.Println("Connected Nats")
	return sc
}
