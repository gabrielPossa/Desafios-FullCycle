package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	log.Println("Starting server...")

	log.Println("Opening database...")
	db, err := sql.Open("sqlite3", "db.sqlite")
	if err != nil {
		log.Println("Failed to open database")
		panic(err)
	}
	defer db.Close()

	log.Println("Preparing database...")
	err = prepareDB(db)
	if err != nil {
		log.Println("Failed to prepare database")
		panic(err)
	}

	log.Println("Configuring server...")
	http.HandleFunc("/cotacao", handler(db))

	log.Println("Listening on port 8080...")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Println("Failed to start server on port 8080")
		panic(err)
	}
}

func handler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		log.Println("Request from", r.RemoteAddr, "received")
		defer log.Println("Request from", r.RemoteAddr, " finished")

		dollar, err := getDollarBID(ctx)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Println(dollar)

		err = saveDollarBID(ctx, db, dollar)
		if err != nil {
			log.Println("Failed to save data into the DB: ", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]string{"Dólar": dollar})
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func getDollarBID(ctx context.Context) (string, error) {
	reqCTX, reqCCTX := context.WithTimeout(ctx, 200*time.Millisecond)
	defer reqCCTX()
	req, err := http.NewRequestWithContext(reqCTX, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		return "", err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	var data map[string]map[string]string
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return "", err
	}

	return data["USDBRL"]["bid"], nil
}

func saveDollarBID(ctx context.Context, db *sql.DB, bid string) error {
	reqCTX, reqCCTX := context.WithTimeout(ctx, 10*time.Millisecond)
	defer reqCCTX()

	stmt, err := db.PrepareContext(reqCTX, "insert into cotacao(value) values(?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(reqCTX, bid)
	if err != nil {
		return err
	}

	return nil
}

func prepareDB(db *sql.DB) error {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS `cotacao` (`id` INTEGER PRIMARY KEY AUTOINCREMENT, `value` VARCHAR(64))")
	if err != nil {
		return err
	}

	return nil
}
