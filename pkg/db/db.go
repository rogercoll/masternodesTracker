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
	Coin 			string  `bson:"coin" json:"coin"`
	PublicKey		string	`bson:"publickey" json:"publickey"`
    ApiEndpoint 	string  `bson:"apiendpoint" json:apiendpoint`
	RegexBalance 	string  `bson:"regexbalance" json:regexbalance` //to get the balance in case the API is different
	LastCheck		int64	`bson:"lastcheck" json:lastcheck`
	LastHash		string	`bson:"lasthash" json:lasthash`
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
	//filter bson.D{{"coin", coin}}
	cursor, err := meska.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	if err = cursor.All(context.TODO(), &masternodes); err != nil {
		return nil, err
	}
	return &masternodes, nil
}

func UpdateCoinInfo(c *mongo.Client, publickey string, new *Masternode) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	mdb := c.Database(database)
	meska := mdb.Collection(collection)
	update := bson.M {
		"$set": *new,
	}
	_, err := meska.UpdateOne(ctx, bson.D{{"publickey", publickey}}, update)
	if err != nil {
		return err
	}
	return nil
}
