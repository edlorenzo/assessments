package handler

import (
	"testing"

	"assessments/config"
	"assessments/internal/assessments"
	"github.com/gin-gonic/gin"
)

func TestSetupAssessmentsHandler(t *testing.T) {
	type args struct {
		router gin.IRouter
		s      assessments.Service
		repo   assessments.Repo
		cnf    *config.Config
	}
	var tests []struct {
		name string
		args args
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetupAssessmentsHandler(tt.args.router, tt.args.s, tt.args.repo, tt.args.cnf)
		})
	}
}

func Test_assessmentsHandler_ReadinessProbe(t *testing.T) {
	type fields struct {
		s    assessments.Service
		repo assessments.Repo
		cnf  *config.Config
	}
	type args struct {
		c *gin.Context
	}
	var tests []struct {
		name   string
		fields fields
		args   args
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hh := &assessmentsHandler{
				s:    tt.fields.s,
				repo: tt.fields.repo,
				cnf:  tt.fields.cnf,
			}
			hh.ReadinessProbe(tt.args.c)
		})
	}
}

func Test_assessmentsHandler_fetchAndStore(t *testing.T) {
	type fields struct {
		s    assessments.Service
		repo assessments.Repo
		cnf  *config.Config
	}
	type args struct {
		ctx *gin.Context
	}
	var tests []struct {
		name   string
		fields fields
		args   args
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hh := assessmentsHandler{
				s:    tt.fields.s,
				repo: tt.fields.repo,
				cnf:  tt.fields.cnf,
			}
			hh.fetchAndStore(tt.args.ctx)
		})
	}
}

func Test_assessmentsHandler_getStationsData(t *testing.T) {
	type fields struct {
		s    assessments.Service
		repo assessments.Repo
		cnf  *config.Config
	}
	type args struct {
		c *gin.Context
	}
	var tests []struct {
		name   string
		fields fields
		args   args
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hh := assessmentsHandler{
				s:    tt.fields.s,
				repo: tt.fields.repo,
				cnf:  tt.fields.cnf,
			}
			hh.getStationsData(tt.args.c)
		})
	}
}
