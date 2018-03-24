package bot

import (
	"fmt"
)

type Application struct {
	Name string
	Version
	ConfigPath string
	Debug      bool `yaml:"debug"`
	Bot        `yaml:"bot"`
}

type Version struct {
	Major int
	Minor int
	Patch int
}

func (self Application) Init(configPath string) Application {
	self.PrintBanner()
	self = self.LoadConfiguration(configPath)
	return self
}

func (self Application) PrintBanner() {
	fmt.Println(self.Name, ":", "v"+self.Version.ToString())
	fmt.Println("==============================")
}

func (self Version) ToString() string {
	return fmt.Sprintf("%v.%v.%v", self.Major, self.Minor, self.Patch)
}
