package client

import (
	"context"
	"log"

	circleci "github.com/grezar/go-circleci"
)

type pipelines struct {
	client circleci.Client
}

func New(token string) *pipelines {
	config := circleci.DefaultConfig()
	config.Token = token

	client, err := circleci.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}

	// contexts, err := client.Contexts.List(context.Background(), circleci.ContextListOptions{
	// 	OwnerSlug: circleci.String("fb972d34-235e-432b-bb61-36fedc6445d2"),
	// })
 //  log.Printf("contexts: %v", contexts)

	// if err != nil {
	// 	log.Fatal(err)
	// }
  return &pipelines{
    client: *client,
  }
}

func (p *pipelines) FetchMyPipelines(ctx context.Context, orgSlug *string, pageToken *string) (*circleci.PipelineList, error) {
  mine := false
  options := circleci.PipelineListOptions{
  	OrgSlug:   orgSlug,
  	Mine:      &mine,
  	PageToken: pageToken,
  }
  pipelines, err := p.client.Pipelines.List(ctx, options)

  if err != nil {
    return nil, err
  }

  return pipelines, nil
  
}
