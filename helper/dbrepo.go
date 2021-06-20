package helper

import (
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
)

type DBRepo struct {
	DB *mongo.Client
}

func (r *DBRepo) DeleteByID(colname string, id string) error {
	fmt.Println("DeleteByID")
	return nil
}

func (r *DBRepo) StoreOne(colname string, data interface{}) error {
	fmt.Println("StoreOne")
	return nil
}

func (r *DBRepo) FindOneByID(colname string, id string) (interface{}, error) {
	fmt.Println("FindOneByID")
	return nil, nil
}

func (r *DBRepo) FindOneByField(colname string, field string, data string) (interface{}, error) {
	fmt.Println("FindOneByID")
	return nil, nil
}

func (r *DBRepo) FindManyByLastID(colname string, lastid string, limit int64) ([]interface{}, error) {
	fmt.Println("FindManyByLastID")
	return []interface{}{}, nil
}
