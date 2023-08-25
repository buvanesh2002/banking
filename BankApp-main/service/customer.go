package service

import (
	"bankDemo/interfaces"
	"bankDemo/models"
	"context"
	"fmt"
	"log"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type Cust struct {
	ctx             context.Context
	mongoCollection *mongo.Collection
}

func InitCustomer(collection *mongo.Collection, Ctx context.Context) interfaces.Icustomer {
	return &Cust{Ctx, collection}
}
func (c *Cust) CreateCustomer(user *models.Customer) (*mongo.InsertOneResult, error) {
	indexModel := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "account_id", Value: 1}, {Key: "customer_id", Value: 1}}, // 1 for ascending, -1 for descending
			Options: options.Index().SetUnique(true),
		},
	}
	_, err := c.mongoCollection.Indexes().CreateMany(c.ctx, indexModel)
	if err != nil {
		return nil, err
	}
	user.Transaction[0].Date = time.Now()
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 7)
	user.Password = string(hashedPassword)
	res, err := c.mongoCollection.InsertOne(c.ctx, &user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			log.Fatal("Duplicate key error")
		}
		return nil, err
	}

	return res, nil
}

func (c *Cust) GetCustomerById(id int) (*models.Customer, error) {
	filter := bson.D{{Key: "customer_id", Value: id}}
	var customer *models.Customer
	res := c.mongoCollection.FindOne(c.ctx, filter)
	err := res.Decode(&customer)
	if err != nil {
		return nil, err
	}
	fmt.Println(customer)
	return customer, nil
}

func (c *Cust) UpdateCustomerById(id int64, n *models.UpdateModel) (*mongo.UpdateResult, error) {
	iv := bson.M{"customer_id": id}
	if n.Topic == "password" {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(string(n.FinalValue.(string))), 8)
		n.FinalValue = string(hashedPassword)
	}
	if reflect.TypeOf(n.FinalValue).String() == "float64" {
		n.FinalValue = int64(n.FinalValue.(float64))
	}
	fv := bson.M{"$set": bson.M{n.Topic: n.FinalValue}}
	res, err := c.mongoCollection.UpdateOne(c.ctx, iv, fv)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Cust) DeleteCustomerById(id int64) (*mongo.DeleteResult, error) {
	del := bson.M{"customer_id": id}
	res, err := c.mongoCollection.DeleteOne(c.ctx, del)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Cust) GetAllCustomerTransaction(id int64) (*[]models.CustTransaction, error) {
	filter := bson.D{{Key: "customer_id", Value: id}}
	var customer *models.Customer
	res := c.mongoCollection.FindOne(c.ctx, filter)
	err := res.Decode(&customer)
	if err != nil {
		return nil, err
	}
	return &customer.Transaction, nil
}
func (c *Cust) GetCustomerTransactionByDate(fromDate, toDate time.Time) ([]*models.CustTransaction, error) {
	filter := bson.M{
		"transaction.date": bson.M{
			"$gte": fromDate,
			"$lte": toDate,
		},
	}

	options := options.Find()
	//options.SetSort(bson.D{{"transaction.date", 1}}) // Sort the results if needed

	cursor, err := c.mongoCollection.Find(c.ctx, filter, options)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(c.ctx)

	var transactions []*models.CustTransaction
	for cursor.Next(c.ctx) {
		var customer models.Customer
		if err := cursor.Decode(&customer); err != nil {
			return nil, err
		}

		// Append each transaction from the customer's transactions to the transactions slice
		for cursor.Next(c.ctx) {
			var customer models.Customer
			if err := cursor.Decode(&customer); err != nil {
				return nil, err
			}

			// Filter and append transactions within the specified date range
			for _, transaction := range customer.Transaction {
				if transaction.Date.After(fromDate) && transaction.Date.Before(toDate) {
					transactions = append(transactions, &transaction)
				}
			}
		}

	}
	return transactions, nil
}

