package eks

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/spf13/cobra"
	"os/exec"

	"github.com/manifoldco/promptui"
)

var errNoClustersFound = errors.New("no clusters found")

type Client interface {
	ListClusters(ctx context.Context, params *eks.ListClustersInput, optFns ...func(*eks.Options)) (*eks.ListClustersOutput, error)
}

var selectContextCommand = &cobra.Command{
	Use:   "select-context",
	Short: "Interactively select EKS context",
	RunE:  executeSelectContext(),
}

func executeSelectContext() func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		region, err := cmd.Flags().GetString("region")
		if err != nil {
			return fmt.Errorf("failed to get region flag: %w", err)
		}

		cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
		if err != nil {
			return fmt.Errorf("failed to load AWS config: %w", err)
		}

		client := eks.NewFromConfig(cfg)

		clusters, err := retrieveClusters(client)
		if err != nil {
			return fmt.Errorf("failed to retrieve clusters: %w", err)
		}

		selectedCluster, err := selectCluster(clusters)
		if err != nil {
			return fmt.Errorf("failed to select cluster: %w", err)
		}

		if err := updateKubeconfig(selectedCluster, region); err != nil {
			return fmt.Errorf("failed to update kubeconfig: %w", err)
		}
		return nil
	}
}

func updateKubeconfig(result string, region string) error {
	ext := exec.Command("aws", "eks", "update-kubeconfig", "--name", result, "--region", region)
	err := ext.Run()
	if err != nil {
		return err
	}
	fmt.Printf("Successfully connected to cluster: %s\n", result)
	return nil
}

func retrieveClusters(client Client) ([]string, error) {
	output, err := client.ListClusters(context.TODO(), nil)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve clusters: %w", err)
	}

	if len(output.Clusters) == 0 {
		return nil, errNoClustersFound
	}
	return output.Clusters, nil
}

func selectCluster(clusters []string) (string, error) {
	prompt := promptui.Select{
		Label: "Select Cluster",
		Items: clusters,
	}
	_, result, err := prompt.Run()
	if err != nil {
		return "", err
	}
	return result, nil
}
