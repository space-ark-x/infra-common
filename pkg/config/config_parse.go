package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func LoadConfig() *Type {
	env := os.Getenv("Env")
	if env == "" {
		panic("env variable not set")
	}
	reader, err := os.ReadFile(fmt.Sprintf("./config/%s.yaml", env))
	if err != nil {
		panic(err)
	}
	var config *Type
	err = yaml.Unmarshal(reader, &config)
	if err != nil {
		panic(err)
	}
	return config
}
