// Package env provides the utility functions for env values preparation
package env

import (
	"os"

	"github.com/atomgunlk/golang-common/pkg/logger"
)

// RequiredEnv gets value from environment variable and panic if it is not set
func RequiredEnv(k string) string {
	env := os.Getenv(k)
	if len(env) == 0 {
		logger.Panicf("[RequiredEnv]: missing environment variable `%s`", k)
	}
	logger.Debugf("[RequiredEnv]: read environment variable %s: %s", k, env)
	return env
}
