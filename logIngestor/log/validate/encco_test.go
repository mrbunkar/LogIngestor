package validate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"logIngestor/logIngestor/log/logtype"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEncDec(t *testing.T) {
	decoder := NewDecoder()
	log := &logtype.Log{}
	data := []byte("To test")

	reader := io.NopCloser(bytes.NewReader(data))
	err := decoder.Decode(reader, log)

	assert.NotEqual(t, err, nil)

	levels := []string{"INFO", "ERROR", "WARN", "DEBUG"}
	messages := []string{"This is a message 1", "This is a message 2", "This is a message 3", "This is a debug message"}
	logReq := map[string]interface{}{
		"Level":      levels[rand.Intn(4)],
		"Message":    messages[rand.Intn(4)],
		"ResourceId": generateRandomStringNumber(4),
		"Timestamp":  time.Now(),
		"TraceId":    "trace-" + generateRandomStringNumber(4),
		"SpanId":     "span-" + generateRandomStringNumber(4),
		"Commit":     generateRandomStringNumber(6),
		"Metadata": map[string]string{
			"ParentResourceId": generateRandomStringNumber(4),
		},
	}
	bt, _ := json.Marshal(logReq)
	lreader := io.NopCloser(bytes.NewReader(bt))
	assert.Equal(t, decoder.Decode(lreader, log), nil)
	fmt.Println(log)
}

func generateRandomStringNumber(n int) string {
	const letters = "1234567890"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
