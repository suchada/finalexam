package main

import(
	"github.com/suchada/finalexam/customer"
)



func main() {
	r := customer.SetupRouter()
	r.Run(":2019")
}
