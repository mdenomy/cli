package service

import (
	"io"

	"github.com/fastly/cli/pkg/cmd"
	"github.com/fastly/cli/pkg/config"
	"github.com/fastly/cli/pkg/manifest"
	"github.com/fastly/cli/pkg/text"
	"github.com/fastly/go-fastly/v6/fastly"
)

// SearchCommand calls the Fastly API to describe a service.
type SearchCommand struct {
	cmd.Base
	manifest manifest.Data
	Input    fastly.SearchServiceInput
}

// NewSearchCommand returns a usable command registered under the parent.
func NewSearchCommand(parent cmd.Registerer, globals *config.Data, data manifest.Data) *SearchCommand {
	var c SearchCommand
	c.Globals = globals
	c.manifest = data
	c.CmdClause = parent.Command("search", "Search for a Fastly service by name")
	c.CmdClause.Flag("name", "Service name").Short('n').Required().StringVar(&c.Input.Name)
	return &c
}

// Exec invokes the application logic for the command.
func (c *SearchCommand) Exec(_ io.Reader, out io.Writer) error {
	service, err := c.Globals.APIClient.SearchService(&c.Input)
	if err != nil {
		c.Globals.ErrLog.AddWithContext(err, map[string]any{
			"Service Name": c.Input.Name,
		})
		return err
	}

	text.PrintService(out, "", service)
	return nil
}
