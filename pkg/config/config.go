package config

type Config interface {
	GetString(key string) string
}
