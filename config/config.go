package config

type Configuration struct {
	// Seed
	Seeds             []string    `yaml:"seeds"`
	// Port
	Port              int         `yaml:"port"`
}
