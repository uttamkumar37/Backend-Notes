package greeting

import "fmt"

// Welcome returns a greeting for the supplied fictional investor name.
func Welcome(name string) string {
	return fmt.Sprintf("Welcome, %s.", name)
}
