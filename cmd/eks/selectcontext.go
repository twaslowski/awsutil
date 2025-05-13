package eks

import (
	util "awsutil/pkg"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/spf13/cobra"
	"log"
	"os/exec"

	"github.com/manifoldco/promptui"
)

var selectContextCommand = &cobra.Command{
	Use:   "select-context",
	Short: "Interactively select EKS context",
	Run:   executeSelectContext(),
}

func executeSelectContext() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		region := util.Require(cmd.Flags().GetString("region"))
		config := util.LoadConfiguration(region)

		client := eks.NewFromConfig(config)
		clusters := retrieveClusters(*client)
		selectedCluster := selectCluster(clusters)
		updateKubeconfig(selectedCluster, region)
	}
}

func updateKubeconfig(result string, region string) {
	ext := exec.Command("aws", "eks", "update-kubeconfig", "--name", result, "--region", region)
	err := ext.Run()
	if err != nil {
		log.Fatalf("Failed to update cluster config: %s", err)
	}
	fmt.Printf("Successfully connected to cluster: %s\n", result)
}

func retrieveClusters(client eks.Client) []string {
	output, err := client.ListClusters(context.TODO(), nil)
	if err != nil {
		log.Fatalf("Could not retrieve clusters: %s", err)
	}
	return output.Clusters
}

func selectCluster(clusters []string) string {
	prompt := promptui.Select{
		Label: "Select Cluster",
		Items: clusters,
	}
	_, result, err := prompt.Run()
	if err != nil {
		log.Fatalf("Failed to select cluster: %s", err)
	}
	return result
}
