package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Conta struct {
	Numero int     `json:"num"`
	Saldo  float64 `json:"sld"`
}

func main() {

	conta := Conta{Numero: 1, Saldo: 100}
	res, err := json.Marshal(conta)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Conta : %v\n", string(res))

	encoder := json.NewEncoder(os.Stdout)
	err = encoder.Encode(conta)
	if err != nil {
		fmt.Println(err)
	}

	jsonPuro := []byte(`{"num":2, "sld":2100}`)
	var conta2 Conta
	err = json.Unmarshal(jsonPuro, &conta2)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(conta2)
}
