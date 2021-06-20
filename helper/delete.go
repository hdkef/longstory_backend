package helper

import (
	"longstory/graph/model"
	"os"
	"path/filepath"
)

func DeleteFile(input *model.DeleteVid) error {

	absPath, err := GetAbsPath()
	if err != nil {
		return err
	}
	videoPath := filepath.Join(absPath, input.Video)
	thumbnailPath := filepath.Join(absPath, input.Thumbnail)

	err = RemoveFile(videoPath, thumbnailPath)
	if err != nil {
		return err
	}

	return nil
}

func RemoveFile(truePaths ...string) error {
	for _, truePath := range truePaths {
		err := os.Remove(truePath)
		if err != nil {
			return err
		}
	}
	return nil
}
