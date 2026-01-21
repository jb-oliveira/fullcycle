package main

import (
	"encoding/json"
	"io"
	"net/http"
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
	http.HandleFunc("/", SearchCepHandler)
	http.ListenAndServe(":8080", nil)
}

func SearchCepHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		// http.NotFound(w, r)
		return
	}
	cepParam := r.URL.Query().Get("cep")
	if cepParam == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	viaCep, err := SearchCep(cepParam)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// result, err := json.Marshal(viaCep)
	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }
	// w.Write(result)

	json.NewEncoder(w).Encode(viaCep)
}

func SearchCep(cep string) (*ViaCep, error) {
	res, err := http.Get("https://viacep.com.br/ws/" + cep + "/json/")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var result ViaCep
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
