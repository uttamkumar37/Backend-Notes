package main

import (
	"fmt"

	"example.com/go-lab/module-basics/greeting"
)

func main() {
	investorName := "Asha"
	message := greeting.Welcome(investorName)
	fmt.Println(message)
}
