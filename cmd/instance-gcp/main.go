package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/anarcher/infrakit.gcp/plugin/instance"
	"github.com/docker/infrakit/pkg/cli"
	instance_plugin "github.com/docker/infrakit/pkg/rpc/instance"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/compute/v1"
	//"strings"
)

func main() {

	var logLevel int
	var name string
	var namespaceTags []string

	cmd := &cobra.Command{
		Use:   os.Args[0],
		Short: "GCE instance plugin",
		Run: func(c *cobra.Command, args []string) {

			ctx := context.Background()

			client, err := google.DefaultClient(ctx, compute.ComputeScope)
			if err != nil {
				log.Error(err)
				os.Exit(1)
			}
			service, err := compute.New(client)
			if err != nil {
				log.Error(err)
				os.Exit(1)
			}

			instancePlugin := instance.NewInstancePlugin(service)
			cli.SetLogLevel(logLevel)
			cli.RunPlugin(name, instance_plugin.PluginServer(instancePlugin))
		},
	}

	cmd.Flags().IntVar(&logLevel, "log", cli.DefaultLogLevel, "Logging level. 0 is least verbose. Max is 5")
	cmd.Flags().StringVar(&name, "name", "instance-gcp", "Plugin name to advertise for discovery")
	cmd.Flags().StringSliceVar(
		&namespaceTags,
		"namespace-tags",
		[]string{},
		"A list of key=value resource tags to namespace all resources created")

	cmd.AddCommand(cli.VersionCommand())

	err := cmd.Execute()
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
}
