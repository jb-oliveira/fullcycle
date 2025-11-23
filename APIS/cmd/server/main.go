package main

import (
	"log"

	"github.com/jb-oliveira/fullcycle/tree/main/APIS/configs"
)

func main() {
	_, err := configs.LoadWebConfig(".")
	if err != nil {
		log.Fatalf("falha ao carregar configuração web: %v", err)
	}
	InitDB()
	log.Println("Configuração carregada com sucesso")
}

func InitDB() {
	_, err := configs.LoadDbConfig(".")
	if err != nil {
		log.Fatalf("falha ao carregar configuração do banco: %v", err)
	}
	err = configs.InitGorm()
	if err != nil {
		log.Fatalf("falha ao conectar ao banco: %v", err)
	}

	db := configs.GetDB()
	if db == nil {
		log.Fatal("instância do banco é nula")
	}
	log.Println("Conexão com banco estabelecida")
}
