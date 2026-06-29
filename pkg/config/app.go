package config

type AppConfig struct {
	EnabledModules []string `json:"enabled_modules"`
}

func DefaultConfig() AppConfig {
	return AppConfig{
		EnabledModules: []string{"core", "tui", "monitor"},
	}
}
