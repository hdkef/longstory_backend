package rest

import (
	"fmt"
	"longstory/utils"
	"mime/multipart"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
)

const (
	ERROR_OCCURED   = "error occured"
	DELETE_FILE     = "delete file"
	UPLOAD_COMPLETE = "upload complete"
	STORE_PATH_DB   = "store path to db"
	CONVERT_TO_HLS  = "convert to hls"
)

type progress struct {
	Res    *http.ResponseWriter
	ID     string
	Status string
	Path   string
	Error  error
	DB     *mongo.Client
}

func UploadVideo(db *mongo.Client) http.HandlerFunc {
	return func(res http.ResponseWriter, r *http.Request) {

		file, fileHeader, err := r.FormFile("video")
		if err != nil {
			utils.ResErr(&res, http.StatusInternalServerError, err)
			return
		}

		var progressChan chan progress = make(chan progress)
		var responseChan chan progress = make(chan progress)

		go progressChanRouter(progressChan, responseChan)

		go storeFile(progress{
			Res: &res,
			DB:  db,
		}, file, *fileHeader, progressChan)

		for progress := range responseChan {
			switch progress.Status {
			case ERROR_OCCURED:
				utils.ResErr(progress.Res, http.StatusInternalServerError, progress.Error)
				close(responseChan)
				return
			case UPLOAD_COMPLETE:
				utils.ResOK(progress.Res, "OK")
				close(responseChan)
				return
			}
		}
	}
}

func progressChanRouter(progressChan chan progress, responseChan chan progress) {
	for progress := range progressChan {
		switch progress.Status {
		case DELETE_FILE:
			go deleteFile(progress, progressChan)
		case STORE_PATH_DB:
			go storePathToDB(progress, progressChan)
		case CONVERT_TO_HLS:
			go convertToHLS(progress, progressChan)
		case UPLOAD_COMPLETE:
			responseChan <- progress
			close(progressChan)
			return
		case ERROR_OCCURED:
			responseChan <- progress
			close(progressChan)
			return
		}
	}
}

func storeFile(progress progress, file multipart.File, fileHeader multipart.FileHeader, progressChan chan progress) {

	absPath, err := getAbsPath()
	if err != nil {
		progress.Status = ERROR_OCCURED
		progress.Error = err
		progressChan <- progress
		return
	}
	fmt.Println("abspath : ", absPath)
	//TOBEIMPLEMENTED
	//STORE FILE TO DISK AND RETURN PATH

	progress.Status = CONVERT_TO_HLS
	progress.Path = "path from file stored"
	progressChan <- progress
}

func getAbsPath() (string, error) {

	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return dir, nil
}

func convertToHLS(progress progress, progressChan chan progress) {
	//TOBEIMPLEMENTED
	//FFMPEG output format is hls

	progress.Status = STORE_PATH_DB
	progress.Path = "path from hls"
	progressChan <- progress
}

func storePathToDB(progress progress, progressChan chan progress) {
	//TOBEIMPLEMENTED
	//store hls path to database

	progress.Status = UPLOAD_COMPLETE
	progressChan <- progress
}

func deleteFile(progress progress, progressChan chan progress) {

}
