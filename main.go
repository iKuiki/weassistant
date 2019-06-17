package main

import (
	"weassistant/conf"
)

func main() {
	config := conf.MustNewConfig()
	err := config.Load("config.json")
	if err != nil {
		panic(err)
	}
}
