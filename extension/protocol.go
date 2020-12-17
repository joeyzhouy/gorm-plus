package extension

import (
	"database/sql"
	"errors"
	"github.com/joeyzhouy/gorm-plus/utils"
)

type Protocol interface {
	GetConnection() (*sql.DB, error)
	GetTableNames() ([]string, error)
	GetColumnWithTableNames(tableNames []string) ([]utils.TableInfo, error)
}

var protocols = make(map[string]func(param utils.Param) (Protocol, error))

func RegisterProtocol(name string, v func(param utils.Param) (Protocol, error)) error {
	_, ok := protocols[name]
	if ok {
		return errors.New("protocol name: " + name + ", already registered!")
	}
	protocols[name] = v
	return nil
}

func GetProtocol(name string, param utils.Param) (Protocol, error) {
	p, ok := protocols[name]
	if !ok {
		return nil, errors.New("protocol name: " + name + ", not found!")
	}
	return p(param)
}
