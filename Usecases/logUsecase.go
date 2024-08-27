// log_usecase.go
package Usecases

import (
	"Loan_manager/Domain"
	"Loan_manager/Repository"
)

type LogUsecase interface {
	GetAllLogs() ([]Domain.Log, error)
}

type logUsecase struct {
	logRepo Repository.LogRepository
}

func NewLogUsecase(logRepo Repository.LogRepository) LogUsecase {
	return &logUsecase{logRepo: logRepo}
}

func (lu *logUsecase) GetAllLogs() ([]Domain.Log, error) {
	return lu.logRepo.GetAllLogs()
}
