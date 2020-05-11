package main

import (
	"log"
	"github.com/rogercoll/masternodesTracker/pkg/db"
)


func main() {
	c, err := db.NewMongoClient()
	if err != nil { log.Fatal(err)}
	mcoins, err := db.GetCoinInfo(c, "eska")
	if err != nil { log.Fatal(err)}
	log.Printf("Your masternodes: %v", *mcoins)
}