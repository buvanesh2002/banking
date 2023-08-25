package interfaces

import (
	"bankDemo/models"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type Icustomer interface {
	CreateCustomer(*models.Customer)(*mongo.InsertOneResult,error)
	GetCustomerById(int) (*models.Customer, error)
	UpdateCustomerById(int64, *models.UpdateModel) (*mongo.UpdateResult, error)
	DeleteCustomerById(int64) (*mongo.DeleteResult, error)
	GetAllCustomerTransaction(int64) (*[]models.CustTransaction, error)
	GetCustomerTransactionByDate(fromDate, toDate time.Time)([]*models.CustTransaction, error)

}