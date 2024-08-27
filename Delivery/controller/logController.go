// log_controller.go
package controller

import (
	"Loan_manager/Usecases"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LogController struct {
	logUsecase Usecases.LogUsecase
}

func NewLogController(logUsecase Usecases.LogUsecase) *LogController {
	return &LogController{logUsecase: logUsecase}
}

func (lc *LogController) ViewSystemLogs(c *gin.Context) {
	logs, err := lc.logUsecase.GetAllLogs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve logs"})
		return
	}
	c.JSON(http.StatusOK, logs)
}
