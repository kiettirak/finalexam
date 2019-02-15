package main

import (
	"github.com/kiettirak/finalexam/customer"
)

func main() {

	customer.CreateTb()
	r := customer.Router()
	r.Run(":2019")
}
