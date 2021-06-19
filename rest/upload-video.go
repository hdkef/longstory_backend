package rest

import (
	"fmt"
	"io"
	"longstory/utils"
	"net/http"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	godotenv.Load()
}

const (
	ERROR_OCCURED   = "error occured"
	DELETE_FILE     = "delete file"
	UPLOAD_COMPLETE = "upload complete"
	STORE_PATH_DB   = "store path to db"
	CONVERT_TO_HLS  = "convert to hls"
)

var STATIC_PATH = os.Getenv("STATIC_PATH")
var VIDEOS_PATH = os.Getenv("VIDEOS_PATH")

type progress struct {
	Res    *http.ResponseWriter
	ID     string
	Status string
	Path   string
	Error  error
	DB     *mongo.Client
}

func UploadVideo(db *mongo.Client) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		var progressChan chan progress = make(chan progress)
		var responseChan chan progress = make(chan progress)

		go progressChanRouter(progressChan, responseChan)

		go storeFile(progress{
			Res: &res,
			DB:  db,
		}, req, progressChan)

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

func storeFile(progress progress, req *http.Request, progressChan chan progress) {

	err := req.ParseMultipartForm(1024)
	if err != nil {
		go sendErrorSignal(&err, &progress, progressChan)
		return
	}

	file, fileHeader, err := req.FormFile("video")
	if err != nil {
		go sendErrorSignal(&err, &progress, progressChan)
		return
	}
	defer file.Close()

	absPath, err := getAbsPath()
	if err != nil {
		go sendErrorSignal(&err, &progress, progressChan)
		return
	}

	filename := fileHeader.Filename
	userID := req.FormValue("id")
	folderPath := filepath.Join(absPath, STATIC_PATH, VIDEOS_PATH, userID)
	fileloc := filepath.Join(folderPath, filename)

	err = createFolder(folderPath)
	if err != nil {
		go sendErrorSignal(&err, &progress, progressChan)
		return
	}

	targetFile, err := os.OpenFile(fileloc, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		go sendErrorSignal(&err, &progress, progressChan)
		return
	}
	defer targetFile.Close()

	if _, err := io.Copy(targetFile, file); err != nil {
		go sendErrorSignal(&err, &progress, progressChan)
		return
	}

	progress.Status = CONVERT_TO_HLS
	progress.Path = fmt.Sprintf(VIDEOS_PATH, filename)
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

func createFolder(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return err
		} else {
			return nil
		}
	} else {
		return nil
	}
}

func sendErrorSignal(err *error, progress *progress, progressChan chan progress) {
	progress.Status = ERROR_OCCURED
	progress.Error = *err
	progressChan <- *progress
}
