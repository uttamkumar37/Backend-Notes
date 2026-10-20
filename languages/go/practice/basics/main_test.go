package main

import "testing"

func TestWelcomeMessageUsesArguments(t *testing.T) {
	got := welcomeMessage("Java", "Investor Platform")
	want := "Learning Java for Investor Platform."

	if got != want {
		t.Errorf("welcomeMessage() = %q; want %q", got, want)
	}
}

func Example() {
	main()
	// Output:
	// Learning Go for backend engineering.
	// Next: practice variables, functions, structs, and errors.
}
