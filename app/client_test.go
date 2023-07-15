package app_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/bigairjosh/circlui/app"
)

func TestFetchPipelines(t *testing.T) {
  c := app.New("8902bd808146a4a096b40f5289c7d449c77ec85c")
  orgSlug := "github/BigAirJosh"
  pipelines, err := c.FetchMyPipelines(context.Background(), &orgSlug, nil)
  if err != nil {
    t.Errorf("Expected to have no error but got %v", err)
  }
  fmt.Printf("State: %v\n", pipelines.Items[0].State)
  fmt.Printf("Who: %v\n", pipelines.Items[0].Trigger.Actor.Login)
  fmt.Printf("Branch: %v\n", pipelines.Items[0].Vcs.Branch)
  fmt.Printf("Errors: %v\n", len(pipelines.Items[0].Errors))
  fmt.Printf("Type %v\n", pipelines.Items[0].Trigger.Type)

  if len(pipelines.Items) == 0 {
    t.Errorf("Expected to have pipelines but got %v", pipelines)
  }

  workflows, err := c.FetchWorkflows(context.Background(), pipelines.Items[0].ID)
  if err != nil {
    t.Errorf("Expected to have no error but got %v", err)
  }

  fmt.Printf("Workflows: %v\n", workflows.Items[0])
}
