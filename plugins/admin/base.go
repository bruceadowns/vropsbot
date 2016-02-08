package plugins

// logs is the base struct for admin plugins
type base struct {
}

// Prepare initializes logs context
func (p *base) Prepare() error {
	return nil
}
