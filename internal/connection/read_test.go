package connection

import (
	"context"
	"github.com/pashagolub/pgxmock"
	"reflect"
	"testing"
)

func TestClient_Read(t *testing.T) {

	ctx := context.Background()

	mock, err := pgxmock.NewConn()

	articleMockRows := mock.NewRows([]string{"raw", "display"}).AddRow("USD", "BTC")

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close(context.Background())

	mock.ExpectQuery("SELECT raw, display FROM crypto.quotes").WithArgs("USD", "BTC").WillReturnRows(articleMockRows).RowsWillBeClosed()

	type fields struct {
		ctx  *context.Context
		pool DB
	}
	type args struct {
		query string
		args  []interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string]string
		wantErr bool
	}{
		{
			name: "test 1",
			fields: fields{
				ctx:  &ctx,
				pool: mock,
			},
			args: args{
				query: "SELECT raw, display FROM crypto.quotes WHERE fsyms=? and tsyms=?",
				args:  []interface{}{"USD", "BTC"},
			},
			want:    map[string]string{"raw": "USD", "display": "BTC"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				ctx:  tt.fields.ctx,
				pool: tt.fields.pool,
			}
			got, err := c.Read(tt.args.query, tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Read() got = %v, want %v", got, tt.want)
			}
		})
	}
}
