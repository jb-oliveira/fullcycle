package main

import (
	"log"

	"github.com/jb-oliveira/fullcycle/tree/main/APIS/configs"
)

func main() {
	_, err := configs.LoadDbConfig(".")
	if err != nil {
		log.Fatalf("falha ao carregar configuração do banco: %v", err)
	}

	_, err = configs.LoadWebConfig(".")
	if err != nil {
		log.Fatalf("falha ao carregar configuração web: %v", err)
	}

	dsn, err := configs.GetDSN()
	if err != nil {
		log.Printf("aviso: não foi possível obter DSN: %v", err)
	} else {
		log.Printf("DSN do banco configurado: %s", dsn)
	}

	_, err = configs.NewDB()
	if err != nil {
		log.Fatalf("falha ao conectar ao banco: %v", err)
	}
	log.Println("Conexão com banco estabelecida")

	db := configs.GetDB()
	if db == nil {
		log.Fatal("instância do banco é nula")
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("falha ao obter instância do banco: %v", err)
	}
	defer sqlDB.Close()

	log.Println("Configuração carregada com sucesso")
}
