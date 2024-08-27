package controller

import (
	"Loan_manager/Domain"
	"Loan_manager/Usecases"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LoanController struct {
	loanUsecase Usecases.LoanUsecase
}

func NewLoanController(loanUsecase Usecases.LoanUsecase) *LoanController {
	return &LoanController{
		loanUsecase: loanUsecase,
	}
}

// Apply for Loan
func (lc *LoanController) ApplyLoan(c *gin.Context) {
	var loanRequest Domain.Loan
	if err := c.ShouldBindJSON(&loanRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	loanRequest.UserID = primitive.NewObjectID() // Replace with actual user ID after authentication
	loanStatus, err := lc.loanUsecase.ApplyLoan(loanRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": loanStatus})
}

// View Loan Status
func (lc *LoanController) ViewLoanStatus(c *gin.Context) {
	loanID := c.Param("id")
	loanObjectID, err := primitive.ObjectIDFromHex(loanID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid loan ID"})
		return
	}

	loan, err := lc.loanUsecase.ViewLoanStatus(loanObjectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, loan)
}

// View All Loans (Admin)
func (lc *LoanController) ViewAllLoans(c *gin.Context) {
	status := c.DefaultQuery("status", "all")
	order := c.DefaultQuery("order", "asc")

	loans, err := lc.loanUsecase.ViewAllLoans(status, order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, loans)
}

// Approve/Reject Loan (Admin)
func (lc *LoanController) ApproveRejectLoan(c *gin.Context) {
	loanID := c.Param("id")
	loanObjectID, err := primitive.ObjectIDFromHex(loanID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid loan ID"})
		return
	}

	var statusUpdate struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&statusUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedLoan, err := lc.loanUsecase.ApproveRejectLoan(loanObjectID, statusUpdate.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedLoan)
}

// Delete Loan (Admin)
func (lc *LoanController) DeleteLoan(c *gin.Context) {
	loanID := c.Param("id")
	loanObjectID, err := primitive.ObjectIDFromHex(loanID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid loan ID"})
		return
	}

	if err := lc.loanUsecase.DeleteLoan(loanObjectID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Loan deleted successfully"})
}
