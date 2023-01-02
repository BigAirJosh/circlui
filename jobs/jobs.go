package jobs

import (
	"github.com/brianvoe/gofakeit/v6"
  "github.com/hako/durafmt"
  "time"
)

type Job struct {
  Id string
  Name string
  Branch string
  RemianingSeconds int
}

func RandomJobs(count int) []Job {

  jobs := make([]Job, count)
  for i := 0; i < count; i++ {
    jobs[i] = Job {
      Id: gofakeit.AchAccount(),
      Name: gofakeit.AppName(),
      Branch: gofakeit.DomainSuffix(),
      RemianingSeconds: gofakeit.Number(0, 9999),
    }
  } 

  return jobs
}

func (j Job) FormatRemainingTime() string {
  timeduration := time.Second * time.Duration(j.RemianingSeconds)
  return durafmt.Parse(timeduration).String()
}

func (j Job) CountDown() {
  j.RemianingSeconds--
}

