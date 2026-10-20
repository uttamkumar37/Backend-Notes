package greeting_test

import (
	"testing"

	"example.com/backend-notes/practice/packages/greeting"
)

func TestWelcomeUsesInvestorName(t *testing.T) {
	got := greeting.Welcome("Mira")
	if got != "Welcome, Mira." {
		t.Errorf("Welcome() = %q; want %q", got, "Welcome, Mira.")
	}
}
