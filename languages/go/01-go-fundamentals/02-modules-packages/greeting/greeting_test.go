package greeting_test

import (
	"testing"

	"example.com/go-lab/module-basics/greeting"
	"github.com/google/go-cmp/cmp"
)

func TestWelcomeUsesInvestorName(t *testing.T) {
	got := greeting.Welcome("Mira")
	want := "Welcome, Mira."
	diff := cmp.Diff(want, got)

	if diff != "" {
		t.Errorf("Welcome() mismatch (-want +got):\n%s", diff)
	}
}
