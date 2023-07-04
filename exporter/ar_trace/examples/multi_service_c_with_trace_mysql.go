package examples

import (
	"context"
	"database/sql"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporter/ar_trace"
	_ "github.com/go-sql-driver/mysql"
)

type Address struct {
	id    int
	value string
}

var db *sql.DB

func DBInit() error {
	var err error
	db, err = sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/example")
	return err
}

func GetProvince(ctx context.Context, id string) string {
	_, span := ar_trace.Tracer.Start(ctx, "GetProvince")
	defer span.End()
	var address Address
	err := db.QueryRow("select * from address where id = "+id).Scan(&address.id, &address.value)
	if err != nil {
		return err.Error()
	}
	return address.value
}

func GetCity(ctx context.Context, id string) string {
	_, span := ar_trace.Tracer.Start(ctx, "GetCity")
	defer span.End()
	var address Address
	err := db.QueryRow("select * from address where id = "+id).Scan(&address.id, &address.value)
	if err != nil {
		return err.Error()
	}
	return address.value
}

func MockGetProvince(ctx context.Context, id string) string {
	_, span := ar_trace.Tracer.Start(ctx, "MockGetProvince")
	defer span.End()
	return "SiChuan"
}

func MockGetCity(ctx context.Context, id string) string {
	_, span := ar_trace.Tracer.Start(ctx, "MockGetCity")
	defer span.End()
	return "ChengDu"
}
