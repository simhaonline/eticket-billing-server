package server

type Config struct {
    Environment string
    DataBaseName string
    RequestLogDir string
}

func NewConfig() *Config {
    return &Config{}
}
