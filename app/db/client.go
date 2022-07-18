package db

import (
	"github.com/jmoiron/sqlx"
	yaml "gopkg.in/yaml.v2"
	"log"
)

type Client interface {
	Get(id string, tableName string) interface{}
	GetList(id string, tableName string) []interface{}
	Create(id string, object interface{}) error
}

type DatabaseConn struct {
	*sqlx.DB
}

/*
	client_name: "sqlx"
  	driver: "mysql"
  	dsn: "root:@/database?parseTime=true"
  	db_conn_max_life_time: "20s"
*/
type DatabaseConfig struct {
	driver          string `yaml:"driver"`
	dsn             string `yaml:"dsn"`
	ConnMaxLifeTime string `yaml:"db_Conn_Max_Life_Time"`
}

func NewConfig(path string) (*DatabaseConfig, error) {
	config := &DatabaseConfig{}
	yaml.Unmarshal([]byte(""), config)

	return config, nil
}

func NewDatabaseConn(config *DatabaseConfig) (Client, error) {
	db, err := sqlx.Connect(config.driver, config.dsn)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	return &DatabaseConn{db}, nil
}

func (db *DatabaseConn) Get(id string, tableName string) interface{} {
	return nil
}

func (db *DatabaseConn) GetList(id string, tableName string) []interface{} {

	return nil
}

func (db *DatabaseConn) Create(id string, object interface{}) error {
	return nil
}
