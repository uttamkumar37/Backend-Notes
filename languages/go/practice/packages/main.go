package main

import (
	"fmt"

	"example.com/backend-notes/practice/packages/greeting"
)

func main() {
	investorName := "Asha"
	message := greeting.Welcome(investorName)
	fmt.Println(message)
}
