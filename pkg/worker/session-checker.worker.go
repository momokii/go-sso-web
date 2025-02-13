package worker

import (
	"log"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/momokii/go-sso-web/internal/database"
	"github.com/momokii/go-sso-web/internal/repository/session"
)

type SessionChecker struct {
	sessionRepo session.SessionRepo
}

func NewSessionChecker(sessionRepo session.SessionRepo) *SessionChecker {
	return &SessionChecker{
		sessionRepo: sessionRepo,
	}
}

func (s *SessionChecker) CheckSession() error {
	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}
	defer func() {
		database.CommitOrRollback(tx, nil, err)
		log.Println("Worker Session Checker Executed")
	}()

	time_now := time.Now().Format(time.RFC3339)
	if err := s.sessionRepo.DeleteExpiredSession(tx, time_now); err != nil {
		return err
	}

	return nil
}

func (s *SessionChecker) StartChecker(durationJob time.Duration) {
	// Setup GoCron scheduler
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		log.Println("Error creating scheduler: ", err)
		panic(err)
	}

	// add a job to the scheduler
	if _, err := scheduler.NewJob(
		gocron.DurationJob(
			durationJob,
		),
		gocron.NewTask(s.CheckSession),
	); err != nil {
		log.Println("Error creating job: ", err)
		panic(err)
	}

	scheduler.Start() // start the scheduler
	log.Println("Worker Session Checker Started")
}
