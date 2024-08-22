package manager

import (
	"github.com/Somesh/go-boilerplate/common/config"
)

type Module struct {
	// Add your modules here
	cfg *config.Config
}

func New() *Module {
	return &Module{}
}

func (mod *Module) Init() {
	mod.cfg = config.GetConfig()

	mod.load()
}

func (mod *Module) load() {
	// start Expedia Feed

	// Enable modules or write interfaces

}
