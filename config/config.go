package config

import (
	"code.google.com/p/gcfg"
	"github.com/golang/glog"
	"os"
	"path/filepath"
)

type Config struct {
	DatabaseName     string `gcfg:"database-name"`
	DatabaseUser     string `gcfg:"database-user"`
	DatabasePassword string `gcfg:"database-password"`
	RequestLogDir    string `gcfg:"request-log-dir"`
}

type configList struct {
	Environment map[string]*Config
}

func NewConfig(env string, configFile string) *Config {
	cfg := configList{}

	err := gcfg.ReadFileInto(&cfg, configFile)
	if err != nil {
		glog.Fatal(err)
		panic(err)
	}

	// TODO hold situation when there are no section for current env
	if cfg.Environment[env].RequestLogDir == "" {
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			glog.Fatal(err)
			os.Exit(1)
		}

		cfg.Environment[env].RequestLogDir = dir
		glog.Infof("Use directory `%v' as root for storing log files", dir)
	}

	glog.Infof("Configured to run in %v environment", env)

	if err != nil {
		glog.Fatal(err)
		panic(err)
	}

	return cfg.Environment[env]
}
