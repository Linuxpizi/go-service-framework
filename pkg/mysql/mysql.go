package mysql

import (
	"fmt"
	"time"

	"golang.org/x/xerrors"

	"database/sql" //nolint

	entsql "entgo.io/ent/dialect/sql"

	_ "github.com/go-sql-driver/mysql" //nolint

	"github.com/NpoolPlatform/go-service-framework/pkg/config"
	constant "github.com/NpoolPlatform/go-service-framework/pkg/mysql/const" //nolint
)

type Client struct {
	Driver *entsql.Driver
}

const (
	keyUsername = "username"
	keyPassword = "password"
	keyDBName   = "database_name"
)

func NewMysqlClient() (*Client, error) {
	service, err := config.PeekService(constant.MysqlServiceName)
	if err != nil {
		return nil, xerrors.Errorf("Fail to query mysql service: %v", err)
	}

	username := config.GetStringValueWithNameSpace(constant.MysqlServiceName, keyUsername)
	password := config.GetStringValueWithNameSpace(constant.MysqlServiceName, keyPassword)
	myServiceName := config.GetStringValueWithNameSpace("", config.KeyHostname)
	dbname := config.GetStringValueWithNameSpace(myServiceName, keyDBName)

	if username == "" {
		return nil, xerrors.Errorf("invalid username")
	}
	if password == "" {
		return nil, xerrors.Errorf("invalid password")
	}
	if myServiceName == "" {
		return nil, xerrors.Errorf("invalid service name")
	}
	if dbname == "" {
		return nil, xerrors.Errorf("invalid dbname")
	}

	dsl := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=true&interpolateParams=true", username, password, service.Address, service.Port, dbname)
	db, err := sql.Open("mysql", dsl)
	if err != nil {
		return nil, xerrors.Errorf("Fail to initialize sql driver: %v", err)
	}

	// https://github.com/go-sql-driver/mysql
	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return &Client{
		Driver: entsql.OpenDB("mysql", db),
	}, nil
}
