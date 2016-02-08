package plugins

import (
	"github.com/bruceadowns/vropsbot/common"
)

type base struct {
	config *common.Config
}

// Prepare initializes context
func (plugin *base) Prepare() (err error) {
	plugin.config, err = common.NewConfig()

	return
}
