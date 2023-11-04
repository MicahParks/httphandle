package postgres

import (
	"context"
	"fmt"
	"time"

	jt "github.com/MicahParks/jsontype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	DSN                   string                      `json:"dsn"`
	Health                *jt.JSONType[time.Duration] `json:"health"`
	InitialTimeout        *jt.JSONType[time.Duration] `json:"initialTimeout"`
	MaxIdle               *jt.JSONType[time.Duration] `json:"maxIdle"`
	MaxConnLifetime       *jt.JSONType[time.Duration] `json:"maxConnLifetime"`
	MaxConnLifetimeJitter *jt.JSONType[time.Duration] `json:"maxConnLifetimeJitter"`
	MinConns              int32                       `json:"minConns"`
}

func (c Config) DefaultsAndValidate() (Config, error) {
	if c.DSN == "" {
		return c, fmt.Errorf("%w: dsn is required", jt.ErrDefaultsAndValidate)
	}
	if c.Health.Get() == 0 {
		c.Health = jt.New(5 * time.Second)
	}
	if c.InitialTimeout.Get() == 0 {
		c.InitialTimeout = jt.New(10 * time.Second)
	}
	if c.MaxIdle.Get() == 0 {
		c.MaxIdle = jt.New(5 * time.Minute)
	}
	if c.MaxConnLifetime.Get() == 0 {
		c.MaxConnLifetime = jt.New(30 * time.Minute)
	}
	if c.MaxConnLifetimeJitter.Get() == 0 {
		c.MaxConnLifetimeJitter = jt.New(5 * time.Minute)
	}
	if c.MinConns == 0 {
		c.MinConns = 2
	}
	return c, nil
}

func Pool(ctx context.Context, config Config) (*pgxpool.Pool, error) {
	c, err := pgxpool.ParseConfig(config.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to parse PostgreSQL DSN: %w", err)
	}
	c.HealthCheckPeriod = config.Health.Get()
	c.MaxConnIdleTime = config.MaxIdle.Get()
	c.MaxConnLifetime = config.MaxConnLifetime.Get()
	c.MaxConnLifetimeJitter = config.MaxConnLifetimeJitter.Get()
	c.MinConns = config.MinConns

	var conn *pgxpool.Pool
	const retries = 5
	for i := 0; i < retries; i++ {
		conn, err = pgxpool.NewWithConfig(ctx, c)
		if err != nil {
			return nil, fmt.Errorf("failed to create pool with given configuration: %w", err)
		}
		err = conn.Ping(ctx)
		if err != nil {
			select {
			case <-ctx.Done():
				return nil, fmt.Errorf("failed to connect to Postgres after waiting %d seconds: %w", retries, err)
			case <-time.After(time.Second):
			}
			continue
		}
		break
	}
	if err != nil {
		return nil, fmt.Errorf("failed to create database connection: %w", err)
	}

	return conn, nil

}
