package routes

import (
	"bankDemo/controllers"

	"github.com/gin-gonic/gin"
)

func CustRoute(router *gin.Engine, controller controllers.TransactionController) {
	router.POST("/customer", controller.CreateCustomer)
	router.GET("/customer/:id", controller.GetCustomerById)
	router.PUT("/customer/:id", controller.UpdateCustomerById)
	router.DELETE("/customer/:id", controller.DeleteCustomerById)
	//router.GET("/customertrans/:id", controller.GetAllCustomerTransaction)
	router.POST("/customer/date",controller.GetCustomerTransactionByDate)
	//router.POST("/customer/:id",controller.GetCustomerTransactionById)
}