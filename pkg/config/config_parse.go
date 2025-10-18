package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func LoadConfig() *Type {
	env := os.Getenv("Env")
	configInject := os.Getenv("ConfigInject")
	injectMode := configInject == "true"
	config := &Type{
		Record: Record{},
	}
	if !injectMode {
		if env == "" {
			panic("env variable not set")
		}
		reader, err := os.ReadFile(fmt.Sprintf("./config/%s.yaml", env))
		if err != nil {
			panic(err)
		}
		err = yaml.Unmarshal(reader, &config)
		if err != nil {
			panic(err)
		}
		err = yaml.Unmarshal(reader, &config.Record)
		if err != nil {
			panic(err)
		}
		return config
	}
	return config
}
