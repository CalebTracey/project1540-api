package postgres

import (
	"context"
	"go.uber.org/mock/gomock"
	"project1540-api/external/models/postgres"
	postgres2 "project1540-api/internal/dao/postgres/mock"
	"reflect"
	"testing"
)

func TestService_SearchFilesByTag(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx  context.Context
		tags []string
	}
	tests := []struct {
		name    string
		args    args
		mockErr error
		want    postgres.FileResponse
	}{
		{
			name: "Happy Path",
			args: args{
				ctx:  context.Background(),
				tags: []string{"TEST"},
			},
			mockErr: nil,
			want:    postgres.FileResponse{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockDAO := postgres2.NewMockIDAO(ctrl)

			s := Service{
				PSQLDAO: mockDAO,
			}

			mockDAO.EXPECT().SearchFilesByTag(tt.args.ctx, SearchByTagQuery, tt.args.tags).
				Return(tt.want.Files, tt.mockErr).MaxTimes(1)

			if got := s.SearchFilesByTag(tt.args.ctx, tt.args.tags); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchFilesByTag() = %v, want %v", got, tt.want)
			}
		})
	}
}
