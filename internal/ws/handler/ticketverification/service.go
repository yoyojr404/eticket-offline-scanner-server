package ticketverification

import (
	"context"

	"github.com/yoyojr404/internal/common/datastore"
)

type Service interface {
	VerifyTicket(ctx context.Context, verificationId int64) error
}

type verifyJob struct {
	ctx            context.Context
	verificationID int64
	errchan        chan error
}

type service struct {
	verifyJobChan chan verifyJob
	transactor    datastore.Transactor
	repo          Repository
}

func NewService(transactor datastore.Transactor, repo Repository) Service {
	service := service{
		verifyJobChan: make(chan verifyJob, 256),
		transactor:    transactor,
		repo:          repo,
	}

	go service.process()

	return &service
}

func (s *service) VerifyTicket(ctx context.Context, verificationId int64) error {
	job := verifyJob{
		verificationID: verificationId,
		errchan:        make(chan error, 1),
		ctx:            ctx,
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case s.verifyJobChan <- job:
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-job.errchan:
		return err
	}
}

func (s *service) process() {
	for job := range s.verifyJobChan {
		err := s.transactor.WithinTx(job.ctx, func(tx datastore.Executor) error {
			repo := s.repo.WithTx(tx)

			tv, err := repo.FindTicketVerificationByID(job.ctx, job.verificationID)
			if err != nil {
				return err
			}

			if tv == nil {
				return ErrTicketVerificationNotFound
			}

			if tv.Verified {
				return ErrTicketVerificationHasBeenVerified
			}

			ticketExist, err := repo.UpdateVerifiedTicketVerification(job.ctx, job.verificationID)

			if err != nil {
				return err
			}

			if !ticketExist {
				return ErrTicketVerificationNotFound
			}

			return nil
		})

		job.errchan <- err
	}

}
