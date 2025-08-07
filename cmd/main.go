package cmd

import (
	"github.com/k0haku1/order-service/internal/db"
	"log"
)

func main() {
	database, err := db.InitDB()
	if err != nil {
		log.Fatal(err)
	}
}
