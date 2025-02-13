package azureblob

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/fastly/cli/pkg/cmd"
	"github.com/fastly/cli/pkg/config"
	fsterr "github.com/fastly/cli/pkg/errors"
	"github.com/fastly/cli/pkg/manifest"
	"github.com/fastly/go-fastly/v6/fastly"
)

// DescribeCommand calls the Fastly API to describe an Azure Blob Storage logging endpoint.
type DescribeCommand struct {
	cmd.Base
	manifest       manifest.Data
	Input          fastly.GetBlobStorageInput
	json           bool
	serviceName    cmd.OptionalServiceNameID
	serviceVersion cmd.OptionalServiceVersion
}

// NewDescribeCommand returns a usable command registered under the parent.
func NewDescribeCommand(parent cmd.Registerer, globals *config.Data, data manifest.Data) *DescribeCommand {
	var c DescribeCommand
	c.Globals = globals
	c.manifest = data
	c.CmdClause = parent.Command("describe", "Show detailed information about an Azure Blob Storage logging endpoint on a Fastly service version").Alias("get")
	c.RegisterFlagBool(cmd.BoolFlagOpts{
		Name:        cmd.FlagJSONName,
		Description: cmd.FlagJSONDesc,
		Dst:         &c.json,
		Short:       'j',
	})
	c.RegisterFlag(cmd.StringFlagOpts{
		Name:        cmd.FlagServiceIDName,
		Description: cmd.FlagServiceIDDesc,
		Dst:         &c.manifest.Flag.ServiceID,
		Short:       's',
	})
	c.RegisterFlag(cmd.StringFlagOpts{
		Action:      c.serviceName.Set,
		Name:        cmd.FlagServiceName,
		Description: cmd.FlagServiceDesc,
		Dst:         &c.serviceName.Value,
	})
	c.RegisterFlag(cmd.StringFlagOpts{
		Name:        cmd.FlagVersionName,
		Description: cmd.FlagVersionDesc,
		Dst:         &c.serviceVersion.Value,
		Required:    true,
	})
	c.CmdClause.Flag("name", "The name of the Azure Blob Storage logging object").Short('n').Required().StringVar(&c.Input.Name)
	return &c
}

// Exec invokes the application logic for the command.
func (c *DescribeCommand) Exec(_ io.Reader, out io.Writer) error {
	if c.Globals.Verbose() && c.json {
		return fsterr.ErrInvalidVerboseJSONCombo
	}

	serviceID, serviceVersion, err := cmd.ServiceDetails(cmd.ServiceDetailsOpts{
		AllowActiveLocked:  true,
		APIClient:          c.Globals.APIClient,
		Manifest:           c.manifest,
		Out:                out,
		ServiceNameFlag:    c.serviceName,
		ServiceVersionFlag: c.serviceVersion,
		VerboseMode:        c.Globals.Flag.Verbose,
	})
	if err != nil {
		c.Globals.ErrLog.AddWithContext(err, map[string]any{
			"Service ID":      serviceID,
			"Service Version": fsterr.ServiceVersion(serviceVersion),
		})
		return err
	}

	c.Input.ServiceID = serviceID
	c.Input.ServiceVersion = serviceVersion.Number

	azureblob, err := c.Globals.APIClient.GetBlobStorage(&c.Input)
	if err != nil {
		c.Globals.ErrLog.AddWithContext(err, map[string]any{
			"Service ID":      serviceID,
			"Service Version": serviceVersion.Number,
		})
		return err
	}

	if c.json {
		data, err := json.Marshal(azureblob)
		if err != nil {
			return err
		}
		_, err = out.Write(data)
		if err != nil {
			c.Globals.ErrLog.Add(err)
			return fmt.Errorf("error: unable to write data to stdout: %w", err)
		}
		return nil
	}

	if !c.Globals.Verbose() {
		fmt.Fprintf(out, "\nService ID: %s\n", azureblob.ServiceID)
	}
	fmt.Fprintf(out, "Version: %d\n", azureblob.ServiceVersion)
	fmt.Fprintf(out, "Name: %s\n", azureblob.Name)
	fmt.Fprintf(out, "Container: %s\n", azureblob.Container)
	fmt.Fprintf(out, "Account name: %s\n", azureblob.AccountName)
	fmt.Fprintf(out, "SAS token: %s\n", azureblob.SASToken)
	fmt.Fprintf(out, "Path: %s\n", azureblob.Path)
	fmt.Fprintf(out, "Period: %d\n", azureblob.Period)
	fmt.Fprintf(out, "GZip level: %d\n", azureblob.GzipLevel)
	fmt.Fprintf(out, "Format: %s\n", azureblob.Format)
	fmt.Fprintf(out, "Format version: %d\n", azureblob.FormatVersion)
	fmt.Fprintf(out, "Response condition: %s\n", azureblob.ResponseCondition)
	fmt.Fprintf(out, "Message type: %s\n", azureblob.MessageType)
	fmt.Fprintf(out, "Timestamp format: %s\n", azureblob.TimestampFormat)
	fmt.Fprintf(out, "Placement: %s\n", azureblob.Placement)
	fmt.Fprintf(out, "Public key: %s\n", azureblob.PublicKey)
	fmt.Fprintf(out, "File max bytes: %d\n", azureblob.FileMaxBytes)
	fmt.Fprintf(out, "Compression codec: %s\n", azureblob.CompressionCodec)

	return nil
}
