package worker

import (
	"log"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/momokii/go-sso-web/internal/database"
	"github.com/momokii/go-sso-web/pkg/repository/user"
)

type UserCreditChecker struct {
	userRepo user.UserRepo
}

func NewUserCreditChecker(userRepo user.UserRepo) *UserCreditChecker {
	return &UserCreditChecker{
		userRepo: userRepo,
	}
}

func (uc *UserCreditChecker) CheckUserCredit() error {
	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}
	defer func() {
		database.CommitOrRollback(tx, nil, err)
		log.Println("Worker User Credit Checker Executed")
	}()

	if err := uc.userRepo.ResetUserDailyToken(tx); err != nil {
		return err
	}

	return nil
}

func (uc *UserCreditChecker) StartChecker(durationJob time.Duration) {
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
		gocron.NewTask(uc.CheckUserCredit),
	); err != nil {
		log.Println("Error creating job: ", err)
		panic(err)
	}

	scheduler.Start() // start the scheduler
	log.Println("Worker User Credit Checker Started")
}
