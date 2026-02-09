package bootstrap

import (
	"news-scrabber/internal/config"
	"os"
	"strconv"
)

type InstanceID string

func NewInstanceID(cfg *config.Config) InstanceID {
	hostname, _ := os.Hostname()

	return InstanceID(cfg.AppName() + "-" + hostname + "-" + strconv.Itoa(os.Getpid()))
}

func (i InstanceID) String() string {
	return string(i)
}
