package helper

import (
	"context"
	"fmt"
	"longstory/graph/model"

	"go.mongodb.org/mongo-driver/mongo"
)

type DBRepo struct {
	DB *mongo.Client
}

func (r *DBRepo) DeleteByID(ctx context.Context, string, id string) error {
	fmt.Println("DeleteByID")
	return nil
}

func (r *DBRepo) StoreOne(ctx context.Context, colname string, data interface{}) error {
	fmt.Println("StoreOne")
	return nil
}

func (r *DBRepo) FindOneByID(ctx context.Context, colname string, id string) (interface{}, error) {
	fmt.Println("FindOneByID")
	return nil, nil
}

func (r *DBRepo) FindOneByField(ctx context.Context, colname string, field string, data string) (interface{}, error) {
	fmt.Println("FindOneByID")

	user := model.User{
		ID:        "1",
		Username:  "agus",
		Pass:      StrPointer("pass"),
		Avatarurl: "avatarurl",
	}

	return user, nil
}

func (r *DBRepo) FindManyByLastID(ctx context.Context, colname string, lastid string, limit int64) (interface{}, error) {
	fmt.Println("FindManyByLastID")

	return nil, nil
}

func StrPointer(s string) *string {
	return &s
}
