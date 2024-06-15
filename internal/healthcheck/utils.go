package healthcheck

import "fmt"

func formatHealthError(key, value, expectedValue string) string {
	return fmt.Sprintf("%s returned %s, while %s was expected", key, value, expectedValue)
}

func formatKeyError(key string) string {
	return fmt.Sprintf("%s field was not expected in JSON response", key)
}
