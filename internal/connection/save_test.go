package connection

import (
	"context"
	"github.com/pashagolub/pgxmock"
	"regexp"
	"testing"
)

func TestClient_Save(t *testing.T) {

	ctx := context.Background()

	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close(context.Background())

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO crypto.quotes")).WillReturnResult(pgxmock.NewResult("UPDATE", 1))

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
		wantErr bool
	}{
		{
			name: "test1 проверка сохранения, должно быть без ошибок",
			fields: fields{
				ctx:  &ctx,
				pool: mock,
			},
			args: args{
				query: "INSERT INTO crypto.quotes",
				args:  []interface{}{"USD", "BTC"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				ctx:  tt.fields.ctx,
				pool: tt.fields.pool,
			}
			if err := c.Save(tt.args.query, tt.args.args); (err != nil) != tt.wantErr {
				t.Errorf("Save() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
