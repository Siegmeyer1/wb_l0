package utils

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-yaml/yaml"
	"github.com/nats-io/stan.go"
)

func NewConfig(configPath string) Config {
	config := Config{}
	file, err := os.Open(configPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	d := yaml.NewDecoder(file)
	if err := d.Decode(&config); err != nil {
		log.Fatal(err)
	}
	return config
}

var Cfg = NewConfig("config.yml")

func ErrHandle(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func ConnectStan(clientID string) stan.Conn {
	clusterID := Cfg.Nats_.ClusterID // nats cluster id
	url := Cfg.Nats_.Url             // nats url

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

func ConnectPG() *sql.DB {
	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", Cfg.Postgr.User, Cfg.Postgr.Pass, Cfg.Postgr.Addr, Cfg.Postgr.Db)
	db, err := sql.Open("postgres", connStr)
	ErrHandle(err)
	ErrHandle(db.Ping())
	fmt.Println("Successfully connected to Postgres")
	return db
}
