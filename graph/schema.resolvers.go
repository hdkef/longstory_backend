package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"longstory/graph/generated"
	"longstory/graph/model"
	"longstory/helper"

	"go.mongodb.org/mongo-driver/mongo"
)

func (r *mutationResolver) Delete(ctx context.Context, input *model.DeleteVid) (*model.Status, error) {
	err := helper.DeleteFile(input)
	if err != nil {
		return &model.Status{}, err
	}
	err = r.DBRepo.DeleteByID(ctx, COL_VIDEOS, input.ID)
	if err != nil {
		return &model.Status{}, err
	}

	return &model.Status{
		Status: true,
	}, nil
}

func (r *queryResolver) Hotspotvideos(ctx context.Context, id string) ([]*model.Video, error) {
	//TOBEIMPLEMENTEDSORT
	result, err := r.DBRepo.FindManyByLastID(ctx, COL_VIDEOS, id, 6)
	if err != nil {
		return nil, err
	}
	var videos []*model.Video = result.([]*model.Video)

	return videos, nil
}

func (r *queryResolver) Foryouvideos(ctx context.Context, id string) ([]*model.Video, error) {
	//TOBEIMPLEMENTEDSORT
	result, err := r.DBRepo.FindManyByLastID(ctx, COL_VIDEOS, id, 6)
	if err != nil {
		return nil, err
	}
	var videos []*model.Video = result.([]*model.Video)

	return videos, nil
}

func (r *queryResolver) Login(ctx context.Context, input *model.NewLogin) (*model.Token, error) {
	result, err := r.DBRepo.FindOneByField(ctx, COL_USERS, COL_USERS_FIELD_USERNAME, input.Username)
	if err != nil {
		return &model.Token{}, err
	}

	user := result.(model.User)

	if *user.Pass != input.Pass {
		return &model.Token{}, errors.New(ERR_PASS_NOT_MATCH)
	}

	token, err := helper.CreateToken(&user)
	if err != nil {
		return &model.Token{}, err
	}
	return &model.Token{
		User:  &user,
		Type:  TOKEN_TYPE_NEW,
		Token: &token,
	}, nil
}

func (r *queryResolver) Autologin(ctx context.Context, input *model.NewAutoLogin) (*model.Token, error) {
	parsedToken, err := helper.ParseTokenString(&input.Token)
	if err == nil {
		user, err := helper.ParseMapClaims(parsedToken)
		if err != nil {
			return &model.Token{
				Type: TOKEN_TYPE_CLEAR,
			}, nil
		}
		return &model.Token{
			User: user,
			Type: TOKEN_TYPE_NO_NEW,
		}, nil
	} else if err.Error() == helper.ERR_NEED_NEW_TOKEN {
		user, err := helper.ParseMapClaims(parsedToken)
		if err != nil {
			return &model.Token{
				Type: TOKEN_TYPE_CLEAR,
			}, nil
		}
		token, err := helper.CreateToken(user)
		if err != nil {
			return &model.Token{
				Type: TOKEN_TYPE_CLEAR,
			}, nil
		}
		return &model.Token{
			User:  user,
			Type:  TOKEN_TYPE_REFRESH,
			Token: &token,
		}, nil
	} else {
		return &model.Token{
			Type: TOKEN_TYPE_CLEAR,
		}, nil
	}
}

func (r *queryResolver) Checkemail(ctx context.Context, input *model.ChekEmail) (*model.Status, error) {
	_, err := r.DBRepo.FindOneByField(ctx, COL_USERS, COL_USERS_FIELD_EMAIL, input.Email)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &model.Status{
				Status: true,
			}, nil
		}
		return &model.Status{}, err
	}
	return &model.Status{Status: false}, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
const (
	ERR_EMAIL_EXIST          = "error email already exist"
	ERR_PASS_NOT_MATCH       = "pass not match"
	DB_NAME                  = "longstory"
	USERS_DOC                = "users"
	VIDEOS_DOC               = "videos"
	TOKEN_TYPE_REFRESH       = "refresh"
	TOKEN_TYPE_CLEAR         = "clear"
	TOKEN_TYPE_NEW           = "new"
	TOKEN_TYPE_NO_NEW        = "nonew"
	COL_VIDEOS               = "videos"
	COL_USERS                = "users"
	COL_USERS_FIELD_EMAIL    = "email"
	COL_USERS_FIELD_USERNAME = "username"
)
