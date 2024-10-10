package ginkit

import (
	"gitlab.com/sunflower4026/interview-flash-mobile/toolkit/config"
)

func NewRuntimeConfig(cfg config.KVStore) *RuntimeConfig {
	r := RuntimeConfig{}

	r.Port = cfg.GetInt("RUNNING_PORT")
	r.ShutdownTimeoutDuration = cfg.GetDuration("TIMEOUT_DURATION")
	r.ShutdownWaitDuration = cfg.GetDuration("WAIT_DURATION")
	r.HealthCheckPath = cfg.GetString("HEALTHCHECK_PATH")
	r.InfoCheckPath = cfg.GetString("INFO_PATH")

	return &r
}
