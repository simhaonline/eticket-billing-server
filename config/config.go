package config

import (
    "code.google.com/p/gcfg"
    "github.com/golang/glog"
    "path/filepath"
    "os"
)

type Config struct {
    environment string
    DatabaseName string `gcfg:"database-name"`
    DatabaseUser string `gcfg:"database-user"`
    DatabasePassword string `gcfg:"database-password"`
    RequestLogDir string `gcfg:"request-log-dir"`
}


var cfg = struct {
    Environment map[string]*Config
}{}

var currentEnvironment string

func ParseConfig(env string, configFile string) {
    currentEnvironment = env
    err := gcfg.ReadFileInto(&cfg, configFile)
    if err != nil {
        glog.Fatal(err)
        panic (err)
    }

    // TODO hold situation when there are no section for current env
    if cfg.Environment[currentEnvironment].RequestLogDir == "" {
        dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
        if err != nil {
            glog.Fatal(err)
            os.Exit(1)
        }

        cfg.Environment[currentEnvironment].RequestLogDir = dir
        glog.Infof("Use directory `%v' as root for storing log files", dir)
    }

    cfg.Environment[currentEnvironment].environment = currentEnvironment

    glog.Infof("Configured to run in %v environment", currentEnvironment)

    if err != nil {
        glog.Fatal(err)
        panic(err)
    }
}

func GetConfig() *Config {
    return cfg.Environment[currentEnvironment]
}
