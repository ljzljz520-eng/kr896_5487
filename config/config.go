package config

import "os"

type Config struct {
	Address, DBPath, StaticDir string
	MaxBodyBytes               int64
}

func Defaults() Config {
	return Config{Address: ":8080", DBPath: "ruralfolk.db", StaticDir: "web", MaxBodyBytes: 2 << 20}
}
func FromEnv() Config {
	c := Defaults()
	if v := os.Getenv("RURALFOLK_ADDR"); v != "" {
		c.Address = v
	}
	if v := os.Getenv("RURALFOLK_DB"); v != "" {
		c.DBPath = v
	}
	if v := os.Getenv("RURALFOLK_WEB"); v != "" {
		c.StaticDir = v
	}
	return c
}
func (c Config) Validate() error {
	if c.Address == "" || c.DBPath == "" || c.StaticDir == "" || c.MaxBodyBytes <= 0 {
		return os.ErrInvalid
	}
	return nil
}
func (c Config) PublicURL() string { return "http://localhost" + c.Address }
