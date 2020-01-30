package main

import (
	"github.com/spf13/pflag"
	"time"
)

type Config struct {
	Host string
	Port int

	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func (c *Config) Flags() *pflag.FlagSet {
	f := pflag.NewFlagSet("Swagger API Config", pflag.PanicOnError)

	f.StringVar(&c.Host, "host", "0.0.0.0", "ip host")
	f.IntVar(&c.Port, "port", 80, "port")
	f.DurationVar(&c.ReadTimeout, "readtimeout", time.Duration(0), "api read timeout (default 0s)")
	f.DurationVar(&c.WriteTimeout, "writetimeout", time.Duration(0), "api write timeout (default 0s)")

	return f
}
