package httpserver

import (
	"bytes"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func generateRandomStringNumber(n int) string {
	const letters = "1234567890"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func TestHttpAddEndpoint(t *testing.T) {

	levels := []string{"INFO", "ERROR", "WARN", "DEBUG"}
	messages := []string{"This is a message 1", "This is a message 2", "This is a message 3", "This is a debug message"}
	for j := 0; j < 10; j++ {
		logReq := map[string]interface{}{
			"level":      levels[rand.Intn(4)],
			"message":    messages[rand.Intn(4)],
			"resourceId": generateRandomStringNumber(4),
			"timestamp":  time.Now(),
			"traceId":    "trace-" + generateRandomStringNumber(4),
			"spanId":     "span-" + generateRandomStringNumber(4),
			"commit":     generateRandomStringNumber(6),
			"metadata": map[string]string{
				"parentResourceId": generateRandomStringNumber(4),
			},
		}
		jsonLogReq, err := json.Marshal(logReq)
		if err != nil {
			panic("error marshalling log request")
		}

		res, err := http.Post("http://localhost:3000/Add", "application/json", bytes.NewBuffer(jsonLogReq))
		if err != nil {
			panic("error posting log request")
		}
		log.Println("Status code: ", res.StatusCode)
		assert.Equal(t, res.StatusCode, 200)
	}
}
