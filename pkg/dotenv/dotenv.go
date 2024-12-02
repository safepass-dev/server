package dotenv

import (
	"bufio"
	"os"
	"strings"
)

// LoadEnv loads environment variables from a specified .env file.
// If the file does not exist or contains invalid entries, an error is returned.
func LoadEnv(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Ignore comments and empty lines
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Split the line into key and value
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return &InvalidEnvLineError{Line: line}
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Remove surrounding quotes if present
		value = strings.Trim(value, `"'`)

		// Set the environment variable
		if err := os.Setenv(key, value); err != nil {
			return err
		}
	}

	// Handle scanner errors
	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

// InvalidEnvLineError represents an error for invalid .env file lines.
type InvalidEnvLineError struct {
	Line string
}

func (e *InvalidEnvLineError) Error() string {
	return "invalid .env line: " + e.Line
}
