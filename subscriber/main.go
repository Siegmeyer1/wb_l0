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

var cash = make(map[string]utils.Order)
var db *sql.DB

func getOrderByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Got request")
	switch r.Method {
	case "POST":
		reqBody, err := ioutil.ReadAll(r.Body)
		utils.ErrHandle(err)
		fmt.Println(string(reqBody))
		if val, ok := cash[string(reqBody)]; ok {
			response, err := json.Marshal(val)
			utils.ErrHandle(err)
			if _, err := io.WriteString(w, string(response)); err != nil {
				utils.ErrHandle(err)
			}
		} else if _, err := io.WriteString(w, "Requested order ID not found"); err != nil {
			utils.ErrHandle(err)
		}

	case "GET":
		if _, err := io.WriteString(w, "Hello! You can get order info in .json format. To get it,\n"+
			"send POST request with desired order_uid in request body."); err != nil {
			utils.ErrHandle(err)
		}
	}
}

func init() {
	db = utils.ConnectPG()

	//loading orders from db to memory
	rows, err := db.Query(`SELECT * FROM "order"`)
	utils.ErrHandle(err)
	defer rows.Close()
	for rows.Next() {
		var order_uid string
		var data []byte
		var unm_data utils.Order

		if err := rows.Scan(&order_uid, &data); err != nil {
			utils.ErrHandle(err)
		}
		if err := json.Unmarshal(data, &unm_data); err != nil {
			utils.ErrHandle(err)
		}
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

		var order utils.Order

		err := decoder.Decode(&order)
		switch {
		case err != nil:
			log.Println(err)
		case order.Order_uid == "": // here we can add more specifics on how full .json should be
			log.Println("Data in .json is not full enough")

		default: //adding order to cash and inserting into pg
			cash[order.Order_uid] = order
			sql := `INSERT INTO "order"("order_uid", "data") VALUES($1, $2)`
			if _, err := db.Exec(sql, order.Order_uid, string(msg.Data)); err != nil {
				utils.ErrHandle(err)
			}
		}
	}

	Sc.Subscribe("JSON_channel", dealWithMsg)
	http.HandleFunc("/", getOrderByID)
	addr := fmt.Sprintf("%s:%s", utils.Cfg.Http.Host, utils.Cfg.Http.Port)
	err := http.ListenAndServe(addr, nil)
	utils.ErrHandle(err)
}
