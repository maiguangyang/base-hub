package main

import (
	"net/url"
	"os"

	"base-engine/gen"
)

func openEngineDBFromEnv() (*gen.DB, error) {
	rawURL := os.Getenv("DATABASE_URL")
	if rawURL == "" {
		return gen.OpenDBFromEnvVars("")
	}
	return openEngineDB(rawURL)
}

// openEngineDB enables SQLite foreign keys on every pooled connection through
// the driver DSN. The generated database package intentionally remains untouched.
func openEngineDB(rawURL string) (*gen.DB, error) {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return gen.OpenDBWithString(rawURL)
	}
	if parsed.Scheme == "sqlite3" {
		query := parsed.Query()
		query.Set("_foreign_keys", "on")
		parsed.RawQuery = query.Encode()
		rawURL = parsed.String()
	}
	return gen.OpenDBWithString(rawURL)
}
