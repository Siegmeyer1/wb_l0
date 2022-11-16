package main

import (
	"fmt"
	"os"

	"github.com/Siegmeyer1/wb_l0/utils"
)

func main() {
	Sc := utils.ConnectStan("sender")
	defer Sc.Close()
	var input string
	for {
		fmt.Scanln(&input)
		if input == "exit" {
			break
		}
		file, err := os.ReadFile(input)
		if err != nil {
			fmt.Printf("Can`t read file %s, error: %s\n", input, err)
		} else {
			Sc.Publish("foo", file)
		}
	}
}
