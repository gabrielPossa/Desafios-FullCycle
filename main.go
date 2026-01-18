package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

type CEP struct {
	Cep        string
	Logradouro string
	Bairro     string
	Localidade string
	Estado     string
}

type fetchResponse struct {
	Cep     CEP
	apiName string
	err     error
}

type apiResponse struct {
	Cep          string `json:"cep"`
	Logradouro   string `json:"logradouro"`
	Complemento  string `json:"complemento"`
	Unidade      string `json:"unidade"`
	Bairro       string `json:"bairro"`
	Localidade   string `json:"localidade"`
	Uf           string `json:"uf"`
	Estado       string `json:"estado"`
	Regiao       string `json:"regiao"`
	Ibge         string `json:"ibge"`
	Gia          string `json:"gia"`
	Ddd          string `json:"ddd"`
	Siafi        string `json:"siafi"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}

var digitCheck = regexp.MustCompile("^\\d{8}$")

func main() {

	responseCH := make(chan fetchResponse)

	if len(os.Args) != 2 {
		log.Fatal("Uso: go run main.go CEP")
	}

	cep := os.Args[1]

	cep = strings.Replace(cep, "-", "", -1)

	if !digitCheck.MatchString(cep) {
		log.Fatal("CEP invalido,CEP deve ser composto por 8 números. Formatos aceitos: 12345678 ou 12345-678 ")
	}

	ctx, cancelCTX := context.WithTimeout(context.Background(), time.Second)
	defer cancelCTX()

	go fetchCEPData(ctx, fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep), "viacep", responseCH)
	go fetchCEPData(ctx, fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep), "brasilapi", responseCH)

	for {
		select {
		case res := <-responseCH:
			if res.err != nil {
				continue
			}
			log.Printf("\n\tAPI: %s\n\tCEP: %s\n\tLogradouro: %s\n\tBairro: %s\n\tCidade: %s\n\tEstado: %s\n",
				res.apiName, res.Cep.Cep, res.Cep.Logradouro, res.Cep.Bairro, res.Cep.Localidade, res.Cep.Estado)
			return
		case <-ctx.Done():
			log.Println("Timeout")
			return
		}
	}
}

func fetchCEPData(ctx context.Context, url, apiName string, ch chan<- fetchResponse) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		ch <- fetchResponse{err: err}
		return
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		ch <- fetchResponse{err: err}
		return
	}
	defer res.Body.Close()

	var data apiResponse
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		ch <- fetchResponse{err: err}
		return
	}

	var cep CEP
	switch apiName {
	case "brasilapi":
		cep = CEP{
			Cep:        data.Cep,
			Logradouro: data.Street,
			Bairro:     data.Neighborhood,
			Localidade: data.City,
			Estado:     data.State,
		}

	case "viacep":
		cep = CEP{
			Cep:        data.Cep,
			Logradouro: data.Logradouro,
			Bairro:     data.Bairro,
			Localidade: data.Localidade,
			Estado:     data.Uf,
		}
	}

	ch <- fetchResponse{
		apiName: apiName,
		err:     nil,
		Cep:     cep,
	}
}
