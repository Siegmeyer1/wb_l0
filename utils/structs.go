package utils

type Config struct {
	Nats_ struct {
		Url       string `yaml:"url"`
		ClusterID string `yaml:"cluster ID"`
	} `yaml:"nats"`
	Postgr struct {
		User string `yaml:"username"`
		Pass string `yaml:"password"`
		Addr string `yaml:"address"`
		Db   string `yaml:"database"`
	} `yaml:"postgres"`
	Http struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"http"`
}

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
