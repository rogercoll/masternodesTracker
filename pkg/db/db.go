package db


import (
	"os"
	"time"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Masternode struct {
    Coin string `bson:"coin" json:"coin"`
    Url string `bson:"url" json:url`
    Regex string `bson:"regex" json:regex`
}

var (
	atlasAPI =	os.Getenv("AtlasAPI")
	database = "masternodes"
	collection = "eska"
)

//CHECK NETWORK SECURITY SETTINGS IN ATLAS PANEL! ADD YOUR IP ADDRESS!!
func NewMongoClient() (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(
		atlasAPI,
	))
	if err != nil { return nil, err }


	return client, nil
}

func NewEntry(c *mongo.Client, masternode *Masternode) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	mdb := c.Database(database)
	meska := mdb.Collection(collection)
	result, err := meska.InsertOne(ctx, *masternode)
	if err != nil {
		return nil, err
	}
	return result.InsertedID, nil
}

func GetCoinInfo(c *mongo.Client, coin string) (*[]Masternode, error) {
	var masternodes []Masternode
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	mdb := c.Database(database)
	meska := mdb.Collection(collection)
	cursor, err := meska.Find(ctx, bson.D{{"coin", coin}})
	if err != nil {
		return nil, err
	}
	if err = cursor.All(context.TODO(), &masternodes); err != nil {
		return nil, err
	}
	return &masternodes, nil
}
