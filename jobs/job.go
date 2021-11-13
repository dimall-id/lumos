package jobs

import (
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

type Job interface {
	 GetSchedule () string
	 Execute ()
}

var jobs = make([]Job, 0)
func AddJob (job Job) {
	jobs = append(jobs, job)
}

func StartJob (log *logrus.Entry) error {
	scheduler := cron.New()
	defer scheduler.Stop()

	for _, job := range jobs {
		id, err := scheduler.AddFunc(job.GetSchedule(), job.Execute)
		if err != nil {log.Fatalln(err)}
		log.WithField("EntryID", id).Warn("Entry ID")
	}

	scheduler.Start()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	return nil
}