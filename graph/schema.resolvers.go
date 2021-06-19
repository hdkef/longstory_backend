package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"longstory/graph/generated"
	"longstory/graph/model"
	"longstory/helper"
	"longstory/mock"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	ERR_EMAIL_EXIST = "error email already exist"
	DB_NAME         = "longstory"
	USERS_DOC       = "users"
	VIDEOS_DOC      = "videos"
)

func (r *mutationResolver) Delete(ctx context.Context, id string) (*model.Status, error) {
	//TOBEIMPLEMENT
	//delete video from database and delete video file too
	return &model.Status{
		Status: true,
	}, nil
}

func (r *queryResolver) Hotspotvideos(ctx context.Context, id string) ([]*model.Video, error) {
	//TOBEIMPLEMENT
	//Implements getting hostpotvideos paginated by lastID
	mockvids := mock.MockVideos[0:6]
	/////////////////////////////////////////////////////
	return mockvids, nil
}

func (r *queryResolver) Foryouvideos(ctx context.Context, id string) ([]*model.Video, error) {
	//TOBEIMPLEMENT
	//Implements getting foryouvideos paginated by lastID
	mockvids := mock.MockVideos[0:6]
	/////////////////////////////////////////////////////
	return mockvids, nil
}

func (r *queryResolver) Login(ctx context.Context, input *model.NewLogin) (*model.Token, error) {
	//TOBEIMPLEMENT
	//GET PASS FROM DATABASE AND COMPARE
	//IF SAME, THEN CREATE TOKEN
	user := model.User{
		ID:        "1",
		Username:  "Agus",
		Avatarurl: "no_avatar",
	}
	token, err := helper.CreateToken(&user)
	if err != nil {
		return &model.Token{}, err
	}
	return &model.Token{
		Type:  "new",
		Token: token,
	}, nil
}

func (r *queryResolver) Autologin(ctx context.Context, input *model.NewAutoLogin) (*model.Token, error) {
	parsedToken, err := helper.ParseTokenString(&input.Token)
	if err == nil {
		return nil, nil
	} else if err.Error() == helper.ERR_NEED_NEW_TOKEN {
		//access DB and generate new token
		user, err := helper.ParseMapClaims(parsedToken)
		if err != nil {
			return &model.Token{}, err
		}
		token, err := helper.CreateToken(user)
		if err != nil {
			return &model.Token{}, err
		}
		return &model.Token{
			Type:  "refresh",
			Token: token,
		}, nil
	} else {
		fmt.Println(err)
		return nil, err
	}
}

func (r *queryResolver) CheckUsername(ctx context.Context, input *model.Email) (*model.Status, error) {
	col := r.DB.Database(DB_NAME).Collection(USERS_DOC)
	var user model.User
	filter := bson.D{{Key: "email", Value: input.Email}}
	err := col.FindOne(ctx, filter).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return &model.Status{
			Status: true,
		}, nil
	} else if err != nil {
		return &model.Status{}, nil
	}
	return &model.Status{Status: false}, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
