package assessments

import (
	"context"
	"net/url"
	"reflect"
	"testing"

	"assessments/config"
)

func TestDataSyncWorker_DataSync(t *testing.T) {
	type fields struct {
		s    Service
		repo Repo
		cnf  *config.Config
	}
	type args struct {
		ctx             context.Context
		u               *url.URL
		newIndegoApiUrl string
		city            string
	}
	var tests []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DataSyncWorker{
				s:    tt.fields.s,
				repo: tt.fields.repo,
				cnf:  tt.fields.cnf,
			}
			if err := d.DataSync(tt.args.ctx, tt.args.u, tt.args.newIndegoApiUrl, tt.args.city); (err != nil) != tt.wantErr {
				t.Errorf("DataSync() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewDataSyncWorker(t *testing.T) {
	type args struct {
		s    Service
		repo Repo
		cnf  *config.Config
	}
	var tests []struct {
		name string
		args args
		want *DataSyncWorker
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDataSyncWorker(tt.args.s, tt.args.repo, tt.args.cnf); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDataSyncWorker() = %v, want %v", got, tt.want)
			}
		})
	}
}
