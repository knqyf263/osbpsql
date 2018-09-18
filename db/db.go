package db

import (
	"fmt"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/knqyf263/osbpsql/config"
)

type ServiceInstance struct {
	Id    int64
	State string
}

type ServiceBinding struct {
	Id    int64
	State string
}

type Manager struct {
	db *pg.DB
}

func New(conf config.DBConfig) *Manager {
	db := pg.Connect(&pg.Options{
		Addr: fmt.Sprintf("%s:%s", conf.Host, conf.Port),
		User: conf.User,
	})
	return &Manager{db: db}
}

func (m Manager) Migrate() error {
	for _, model := range []interface{}{(*ServiceInstance)(nil), (*ServiceBinding)(nil)} {
		err := m.db.CreateTable(model, &orm.CreateTableOptions{
			IfNotExists: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (m Manager) CreateDatabase(dbName string) error {
	_, err := m.db.Exec(fmt.Sprintf(`CREATE DATABASE "%s"`, dbName))
	return err
}

func (m Manager) DropDatabase(dbName string) error {
	_, err := m.db.Exec(fmt.Sprintf(`DROP DATABASE "%s"`, dbName))
	return err
}
