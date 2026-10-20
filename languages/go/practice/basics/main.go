package main

import "fmt"

func main() {
	language := "Go"
	goal := "backend engineering"
	message := welcomeMessage(language, goal)

	fmt.Println(message)
	printNextSteps()
}

func welcomeMessage(language string, goal string) string {
	return fmt.Sprintf("Learning %s for %s.", language, goal)
}

func printNextSteps() {
	fmt.Println("Next: practice variables, functions, structs, and errors.")
}
