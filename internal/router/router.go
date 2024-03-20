package router

import (
	"assessments/config"
	svc "assessments/internal/assessments"
	"assessments/internal/assessments/handler"

	"github.com/gin-gonic/gin"
)

func Route(r *gin.Engine, s svc.Service, repo svc.Repo, cnf *config.Config) {
	handler.SetupAssessmentsHandler(r, s, repo, cnf)
}
