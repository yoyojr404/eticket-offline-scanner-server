package config

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

// Config holds all configuration variables for the application.
type Config struct {
	Port                    string
	GinMode                 string
	BasicAuthUser           string
	BasicAuthPass           string
	DBDriver                string
	DBPath                  string
	TicketCategoriesPath    string
	TicketsPath             string
	TicketVerificationsPath string
	BatchSeedSize           int
	ScannerCertPath         string
	ScannerKeyPath          string
	CAPath                  string
}

// Load loads configuration from an optional .env file and environment variables with fallback defaults.
// If filenames are provided, it attempts to load those in order. Otherwise, it searches for .env
// in the current directory and its parent directories.
func Load(filenames ...string) *Config {
	_ = LoadEnv(filenames...)

	return &Config{
		Port:                    getEnv("PORT", "8080"),
		GinMode:                 getEnv("GIN_MODE", "debug"),
		BasicAuthUser:           getEnv("BASIC_AUTH_USER", "admin"),
		BasicAuthPass:           getEnv("BASIC_AUTH_PASS", "secret123"),
		DBDriver:                getEnv("DB_DRIVER", "sqlite3"), // memory or sqlite3
		DBPath:                  getEnv("DB_PATH", "./kartjis.db"),
		TicketCategoriesPath:    getEnv("TICKET_CATEGORIES_PATH", "./data/ticket_data_100k/ticket_categories.csv"),
		TicketsPath:             getEnv("TICKETS_PATH", "./data/ticket_data_100k/tickets.csv"),
		TicketVerificationsPath: getEnv("TICKET_VERIFICATIONS_PATH", "./data/ticket_data_100k/ticket_verifications.csv"),
		BatchSeedSize:           getIntEnv("BATCH_SEED_SIZE", 1000),
		ScannerCertPath:         getEnv("SCANNER_CERT_PATH", "./certs/scanner.crt"),
		ScannerKeyPath:          getEnv("SCANNER_KEY_PATH", "./certs/scanner.key"),
		CAPath:                  getEnv("CA_PATH", "./certs/ca.crt"),
	}
}

// LoadEnv loads environment variables from .env files into the process environment.
// Existing environment variables already set in the environment take precedence and are not overwritten.
func LoadEnv(filenames ...string) error {
	if len(filenames) > 0 {
		for _, fn := range filenames {
			if err := loadFile(fn); err != nil {
				return err
			}
		}
		return nil
	}

	// Default lookup paths when no filenames are explicitly provided
	paths := []string{".env", "../.env", "../../.env"}
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return loadFile(p)
		}
	}

	return nil
}

func loadFile(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "export ") {
			line = strings.TrimSpace(strings.TrimPrefix(line, "export "))
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		if key == "" {
			continue
		}

		val := parseEnvValue(parts[1])

		// Do not overwrite existing environment variables
		if _, exists := os.LookupEnv(key); !exists {
			_ = os.Setenv(key, val)
		}
	}

	return scanner.Err()
}

func parseEnvValue(raw string) string {
	val := strings.TrimSpace(raw)
	if len(val) >= 2 {
		if (val[0] == '"' && val[len(val)-1] == '"') || (val[0] == '\'' && val[len(val)-1] == '\'') {
			quote := val[0]
			val = val[1 : len(val)-1]
			if quote == '"' {
				val = strings.ReplaceAll(val, `\n`, "\n")
				val = strings.ReplaceAll(val, `\r`, "\r")
				val = strings.ReplaceAll(val, `\t`, "\t")
				val = strings.ReplaceAll(val, `\"`, "\"")
				val = strings.ReplaceAll(val, `\\`, "\\")
			}
			return val
		}
	}

	// For unquoted values, ignore inline comments starting with whitespace and #
	if idx := strings.Index(val, " #"); idx != -1 {
		val = strings.TrimSpace(val[:idx])
	} else if idx := strings.Index(val, "\t#"); idx != -1 {
		val = strings.TrimSpace(val[:idx])
	}

	return val
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func getIntEnv(key string, fallback int) int {
	if val := os.Getenv(key); val != "" {
		val, err := strconv.Atoi(val)
		if err != nil {
			return fallback
		}
		return val
	}
	return fallback
}
