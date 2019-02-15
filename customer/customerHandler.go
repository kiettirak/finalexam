package customer

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kiettirak/finalexam/database"
	_ "github.com/lib/pq"
)

var id int = 0
var customers []Customer
var customersMap = make(map[int]*Customer)

func getCustomersHandler(c *gin.Context) {
	pStatus := c.Query("status")
	var temp []Customer
	sql := "SELECT id, name, email, status FROM customers"
	if pStatus != "" {
		sql += " WHERE status=$1"
	}

	stmt, err := database.Conn().Prepare(sql)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message prepare error": err.Error()})
		return
	}
	rows, err := stmt.Query()
	if pStatus != "" {
		rows, err = stmt.Query(pStatus)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message query error": err.Error()})
		return
	}
	for rows.Next() {
		t := Customer{}
		err := rows.Scan(&t.ID, &t.Name, &t.Email, &t.Status)
		if err != nil {
			log.Fatal("can't scan : ", err)
		}
		temp = append(temp, t)
	}
	c.JSON(http.StatusOK, temp)
}

func getCustomerByIdHandler(c *gin.Context) {
	pId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusProcessing, err.Error())
	}
	stmt, err := database.Conn().Prepare("SELECT id, name, email, status FROM customers WHERE id=$1")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message prepare error": err.Error()})
		return
	}
	row := stmt.QueryRow(pId)
	t := Customer{}
	err = row.Scan(&t.ID, &t.Name, &t.Email, &t.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message scan error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, t)

}

func createCustomersHandler(c *gin.Context) {
	var item Customer
	err := c.ShouldBindJSON(&item)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	//row := database.Conn().QueryRow("INSERT INTO todos (title,status) values ($1,$2) RETURNING id", item.Title, "active")
	var id int

	row := database.InsertCustomer(item.Name, item.Email, item.Status)
	err = row.Scan(&id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	item.ID = id
	c.JSON(http.StatusCreated, item)
}

func updateCustomerHandler(c *gin.Context) {
	pId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusProcessing, err.Error())
	}
	var temp Customer
	err = c.ShouldBindJSON(&temp)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	fmt.Println("DEBUG>>>>", temp)
	stmt, err := database.Conn().Prepare("UPDATE customers SET name=$2, email=$3, status=$4 WHERE id=$1;")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	_, err = stmt.Exec(pId, &temp.Name, &temp.Email, &temp.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	temp.ID = pId
	c.JSON(http.StatusOK, temp)
}

func deleteCustomerHandler(c *gin.Context) {
	pId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusProcessing, err.Error())
	}

	fmt.Println("id: ", pId)
	stmt, err := database.Conn().Prepare("DELETE FROM customers WHERE id=$1;")
	if err != nil {
		fmt.Println("Prepare error: ",err.Error())
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	_, err = stmt.Exec(pId)
	if err != nil {
		fmt.Println("Exec error: ")
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "customer deleted"})
	return
}

// func deleteTodosHandler2(c *gin.Context) {
// 	id, _ := strconv.Atoi(c.Param("id"))

// 	stmt, err := database.Conn().Prepare("delete from todos where id=$1")
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"message statement error": err.Error()})
// 		return
// 	}

// 	if _, err = stmt.Exec(id); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"message execute error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"status": "success"})
// 	return
// }

func CreateTb() {

	database.Conn()

	createTb := `
		CREATE TABLE IF NOT EXISTS customers (
		id SERIAL PRIMARY KEY,
		name TEXT,
		email TEXT,
		status TEXT
		);
		`
	_, err := database.Conn().Exec(createTb)
	if err != nil {
		log.Fatal("can't create table : ", err)
	}
}

func loginMiddleware(c *gin.Context) {
	log.Println("Starting Middleware")

	authKey := c.GetHeader("Authorization")

	if authKey != "token2019" {
		c.JSON(http.StatusUnauthorized, "Unauthorized")
	}

	log.Println("Ending Middleware")
}

func Router() *gin.Engine {
	r := gin.Default()
	v1 := r.Group("")
	v1.Use(loginMiddleware)
	v1.POST("/customers", createCustomersHandler)
	v1.GET("/customers", getCustomersHandler)
	v1.GET("/customers/:id", getCustomerByIdHandler)
	v1.PUT("/customers/:id", updateCustomerHandler)
	v1.DELETE("/customers/:id", deleteCustomerHandler)
	return r
}
