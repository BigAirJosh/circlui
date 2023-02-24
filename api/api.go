package api

import (
	_"encoding/json"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var cliConfig CliConfig = LoadCliConfig()

type CliConfig struct {
	Host         string `yaml:"host"`
	Token        string `yaml:"token"`
	RestEndpoint string `yaml:"rest_endpoint"`
  Scheme       string `yaml:"scheme"`
}

func LoadCliConfig() CliConfig {

	cliConfig := CliConfig{}

	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	dat, err := os.ReadFile(userHomeDir + "/.circleci/cli.yml")

	if err != nil {
		log.Fatalf("error loading ~/.circleci/cli.yml: %v", err)
	}

	err = yaml.Unmarshal([]byte(dat), &cliConfig)
	if err != nil {
		log.Fatalf("error parsing ~/.circleci/cli.yml: %v", err)
	}
	return cliConfig
}

type Workflow struct {
	PipelineId     string `"pipeline_id"`
	PipelineNumber string `"pipeline_number"`
	Id             string `"id"`
	Name           string `"name"`
	ProjectSlug    string `"project_slug"`
	Tag            string `"tag"`
	Status         string `"status"`
	StartedBy      string `"started_by"`
	CancelledBy    string `"canceled_by"`
	ErroredBy      string `"errored_by"`
	CreatedAt      string `"created_at"`
	StoppedAt      string `"stopped_at"`
}

type PipelineWorkflowsResponse struct {
	Items         []Workflow `"items"`
	NextPageToken string     `"next_page_token`
}

func LoadPipelinesWorkflows(pipelineId string) {

  url := cliConfig.Host + "/" + cliConfig.RestEndpoint + "/pipeline/" + pipelineId + "/workflow"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Circle-Token", cliConfig.Token)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	log.Println(res)
	log.Println(string(body))

}

func LoadPipelines() {

  url := cliConfig.Host + "/" + cliConfig.RestEndpoint + "/pipeline?org-slug=circleci%2FLs9DEFFLtnRvAyCqbc3AkG"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Circle-Token", cliConfig.Token)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	log.Println(res)
	log.Println(string(body))

}
