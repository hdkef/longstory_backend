package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"longstory/graph/generated"
	"longstory/graph/model"
	"longstory/helper"
)

func (r *mutationResolver) Login(ctx context.Context, input *model.NewLogin) (*model.Token, error) {
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

func (r *mutationResolver) Autologin(ctx context.Context, input *model.NewAutoLogin) (*model.Token, error) {
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
		return nil, err
	}
}

func (r *mutationResolver) Hotspotvid(ctx context.Context, input *model.Paging) ([]*model.Video, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) Foryouvid(ctx context.Context, input *model.Paging) ([]*model.Video, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Videos(ctx context.Context) ([]*model.Video, error) {
	panic(fmt.Errorf("not implemented"))
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
func (r *mutationResolver) Videos(ctx context.Context, input *model.Paging) ([]*model.Video, error) {
	panic(fmt.Errorf("not implemented"))
}
