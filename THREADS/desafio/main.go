package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <cep>")
		return
	}

	cep := os.Args[1]

	urlViaCep := "https://viacep.com.br/ws/" + cep + "/json/"
	urlBrasilApi := "https://brasilapi.com.br/api/cep/v1/" + cep
	ch1 := make(chan Cep)
	ch2 := make(chan Cep)

	go loadDataFromUrl(urlViaCep, &ViaCep{}, ch1)
	go loadDataFromUrl(urlBrasilApi, &BrasilApiCep{}, ch2)

	select {
	case cep := <-ch1:
		fmt.Printf("Received ViaCep: %v\n", cep)
	case cep := <-ch2:
		fmt.Printf("Received BrasilApi: %v\n", cep)
	case <-time.After(1 * time.Second):
		println("timeout")
	}
	close(ch1)
	close(ch2)
}

func loadDataFromUrl(url string, conversor ConversorCep, ch chan<- Cep) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Default().Printf("Error creating request: %v\n", url)
		log.Default().Println(err)
		return
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Default().Printf("Error making request: %v\n", url)
		log.Default().Println(err)
		return
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Default().Printf("Error reading response body: %v\n", url)
		log.Default().Println(err)
		return
	}
	err = json.Unmarshal(body, &conversor)
	if err != nil {
		log.Default().Printf("Error parsing JSON: %v\n", string(body))
		log.Default().Println(err)
		return
	}
	ch <- conversor.ToCep()
}

type ViaCep struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

type BrasilApiCep struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}

type Cep struct {
	Cep        string `json:"cep"`
	Logradouro string `json:"logradouro"`
	Bairro     string `json:"bairro"`
	Cidade     string `json:"cidade"`
	Uf         string `json:"uf"`
	Api        string `json:"api"`
}

type ConversorCep interface {
	ToCep() Cep
}

func (v *ViaCep) ToCep() Cep {
	return Cep{
		Cep:        v.Cep,
		Logradouro: v.Logradouro,
		Bairro:     v.Bairro,
		Cidade:     v.Localidade,
		Uf:         v.Uf,
		Api:        "ViaCep",
	}
}

func (b *BrasilApiCep) ToCep() Cep {
	return Cep{
		Cep:        b.Cep,
		Logradouro: b.Street,
		Bairro:     b.Neighborhood,
		Cidade:     b.City,
		Uf:         b.State,
		Api:        "BrasilApi",
	}
}
