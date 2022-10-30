package datastore

import (
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/logrusadapter"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/sirupsen/logrus"
	"os"
)

type DBStore struct {
	Logger     *logrus.Logger
	DB         *sql.DB
	SQLBuilder squirrel.StatementBuilderType
}

// NewDBStore open a connection to db
func NewDBStore(logger *logrus.Logger) (*DBStore, error) {
	var (
		err  error
		conn *pgx.ConnConfig
	)

	host := os.Getenv("HOST")
	dbUSER := os.Getenv("DB_USERNAME")
	pass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	src := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUSER, pass, host, dbPort, dbName)
	logrus.Println(src)

	conn, err = pgx.ParseConfig(src)
	if err != nil {
		return nil, err
	}
	conn.Logger = logrusadapter.NewLogger(logger)
	db := stdlib.OpenDB(*conn)
	err = validateSchema(db)
	if err != nil {
		return nil, err
	}

	return &DBStore{
		Logger:     logger,
		DB:         db,
		SQLBuilder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).RunWith(db),
	}, err
}

//Close connection
func (d *DBStore) Close() error {
	return d.DB.Close()
}
