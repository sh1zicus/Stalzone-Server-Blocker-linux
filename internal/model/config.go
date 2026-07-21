package model

type WindowConfig struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type Config struct {
	Selected []string     `json:"selected"`
	Theme    string       `json:"theme"`
	Sort     string       `json:"sort"`
	Window   WindowConfig `json:"window"`
}

func DefaultConfig() Config {
	return Config{
		Selected: []string{},
		Theme:    "dark",
		Sort:     "ping",
		Window: WindowConfig{
			Width:  120,
			Height: 35,
		},
	}
}
