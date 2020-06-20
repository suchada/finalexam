package customer

import (
	"fmt"
	// "log"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/suchada/finalexam/database"
)

type Customer struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Status string `json:"status"`
}

func createCustomerHandler(c *gin.Context) {
	cust := Customer{}
	if error := c.ShouldBindJSON(&cust); error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": error.Error()})
		return
	}

	row := database.Conn().QueryRow("INSERT INTO customer (name, email, status) values ($1, $2, $3)  RETURNING id, name, email, status", cust.Name, cust.Email, cust.Status)

	error := row.Scan(&cust.ID, &cust.Name, &cust.Email, &cust.Status)
	if error != nil {
		c.JSON(http.StatusInternalServerError, error)
		return
	}

	c.JSON(http.StatusCreated, cust)
}

func getAllCustomerHandler(c *gin.Context) {

	stmt, error := database.Conn().Prepare("SELECT id, name, email, status FROM customer")
	if error != nil {
		c.JSON(http.StatusInternalServerError, error)
		return
	}

	rows, error := stmt.Query()
	if error != nil {
		c.JSON(http.StatusInternalServerError, error)
		return
	}

	customers := []Customer{}
	for rows.Next() {
		cust := Customer{}

		error := rows.Scan(&cust.ID, &cust.Name, &cust.Email, &cust.Status)
		if error != nil {
			c.JSON(http.StatusInternalServerError, error)
			return
		}

		customers = append(customers, cust)
	}

	c.JSON(http.StatusOK, customers)
}

func getCustomerByIdHandler(c *gin.Context) {
	id := c.Param("id")
	stmt, error := database.Conn().Prepare("SELECT id, name, email, status FROM customer where id=$1")
	if error != nil {
		c.JSON(http.StatusInternalServerError, error)
		fmt.Printf("1: %w", error)
		return
	}

	row := stmt.QueryRow(id)

	cust := &Customer{}

	error = row.Scan(&cust.ID, &cust.Name, &cust.Email, &cust.Status)
	if error != nil {
		if error == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, error)
			return
		}
		c.JSON(http.StatusInternalServerError, error)
		return
	}

	c.JSON(http.StatusOK, cust)
}

func updateCustomerByIdHandler(c *gin.Context) {
	id := c.Param("id")
	// Query
	stmt, error := database.Conn().Prepare("SELECT id, name, email, status FROM customer where id=$1")
	if error != nil {
		c.JSON(http.StatusInternalServerError, error)
		return
	}

	row := stmt.QueryRow(id)

	cust := &Customer{}

	error = row.Scan(&cust.ID, &cust.Name, &cust.Email, &cust.Status)
	if error != nil {
		c.JSON(http.StatusInternalServerError, error)
		return
	}

	// Assign from request
	if error := c.ShouldBindJSON(cust); error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": error.Error()})
		return
	}

	stmt, error = database.Conn().Prepare("UPDATE customer SET name=$2, email=$3, status=$4 WHERE id=$1;")
	if error != nil {
		c.JSON(http.StatusInternalServerError, error)
		return
	}

	if _, error := stmt.Exec(id, &cust.Name, &cust.Email, &cust.Status); error != nil {
		c.JSON(http.StatusInternalServerError, error)
		return
	}

	c.JSON(http.StatusOK, cust)
}

func deleteCustomerByIdHandler(c *gin.Context) {
	id := c.Param("id")

	stmt, error := database.Conn().Prepare("DELETE FROM customer WHERE id = $1")
	if error != nil {
		c.JSON(http.StatusInternalServerError, error)
		return
	}

	if _, error := stmt.Exec(id); error != nil {
		c.JSON(http.StatusInternalServerError, error)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "customer deleted"})
}

func authMDW(c *gin.Context) {
	fmt.Println("start #middleware")
	token := c.GetHeader("Authorization")
	if token == "token2019wrong_token" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no permission to access the application"})
		c.Abort()
		return
	}
	c.Next()
}

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(authMDW)
	r.POST("/customers", createCustomerHandler)
	r.GET("/customers/:id", getCustomerByIdHandler)
	r.GET("/customers", getAllCustomerHandler)
	r.PUT("/customers/:id", updateCustomerByIdHandler)
	r.DELETE("/customers/:id", deleteCustomerByIdHandler)
	return r
}
