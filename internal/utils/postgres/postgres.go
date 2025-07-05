package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
	sqltrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/database/sql"
)

type Postgres struct {
	DB *sql.DB
}

type Config struct {
	ServiceName   string // trace service name
	Dsn           string
	MaxConn       int
	MaxIdle       int
	DataDogTracer bool //trace to be on or off
}

const (
	// postgres driver name
	postgres = "postgres"
)

func New(cfg *Config) (*Postgres, error) {
	var db *sql.DB
	var err error
	if !cfg.DataDogTracer {
		db, err = sql.Open(postgres, cfg.Dsn)
		if err != nil {
			return nil, err
		}

	} else {
		sqltrace.Register(postgres, &pq.Driver{})

		db, err = sqltrace.Open(postgres, cfg.Dsn, sqltrace.WithServiceName(cfg.ServiceName))
		if err != nil {
			return nil, err
		}
	}

	db.SetMaxOpenConns(cfg.MaxConn)
	db.SetMaxIdleConns(cfg.MaxIdle)

	//ping db to ensure the connection is alive and working
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return &Postgres{db}, nil
}

func (p *Postgres) Close() {
	p.DB.Close()
}
