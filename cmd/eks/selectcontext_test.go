package eks

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"testing"
)

type mockEKSClient struct {
	ListClustersFunc func(ctx context.Context, params *eks.ListClustersInput, optFns ...func(*eks.Options)) (*eks.ListClustersOutput, error)
}

func (m *mockEKSClient) ListClusters(ctx context.Context, params *eks.ListClustersInput, optFns ...func(*eks.Options)) (*eks.ListClustersOutput, error) {
	return m.ListClustersFunc(ctx, params, optFns...)
}

type testcase struct {
	clusterResponse []string
	err             error
}

func TestRetrieveClusters(t *testing.T) {
	tests := []testcase{
		{
			clusterResponse: []string{"cluster-1"},
			err:             nil,
		},
		{
			clusterResponse: make([]string, 0),
			err:             errNoClustersFound,
		},
	}

	for _, test := range tests {
		client := mockEksClient(test.clusterResponse)
		_, err := retrieveClusters(client)

		if !errors.Is(err, test.err) {
			t.Errorf("Expected error to be %s, got: %s", test.err, err)
		}
	}
}

func mockEksClient(clusters []string) *mockEKSClient {
	return &mockEKSClient{
		ListClustersFunc: func(ctx context.Context, params *eks.ListClustersInput, optFns ...func(*eks.Options)) (*eks.ListClustersOutput, error) {
			return &eks.ListClustersOutput{
				Clusters: clusters,
			}, nil
		},
	}
}
