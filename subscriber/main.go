package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/Siegmeyer1/wb_l0/utils"
	_ "github.com/lib/pq"
	"github.com/nats-io/stan.go"
)

var cash = make(map[string]Order)
var db *sql.DB

type Delivery struct {
	Name    string
	Phone   string
	Zip     string
	City    string
	Address string
	Region  string
	Email   string
}

type Payment struct {
	Transaction   string
	Request_id    string
	Currency      string
	Provider      string
	Amount        float32
	Payment_dt    int32
	Bank          string
	Delivery_cost float32
	Goods_total   float32
	Custom_fee    float32
}

type Item struct {
	Chrt_id      int
	Track_number string
	Price        float32
	Rid          string
	Name         string
	Sale         float32
	Size         string
	Total_price  float32
	Nm_id        int
	Brand        string
	Status       int
}

type Order struct {
	Order_uid          string
	Track_number       string
	Entry              string
	Delivery           Delivery
	Payment            Payment
	Items              []Item
	Locale             string
	Internal_signature string
	Customer_id        string
	Delivery_service   string
	Shardkey           string
	Sm_id              int
	Date_created       string
	Oof_shard          string
}

func errHandle(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getOrderByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Got request")
	if r.Method == "POST" {
		io.WriteString(w, "test")
		reqBody, err := ioutil.ReadAll(r.Body)
		errHandle(err)
		fmt.Println(string(reqBody))
		response, err := json.Marshal(cash[string(reqBody)])
		errHandle(err)
		_, err = io.WriteString(w, string(response))
		errHandle(err)
	}
	if r.Method == "GET" {
		_, err := io.WriteString(w, "Hello! You can get order info in .json format. To get it,\n"+
			"send POST request with desired order_uid in request body.")
		errHandle(err)
	}
}

func init() {
	var err error

	// connecting to Postgres db
	connStr := "postgres://wb_intern:1029@localhost/wb_intern?sslmode=disable" //TODO: learn some way to NOT hardcode passwords
	db, err = sql.Open("postgres", connStr)
	errHandle(err)
	err = db.Ping()
	errHandle(err)
	fmt.Println("Successfully connected to Postgres")

	//loading orders from db to memory
	rows, err := db.Query(`SELECT * FROM "order"`)
	errHandle(err)
	defer rows.Close()
	for rows.Next() {
		var order_uid string
		var data []byte
		var unm_data Order

		err = rows.Scan(&order_uid, &data)
		errHandle(err)

		json.Unmarshal(data, &unm_data)
		cash[order_uid] = unm_data
	}
}

func main() {
	Sc := utils.ConnectStan("subscriber")
	defer Sc.Close()

	dealWithMsg := func(msg *stan.Msg) {
		log.Printf("Received a message: %s\n", string(msg.Data))

		reader := bytes.NewReader(msg.Data) //configuring .json unmarshalling
		decoder := json.NewDecoder(reader)
		decoder.DisallowUnknownFields()

		var order Order
		err := decoder.Decode(&order) //fill Order var with .json data
		if err != nil {
			log.Println(err)

		} else if order.Order_uid == "" { // here we can add more specifics on how full .json should be
			log.Println("Data in .json is not full enough")

		} else { //adding order to cash and inserting into pg
			cash[order.Order_uid] = order
			sql := `INSERT INTO "order"("order_uid", "data") VALUES($1, $2)`
			_, err = db.Exec(sql, order.Order_uid, string(msg.Data))
			errHandle(err)
		}
	}

	Sc.Subscribe("JSON_channel", dealWithMsg)
	http.HandleFunc("/", getOrderByID)
	err := http.ListenAndServe(":3333", nil)
	errHandle(err)
}
