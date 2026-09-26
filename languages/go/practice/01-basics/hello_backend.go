package main

import "fmt"

func main() {
	language := "Go"
	goal := "backend engineering"

	printWelcome(language, goal)
	printNextSteps()
}

func printWelcome(language string, goal string) {
	fmt.Printf("Learning %s for %s.\n", language, goal)
}

func printNextSteps() {
	fmt.Println("Next: practice variables, functions, structs, and errors.")
}

