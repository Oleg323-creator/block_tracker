package main

import (
	"block_tracker/internal/db"
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
)

func main() {
	cfg := db.ConnectionConfig{
		Host:     "localhost",
		Port:     "5430",
		Username: "postgres",
		Password: "postgres",
		DBName:   "postgres",
		SSLMode:  "disable",
	}

	// DB CONNECT
	dbConn := db.NewDB(cfg)

	// RUN MIGRATIONS
	err := db.RunMigrations(dbConn)
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	//MIGRATIONS CONNECT
	repo := db.NewRepository(dbConn)

	client, err := ethclient.Dial("wss://sepolia.infura.io/ws/v3/8ab0a7925db44ab094fa1b6546409a6b")
	if err != nil {
		log.Fatalf("Failed to connect to Ethereum client: %v", err)
	}
	defer client.Close()

	// CHANNEL FOR BLOCKS LISTENING
	headers := make(chan *types.Header)
	sub, err := client.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		log.Fatalf("Failed to subscribe to new headers: %v", err)
	}

	// LISTENING
	fmt.Println("Listening for new blocks on Sepolia...")
	for {
		select {
		case err = <-sub.Err():
			log.Fatalf("Subscription error: %v", err)
		case header := <-headers:
			fmt.Printf("New block mined! Block number: %d\n", header.Number.Uint64())

			err := repo.SaveBlocksToDB(header.Number.Int64())
			if err != nil {
				log.Printf("Failed to save block number %d: %v\n", header.Number.Int64(), err)
				continue
			}
		}
	}
}
