package mekano

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Payment struct {
	File        string    `bson:"file"`
	Consecutive int       `bson:"consecutive"`
	CreateAt    time.Time `bson:"create_at"`
}

type Billing struct {
	File     string    `bson:"file" json:"file"`
	Debit    int       `bson:"debit" json:"debit"`
	Credit   int       `bson:"credit" json:"credit"`
	Base     int       `bson:"base" json:"base"`
	CreateAt time.Time `bson:"create_at" json:"create_at"`
}

type MekanoData struct {
	Name string `bson:"name" json:"name"`
	Code string `bson:"code" json:"code"`
}

type DatabaseInterface interface {
	GetPayment() (Payment, error)
	SavePayment(payment Payment) (interface{}, error)
	SaveBilling(billing Billing) (interface{}, error)
	GetAccounts(param string) (MekanoData, error)
	GetCashiers(param string) (MekanoData, error)
	GetCostCenter(param string) (MekanoData, error)
}

type database struct {
	db  *mongo.Client
	ctx context.Context
}

func NewDatabase() (DatabaseInterface, error) {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return nil, err
	}

	return &database{
		db:  client,
		ctx: ctx,
	}, nil
}

func (dr *database) GetPayment() (Payment, error) {
	var result Payment

	collection := dr.db.Database("mekano").Collection("payments")
	option := options.FindOneOptions{}
	option.SetSort(bson.D{{Key: "create_at", Value: -1}})
	err := collection.FindOne(dr.ctx, bson.D{}, &option).Decode(&result)
	if err != nil {
		return Payment{}, nil
	}
	return result, nil
}

func (dr *database) SavePayment(payment Payment) (interface{}, error) {
	collection := dr.db.Database("mekano").Collection("payments")
	cur, err := collection.InsertOne(dr.ctx, payment)
	if err != nil {
		return nil, err
	}
	return cur.InsertedID, nil
}

func (dr *database) SaveBilling(billing Billing) (interface{}, error) {
	collection := dr.db.Database("mekano").Collection("billings")
	cur, err := collection.InsertOne(dr.ctx, billing)
	if err != nil {
		return nil, nil
	}
	return cur.InsertedID, nil
}

func (dr *database) GetAccounts(param string) (MekanoData, error) {
	var results MekanoData
	coll := dr.db.Database("mekano").Collection("accounts")
	filter := bson.D{{Key: "name", Value: param}}
	err := coll.FindOne(dr.ctx, filter).Decode(&results)
	if err != nil {
		return MekanoData{}, err
	}

	return results, nil
}

func (dr *database) GetCashiers(param string) (MekanoData, error) {
	var results MekanoData
	coll := dr.db.Database("mekano").Collection("cashiers")
	filter := bson.D{{Key: "name", Value: param}}
	err := coll.FindOne(dr.ctx, filter).Decode(&results)
	if err != nil {
		return MekanoData{}, err
	}

	return results, nil
}

func (dr *database) GetCostCenter(param string) (MekanoData, error) {
	var results MekanoData
	coll := dr.db.Database("mekano").Collection("cost_center")
	filter := bson.D{{Key: "name", Value: param}}
	err := coll.FindOne(dr.ctx, filter).Decode(&results)
	if err != nil {
		return MekanoData{}, err
	}

	return results, nil
}
