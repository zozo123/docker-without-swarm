// +build experimental

package stack

import (
	"fmt"

	"github.com/spf13/cobra"
	"golang.org/x/net/context"

	"github.com/docker/docker/api/client"
	"github.com/docker/docker/api/client/bundlefile"
	"github.com/docker/docker/cli"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/network"
)

const (
	defaultNetworkDriver = "overlay"
)

type deployOptions struct {
	bundlefile       string
	namespace        string
	sendRegistryAuth bool
}

func newDeployCommand(dockerCli *client.DockerCli) *cobra.Command {
	var opts deployOptions

	cmd := &cobra.Command{
		Use:     "deploy [OPTIONS] STACK",
		Aliases: []string{"up"},
		Short:   "Create and update a stack from a Distributed Application Bundle (DAB)",
		Args:    cli.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.namespace = args[0]
			return runDeploy(dockerCli, opts)
		},
	}

	flags := cmd.Flags()
	addBundlefileFlag(&opts.bundlefile, flags)
	addRegistryAuthFlag(&opts.sendRegistryAuth, flags)
	return cmd
}

func runDeploy(dockerCli *client.DockerCli, opts deployOptions) error {
	bundle, err := loadBundlefile(dockerCli.Err(), opts.namespace, opts.bundlefile)
	if err != nil {
		return err
	}

	networks := getUniqueNetworkNames(bundle.Services)
	ctx := context.Background()

	if err := updateNetworks(ctx, dockerCli, networks, opts.namespace); err != nil {
		return err
	}
	return deployServices(ctx, dockerCli, bundle.Services, opts.namespace, opts.sendRegistryAuth)
}

func getUniqueNetworkNames(services map[string]bundlefile.Service) []string {
	networkSet := make(map[string]bool)
	for _, service := range services {
		for _, network := range service.Networks {
			networkSet[network] = true
		}
	}

	networks := []string{}
	for network := range networkSet {
		networks = append(networks, network)
	}
	return networks
}

func updateNetworks(
	ctx context.Context,
	dockerCli *client.DockerCli,
	networks []string,
	namespace string,
) error {
	client := dockerCli.Client()

	existingNetworks, err := getNetworks(ctx, client, namespace)
	if err != nil {
		return err
	}

	existingNetworkMap := make(map[string]types.NetworkResource)
	for _, network := range existingNetworks {
		existingNetworkMap[network.Name] = network
	}

	createOpts := types.NetworkCreate{
		Labels: getStackLabels(namespace, nil),
		Driver: defaultNetworkDriver,
		// TODO: remove when engine-api uses omitempty for IPAM
		IPAM: network.IPAM{Driver: "default"},
	}

	for _, internalName := range networks {
		name := fmt.Sprintf("%s_%s", namespace, internalName)

		if _, exists := existingNetworkMap[name]; exists {
			continue
		}
		fmt.Fprintf(dockerCli.Out(), "Creating network %s\n", name)
		if _, err := client.NetworkCreate(ctx, name, createOpts); err != nil {
			return err
		}
	}
	return nil
}

func deployServices(
	ctx context.Context,
	dockerCli *client.DockerCli,
	services map[string]bundlefile.Service,
	namespace string,
	sendAuth bool,
) error {
	apiClient := dockerCli.Client()
	out := dockerCli.Out()

	existingServices, err := getServices(ctx, apiClient, namespace)
	if err != nil {
		return err
	}

	return nil
}
