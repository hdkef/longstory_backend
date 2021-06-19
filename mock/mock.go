package mock

import "longstory/graph/model"

var mockvid *model.Video = &model.Video{
	ID:        "videoID",
	Thumbnail: "thumbnail",
	Link:      "link",
	Title:     "title",
	User: &model.User{
		ID:        "userID",
		Username:  "agus",
		Avatarurl: "avatarurl",
		Email:     "email",
	},
}

var MockVideos []*model.Video = []*model.Video{mockvid, mockvid, mockvid, mockvid, mockvid, mockvid, mockvid, mockvid, mockvid, mockvid, mockvid, mockvid, mockvid}
