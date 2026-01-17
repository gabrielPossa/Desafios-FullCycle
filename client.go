package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		log.Fatal(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode != 200 {
		log.Printf("status code error: %d - %s", res.StatusCode, string(bodyBytes))
		return
	}

	var data map[string]string
	err = json.Unmarshal(bodyBytes, &data)
	if err != nil {
		log.Println(err)
	}

	err = os.WriteFile("cotacao.txt", []byte(fmt.Sprintf("Dólar: %s", data["Dólar"])), 0644)
	if err != nil {
		log.Fatal(err) // Handle any potential errors
	}

	log.Printf("Cotação salva em ./cotacao.txt. Valor recebido: %s", data["Dólar"])
}
