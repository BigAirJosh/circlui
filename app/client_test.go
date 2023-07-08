package client

import (
	"context"
	"fmt"
	"testing"
)

func TestFetchPipelines(t *testing.T) {
  client := New("8902bd808146a4a096b40f5289c7d449c77ec85c")
  orgSlug := "github/BigAirJosh"
  pipelines, err := client.FetchMyPipelines(context.Background(), &orgSlug, nil)
  if err != nil {
    t.Errorf("Expected to have no error but got %v", err)
  }
  fmt.Printf("pipelines: %v", pipelines)
  if len(pipelines.Items) == 0 {
    t.Errorf("Expected to have pipelines but got %v", pipelines)
  }
}
