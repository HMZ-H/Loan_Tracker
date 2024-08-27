package Usecases

import (
	"Loan_manager/Domain"
	"Loan_manager/Repository"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LoanUsecase interface {
	ApplyLoan(loan Domain.Loan) (string, error)
	ViewLoanStatus(loanID primitive.ObjectID) (*Domain.Loan, error)
	ViewAllLoans(status string, order string) ([]Domain.Loan, error)
	ApproveRejectLoan(loanID primitive.ObjectID, status string) (*Domain.Loan, error)
	DeleteLoan(loanID primitive.ObjectID) error
}

type loanUsecase struct {
	loanRepo Repository.LoanRepository
}

func NewLoanUsecase(loanRepo Repository.LoanRepository) LoanUsecase {
	return &loanUsecase{
		loanRepo: loanRepo,
	}
}

func (lu *loanUsecase) ApplyLoan(loan Domain.Loan) (string, error) {
	loan.ID = primitive.NewObjectID()
	loan.CreatedAt = time.Now()
	loan.Status = "pending"

	if err := lu.loanRepo.CreateLoan(loan); err != nil {
		return "", err
	}

	return loan.Status, nil
}

func (lu *loanUsecase) ViewLoanStatus(loanID primitive.ObjectID) (*Domain.Loan, error) {
	return lu.loanRepo.GetLoanByID(loanID)
}

func (lu *loanUsecase) ViewAllLoans(status string, order string) ([]Domain.Loan, error) {
	return lu.loanRepo.GetAllLoans(status, order)
}

func (lu *loanUsecase) ApproveRejectLoan(loanID primitive.ObjectID, status string) (*Domain.Loan, error) {
	loan, err := lu.loanRepo.GetLoanByID(loanID)
	if err != nil {
		return nil, err
	}

	loan.Status = status
	if status == "approved" {
		now := time.Now()
		loan.ApprovedAt = &now
	}

	if err := lu.loanRepo.UpdateLoan(loan); err != nil {
		return nil, err
	}

	return loan, nil
}

func (lu *loanUsecase) DeleteLoan(loanID primitive.ObjectID) error {
	return lu.loanRepo.DeleteLoan(loanID)
}
