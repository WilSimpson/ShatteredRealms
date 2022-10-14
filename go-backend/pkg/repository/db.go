package repository

import (
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	_ "github.com/lib/pq"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
	"io/ioutil"
	"time"
)

// DBConfig Information on how to connect to the database
type DBConfig struct {
	Host     string `yaml:"hoster"`
	Port     string `yaml:"port"`
	Name     string `yaml:"name"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// DBPoolConfig Defines the master and slave connections to a replicated database. Slaves may be empty.
type DBPoolConfig struct {
	Master DBConfig   `yaml:"master"`
	Slaves []DBConfig `yaml:"slaves"`
}

// Connect Initializes the connection to a Postgres database
func Connect(pool DBPoolConfig) (*gorm.DB, error) {
	config, err := pgx.ParseConfig(pool.Master.PostgresDSN())
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	sqlDB := stdlib.OpenDB(*config)

	sqlDB.SetConnMaxLifetime(time.Second)
	sqlDB.SetMaxOpenConns(0)
	sqlDB.SetMaxIdleConns(10)
	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})

	if err != nil {
		return nil, fmt.Errorf("gorm: %w", err)
	}

	if len(pool.Slaves) > 0 {
		replicas := make([]gorm.Dialector, len(pool.Slaves))
		for _, slave := range pool.Slaves {
			replicas = append(replicas, postgres.Open(slave.PostgresDSN()))
		}

		err = db.Use(dbresolver.Register(dbresolver.Config{
			Replicas: replicas,
			Policy:   dbresolver.RandomPolicy{},
		}))
	}

	if err := db.Use(otelgorm.NewPlugin(otelgorm.WithDBName(pool.Master.Name))); err != nil {
		return nil, err
	}

	return db, err
}

func (c DBConfig) MySQLDSN() string {
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		c.Name,
	)
}

func (c DBConfig) PostgresDSN() string {
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		c.Name,
	)
}

func ConnectFromFile(filePath string) (*gorm.DB, error) {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	c := &DBPoolConfig{}
	err = yaml.Unmarshal(file, c)
	if err != nil {
		return nil, err
	}

	return Connect(*c)
}
