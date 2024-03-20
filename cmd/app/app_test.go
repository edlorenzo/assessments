package app

import (
	"context"
	"database/sql"
	"testing"

	"assessments/config"
	"assessments/internal/assessments"
	"assessments/server"
)

func TestApp_Run(t *testing.T) {
	type fields struct {
		conf          *config.Config
		assessmentSvc assessments.Service
		server        *server.Server
		cancelFunc    context.CancelFunc
		repoSvc       assessments.Repo
	}
	var tests []struct {
		name    string
		fields  fields
		wantErr bool
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &App{
				conf:          tt.fields.conf,
				assessmentSvc: tt.fields.assessmentSvc,
				server:        tt.fields.server,
				cancelFunc:    tt.fields.cancelFunc,
				repoSvc:       tt.fields.repoSvc,
			}
			if err := a.Run(); (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestApp_Setup(t *testing.T) {
	type fields struct {
		conf          *config.Config
		assessmentSvc assessments.Service
		server        *server.Server
		cancelFunc    context.CancelFunc
		repoSvc       assessments.Repo
	}
	var tests []struct {
		name    string
		fields  fields
		wantErr bool
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &App{
				conf:          tt.fields.conf,
				assessmentSvc: tt.fields.assessmentSvc,
				server:        tt.fields.server,
				cancelFunc:    tt.fields.cancelFunc,
				repoSvc:       tt.fields.repoSvc,
			}
			if err := a.Setup(); (err != nil) != tt.wantErr {
				t.Errorf("Setup() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestApp_closeFN(t *testing.T) {
	type fields struct {
		conf          *config.Config
		assessmentSvc assessments.Service
		server        *server.Server
		cancelFunc    context.CancelFunc
		repoSvc       assessments.Repo
	}
	var tests []struct {
		name    string
		fields  fields
		wantErr bool
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &App{
				conf:          tt.fields.conf,
				assessmentSvc: tt.fields.assessmentSvc,
				server:        tt.fields.server,
				cancelFunc:    tt.fields.cancelFunc,
				repoSvc:       tt.fields.repoSvc,
			}
			if err := a.closeFN(); (err != nil) != tt.wantErr {
				t.Errorf("closeFN() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStartAPP(t *testing.T) {
	var tests []struct {
		name string
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			StartAPP()
		})
	}
}

func Test_runMigration(t *testing.T) {
	type args struct {
		db   *sql.DB
		path string
	}
	var tests []struct {
		name    string
		args    args
		wantErr bool
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RunMigration(tt.args.db, tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("RunMigration() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
