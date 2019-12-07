package config

import (
	"os"
	"strconv"
	"strings"
)

// CloudStorageConfig holds cloud storage environment variables
type CloudStorageConfig struct {
	Bucket  string
	URI     string
	Project string
}

// CloudStorageSettingsConfig holds cloud storage environment variables
type CloudStorageSettingsConfig struct {
	ProfileBucket string
	URI           string
}

// Config holds config constructs
type Config struct {
	CloudStorage  CloudStorageConfig
	CloudSettings CloudStorageSettingsConfig
}

// New returns a new Config struct
func New() *Config {
	return &Config{
		CloudStorage: CloudStorageConfig{
			Bucket:  getEnv("BUCKET_NAME", ""),
			URI:     (getEnv("GOOGLE_CLOUD_STORAGE_URI", "")) + (getEnv("BUCKET_NAME", "")) + "/",
			Project: getEnv("GOOGLE_CLOUD_PROJECT", ""),
		},
		CloudSettings: CloudStorageSettingsConfig{
			ProfileBucket: getEnv("PROFILE_BUCKET_NAME", ""),
			URI:           (getEnv("GOOGLE_CLOUD_STORAGE_URI", "")) + (getEnv("PROFILE_BUCKET_NAME", "")) + "/",
		},
	}
}

// Simple helper function to read an environment or return a default value
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

// Simple helper function to read an environment variable into integer or return a default value
func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultVal
}

// Helper to read an environment variable into a bool or return default value
func getEnvAsBool(name string, defaultVal bool) bool {
	valStr := getEnv(name, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}

	return defaultVal
}

// Helper to read an environment variable into a string slice or return default value
func getEnvAsSlice(name string, defaultVal []string, sep string) []string {
	valStr := getEnv(name, "")

	if valStr == "" {
		return defaultVal
	}

	val := strings.Split(valStr, sep)

	return val
}
