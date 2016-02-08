package plugins

func init() {
	// bot.Register(&reportPush{})
}

type reportPush struct {
	base
}

// Name returns report push's name
func (plugin *reportPush) Name() string {
	return "reportPush"
}

// Run starts the message pump for build push
func (plugin *reportPush) Run() (err error) {
	return
}
