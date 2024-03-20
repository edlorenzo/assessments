package assessments

import (
	"context"
	"reflect"
	"testing"

	"assessments/internal/assessments/model"
	"gorm.io/gorm"
)

func Test_assessmentsService_FetchStationData(t *testing.T) {
	type fields struct {
		repo Repo
		db   *gorm.DB
	}
	type args struct {
		ctx context.Context
		url string
	}
	var tests []struct {
		name    string
		fields  fields
		args    args
		want    *bikeAvailabilityResponse
		wantErr bool
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &assessmentsService{
				repo: tt.fields.repo,
				db:   tt.fields.db,
			}
			got, err := s.FetchStationData(tt.args.ctx, tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchStationData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FetchStationData() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewAssessmentsService(t *testing.T) {
	type args struct {
		repo Repo
		db   *gorm.DB
	}
	var tests []struct {
		name string
		args args
		want Service
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAssessmentsService(tt.args.repo, tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAssessmentsService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_assessmentsService_ConsolidatedData(t *testing.T) {
	type fields struct {
		repo Repo
		db   *gorm.DB
	}
	type args struct {
		ctx          context.Context
		weather      *model.Weather
		availability *model.BikeAvailability
	}
	var tests []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &assessmentsService{
				repo: tt.fields.repo,
				db:   tt.fields.db,
			}
			if err := s.ConsolidatedData(tt.args.ctx, tt.args.weather, tt.args.availability); (err != nil) != tt.wantErr {
				t.Errorf("ConsolidatedData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_assessmentsService_ConsolidatedData2(t *testing.T) {
	type fields struct {
		repo Repo
		db   *gorm.DB
	}
	type args struct {
		ctx          context.Context
		weather      *model.Weather
		availability *model.BikeAvailability
	}
	var tests []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &assessmentsService{
				repo: tt.fields.repo,
				db:   tt.fields.db,
			}
			if err := s.ConsolidatedData2(tt.args.ctx, tt.args.weather, tt.args.availability); (err != nil) != tt.wantErr {
				t.Errorf("ConsolidatedData2() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_assessmentsService_FetchStationData1(t *testing.T) {
	type fields struct {
		repo Repo
		db   *gorm.DB
	}
	type args struct {
		ctx context.Context
		url string
	}
	var tests []struct {
		name    string
		fields  fields
		args    args
		want    *bikeAvailabilityResponse
		wantErr bool
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &assessmentsService{
				repo: tt.fields.repo,
				db:   tt.fields.db,
			}
			got, err := s.FetchStationData(tt.args.ctx, tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchStationData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FetchStationData() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_assessmentsService_FetchWeatherData(t *testing.T) {
	type fields struct {
		repo Repo
		db   *gorm.DB
	}
	type args struct {
		ctx context.Context
		url string
	}
	var tests []struct {
		name    string
		fields  fields
		args    args
		want    *weatherResponse
		wantErr bool
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &assessmentsService{
				repo: tt.fields.repo,
				db:   tt.fields.db,
			}
			got, err := s.FetchWeatherData(tt.args.ctx, tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchWeatherData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FetchWeatherData() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_assessmentsService_HealthCheck(t *testing.T) {
	type fields struct {
		repo Repo
		db   *gorm.DB
	}
	type args struct {
		ctx context.Context
	}
	var tests []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &assessmentsService{
				repo: tt.fields.repo,
				db:   tt.fields.db,
			}
			if err := s.HealthCheck(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("HealthCheck() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_assessmentsService_StationSnapShotByDate(t *testing.T) {
	type fields struct {
		repo Repo
		db   *gorm.DB
	}
	type args struct {
		ctx     context.Context
		kioskId string
		date    string
	}
	var tests []struct {
		name    string
		fields  fields
		args    args
		want    *[]bikeAvailabilityResponses
		wantErr bool
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &assessmentsService{
				repo: tt.fields.repo,
				db:   tt.fields.db,
			}
			got, err := s.StationSnapShotByDate(tt.args.ctx, tt.args.kioskId, tt.args.date)
			if (err != nil) != tt.wantErr {
				t.Errorf("StationSnapShotByDate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StationSnapShotByDate() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_assessmentsService_WeatherSnapShotByDate(t *testing.T) {
	type fields struct {
		repo Repo
		db   *gorm.DB
	}
	type args struct {
		ctx     context.Context
		kioskId string
		date    string
	}
	var tests []struct {
		name    string
		fields  fields
		args    args
		want    *[]weatherResponse
		wantErr bool
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &assessmentsService{
				repo: tt.fields.repo,
				db:   tt.fields.db,
			}
			got, err := s.WeatherSnapShotByDate(tt.args.ctx, tt.args.kioskId, tt.args.date)
			if (err != nil) != tt.wantErr {
				t.Errorf("WeatherSnapShotByDate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WeatherSnapShotByDate() got = %v, want %v", got, tt.want)
			}
		})
	}
}
