package logging

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitializeLogger(t *testing.T) {
	stderr := os.Stderr
	reader, writer, _ := os.Pipe()
	os.Stderr = writer

	defer func() {
		os.Stderr = stderr
		writer.Close()
	}()

	InitializeLogger()

	slog.Debug("This is not logged")
	slog.Info("Test message", "key", "value")
	slog.Warn("We", "are", "gonna", "die", "!")
	slog.Error("Omg")

	writer.Close()
	var buffer bytes.Buffer
	buffer.ReadFrom(reader)
	output := buffer.String()

	lines := strings.Split(strings.TrimSpace(output), "\n")
	assert.Equal(t, len(lines), 3, "Info, Warn, Error are the only lines that should be logged")

	expectedLogs := []string{
		`{"level":"INFO","msg":"Test message","key":"value"}`,
		`{"level":"WARN","msg":"We","are":"gonna","die":"!"}`,
		`{"level":"ERROR","msg":"Omg"}`,
	}

	for index, line := range lines {
		var logEntry map[string]interface{}
		err := json.Unmarshal([]byte(line), &logEntry)
		assert.NoError(t, err, "Log output should be valid JSON")

		delete(logEntry, "time")
		outputWithoutTime, err := json.Marshal(logEntry)
		assert.NoError(t, err)

		expected := expectedLogs[index]
		assert.JSONEq(t, expected, string(outputWithoutTime))
	}

}
