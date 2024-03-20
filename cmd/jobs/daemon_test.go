package jobs

import (
	"context"
	"reflect"
	"testing"

	"assessments/config"
	"assessments/internal/assessments"
)

func TestDataSyncDaemon_StartReconcilerDaemon(t *testing.T) {
	type fields struct {
		s    assessments.Service
		repo assessments.Repo
		cnf  *config.Config
	}
	type args struct {
		ctx context.Context
	}
	var tests []struct {
		name   string
		fields fields
		args   args
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DataSyncDaemon{
				s:    tt.fields.s,
				repo: tt.fields.repo,
				cnf:  tt.fields.cnf,
			}
			d.StartReconcilerDaemon(tt.args.ctx)
		})
	}
}

func TestNewDataSyncDaemon(t *testing.T) {
	type args struct {
		s    assessments.Service
		repo assessments.Repo
		cnf  *config.Config
	}
	var tests []struct {
		name string
		args args
		want *DataSyncDaemon
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDataSyncDaemon(tt.args.s, tt.args.repo, tt.args.cnf); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDataSyncDaemon() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetupDataSyncWorker(t *testing.T) {
	var tests []struct {
		name string
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetupDataSyncWorker()
		})
	}
}
