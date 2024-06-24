package appctx

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"strconv"
	"strings"
	"time"
)

type Environment struct {
}

var environment Environment

func Configs() *Environment {
	return &environment
}

// Unmarshal Read configuration information from the environment settings and parse the configuration content into the given object instance
func (Self *Environment) Unmarshal(p any) error {
	return viper.Unmarshal(p)
}

// UnmarshalKey Read the value of a key from the environment configuration and parse the content of that value into the 'rawVal' instance.
func (Self *Environment) UnmarshalKey(key string, rawVal any) error {
	return viper.UnmarshalKey(key, rawVal)
}

// GetString Read a string configuration value from the environment settings.
func (Self *Environment) GetString(key string) string {
	nk := strings.ReplaceAll(key, ".", "_")
	nk = strings.ToUpper(nk)
	if value, exists := os.LookupEnv(nk); exists {
		return value
	}
	return viper.GetString(key)
}

// GetInt64 Read an int64 configuration value from the environment settings.
func (Self *Environment) GetInt64(key string) int64 {
	nk := strings.ReplaceAll(key, ".", "_")
	nk = strings.ToUpper(nk)
	if value, exists := os.LookupEnv(nk); exists {
		if data, err := strconv.ParseInt(value, 10, 64); err == nil {
			return data
		}
	}
	return viper.GetInt64(key)
}

// GetUint64 Read an uint64 configuration value from the environment settings.
func (Self *Environment) GetUint64(key string) uint64 {
	nk := strings.ReplaceAll(key, ".", "_")
	nk = strings.ToUpper(nk)
	if value, exists := os.LookupEnv(nk); exists {
		if data, err := strconv.ParseInt(value, 10, 64); err == nil {
			return uint64(data)
		}
	}
	return viper.GetUint64(key)
}

// GetInt Read an int configuration value from the environment settings.
func (Self *Environment) GetInt(key string) int {
	nk := strings.ReplaceAll(key, ".", "_")
	nk = strings.ToUpper(nk)
	if value, exists := os.LookupEnv(nk); exists {
		if data, err := strconv.ParseInt(value, 10, 64); err == nil {
			return int(data)
		}
	}
	return viper.GetInt(key)
}

// GetBool Read a bool configuration value from the environment settings.
func (Self *Environment) GetBool(key string) bool {
	nk := strings.ReplaceAll(key, ".", "_")
	nk = strings.ToUpper(nk)
	if value, exists := os.LookupEnv(nk); exists {
		return strings.EqualFold(value, "true")
	}
	return viper.GetBool(key)
}

// GetDuration Read a Duration configuration value from the environment settings.
func (Self *Environment) GetDuration(key string) time.Duration {
	nk := strings.ReplaceAll(key, ".", "_")
	nk = strings.ToUpper(nk)
	if value, exists := os.LookupEnv(nk); exists {
		if data, err := strconv.ParseInt(value, 10, 64); err == nil {
			return time.Duration(data)
		}
	}
	return viper.GetDuration(key)
}

// GetStringSlice returns the value associated with the key as a slice of strings.
func (Self *Environment) GetStringSlice(key string) []string {
	nk := strings.ReplaceAll(key, ".", "_")
	nk = strings.ToUpper(nk)
	var i int
	var data []string
	for {
		if value, exists := os.LookupEnv(fmt.Sprintf("%s_%d", nk, i)); exists {
			data = append(data, value)
			i++
		} else {
			break
		}
	}
	if len(data) > 0 {
		return data
	}
	return viper.GetStringSlice(key)
}
