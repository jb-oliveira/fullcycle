package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type ViaCep struct {
	Cep        string `json:"cep"`
	Street     string `json:"logradouro"`
	Complement string `json:"complemento"`
	Unit       string `json:"unidade"`
	District   string `json:"bairro"`
	City       string `json:"localidade"`
	Uf         string `json:"uf"`
	State      string `json:"estado"`
	Region     string `json:"regiao"`
	Ibge       string `json:"ibge"`
	Gia        string `json:"gia"`
	Ddd        string `json:"ddd"`
	Siafi      string `json:"siafi"`
}

func main() {
	for _, cep := range os.Args[1:] {
		req, err := http.Get("https://viacep.com.br/ws/" + cep + "/json/")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error making http request: %v\n", err)
		}
		defer req.Body.Close()
		res, err := io.ReadAll(req.Body)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading response: %v\n", err)
		}
		var data ViaCep
		err = json.Unmarshal(res, &data)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading json: %v\n", err)
		}
		file, err := os.Create("cep.txt")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating file: %v\n", err)
		}
		defer file.Close()
		_, err = file.WriteString(fmt.Sprintf("UF: %s, City: %s, District: %s, Street: %s", data.Uf, data.City, data.District, data.Street))

	}
}
