package logtype

import (
	"fmt"
	"testing"
)

func TestLog(t *testing.T) {
	logData := []byte(`{
		"level": "error",
		"message": "Failed to connect to DB",
		"resourceId": "server-1234",
		"timestamp": "2023-09-15T08:00:00Z",
		"traceId": "abc-xyz-123",
		"spanId": "span-456",
		"commit": "5e5342f",
		"metadata": {
			"parentResourceId": "server-0987"
		}
	}`)

	logEntry := NewLog(logData)
	fmt.Printf("Unmarshaled Log Struct: %+v\n", logEntry)

	jsonOutput := logEntry.GetJsonEncoding()
	fmt.Printf("Marshaled JSON: %s\n", string(jsonOutput))

}
