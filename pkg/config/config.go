// Package config for load any service config
package config

import "github.com/ilyakaznacheev/cleanenv"

func Init(file string, config any) error {
	err := cleanenv.ReadConfig(file, config)
	if err != nil {
		return err
	}

	err = cleanenv.ReadEnv(config)
	if err != nil {
		return err
	}

	return nil
}

func MustInit(file string, config any) {
	err := Init(file, config)
	if err != nil {
		panic(err)
	}
}
