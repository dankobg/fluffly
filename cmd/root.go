package cmd

import (
	"github.com/alecthomas/kong"
	"github.com/dankobg/fluffly/cmd/fluffly"
	"github.com/dankobg/fluffly/cmd/identities"
)

var CLI struct {
	Serve      fluffly.ServeCommand `cmd:"" help:"Run Fluffly Server"`
	Identities identities.RootCmd   `cmd:"" help:"Manage identities"`
}

func Run() {
	c := kong.Parse(
		&CLI,
		kong.Name("fluffly"),
		kong.Description("Fluffly pet finder app, find your best pal"),
	)

	err := c.Run()
	c.FatalIfErrorf(err)
}
