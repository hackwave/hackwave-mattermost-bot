package bot

type Application struct {
	Name    string
	Version string
	Bot     `yaml:"bot"`
}
