package main

import (
	"fmt"

	"github.com/jb-oliveira/fullcycle/tree/main/curso-go/9-APIS/configs"
)

func main() {
	dbConf, _ := configs.LoadDbConfig(".")
	fmt.Printf("DBConfig: %v", &dbConf)
	webConf, _ := configs.LoadWebConfig(".")
	fmt.Printf("WebConfig: %v", &webConf)
}
