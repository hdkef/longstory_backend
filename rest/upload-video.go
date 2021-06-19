package rest

import (
	"errors"
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
	ERROR_OCCURED    = "error occured"
	UPLOAD_COMPLETE  = "upload complete"
	STORE_THUMB      = "store thumb"
	STORE_PATHS_DB   = "store path to db"
	DELETE_VIDEO     = "del video"
	DELETE_THUMBNAIL = "del thumbnail"
	DELETE_BOTH      = "del both"
)

var STATIC_PATH = os.Getenv("STATIC_PATH")
var VIDEOS_PATH = os.Getenv("VIDEOS_PATH")
var THUMBNAILS_PATH = os.Getenv("THUMBNAILS_PATH")

type progress struct {
	Res           *http.ResponseWriter
	Req           *http.Request
	ID            string
	Status        string
	VideoPath     string
	ThumbnailPath string
	Error         error
	DB            *mongo.Client
}

func UploadVideo(db *mongo.Client) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		var progressChan chan progress = make(chan progress)
		var responseChan chan progress = make(chan progress)

		go progressChanRouter(progressChan, responseChan)

		//first action is storing the video
		go storeVid(progress{
			Res: &res,
			Req: req,
			DB:  db,
		}, progressChan)

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
		case STORE_THUMB:
			go storeThumb(progress, progressChan)
		case DELETE_VIDEO:
			go deleteFile(DELETE_VIDEO, progress, progressChan)
		case DELETE_THUMBNAIL:
			go deleteFile(DELETE_THUMBNAIL, progress, progressChan)
		case DELETE_BOTH:
			go deleteFile(DELETE_BOTH, progress, progressChan)
		case STORE_PATHS_DB:
			go storePathsToDB(progress, progressChan)
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

func storeVid(progress progress, progressChan chan progress) {

	filePath, err := storeFile(progress.Req, "video", VIDEOS_PATH)
	if err != nil {
		go sendErrorSignal(&err, &progress, progressChan)
		return
	}

	progress.Status = STORE_THUMB
	progress.VideoPath = filePath
	progressChan <- progress
}

//IN THIS FUNC EVERY ERROR SHOULD send DELETE_VIDEO
func storeThumb(progress progress, progressChan chan progress) {

	filePath, err := storeFile(progress.Req, "thumbnail", THUMBNAILS_PATH)
	if err != nil {
		go sendDeleteSignal(DELETE_VIDEO, &progress, progressChan)
		return
	}

	progress.Status = STORE_PATHS_DB
	progress.ThumbnailPath = filePath
	progressChan <- progress
}

//storeFile will store file and return relative path to file
func storeFile(req *http.Request, formfilename string, relpath string) (string, error) {

	err := req.ParseMultipartForm(1024)
	if err != nil {
		return "", err
	}

	file, fileHeader, err := req.FormFile(formfilename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	absPath, err := getAbsPath()
	if err != nil {
		return "", err
	}

	filename := fileHeader.Filename
	userID := req.FormValue("id")
	folderRelativePath := filepath.Join(STATIC_PATH, relpath, userID)
	folderPath := filepath.Join(absPath, folderRelativePath)
	fileloc := filepath.Join(folderPath, filename)

	err = createFolder(folderPath)
	if err != nil {
		return "", err
	}

	targetFile, err := os.OpenFile(fileloc, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return "", err
	}
	defer targetFile.Close()

	if _, err := io.Copy(targetFile, file); err != nil {
		return "", err
	}

	return filepath.Join(folderRelativePath, filename), nil
}

func getAbsPath() (string, error) {

	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return dir, nil
}

//IN THIS FUNC EVERY ERROR SHOULD send DELETE_VIDEO AND DELETE_THUMBNAIL progress.Status
func storePathsToDB(progress progress, progressChan chan progress) {
	//TOBEIMPLEMENTED
	//store video path to database

	// go sendDeleteSignal(DELETE_BOTH, &progress, progressChan)

	progress.Status = UPLOAD_COMPLETE
	progressChan <- progress
}

func deleteFile(deltype string, progress progress, progressChan chan progress) {

	absPath, err := getAbsPath()
	if err != nil {
		go sendErrorSignal(&err, &progress, progressChan)
		return
	}

	var truePath string

	if deltype == DELETE_VIDEO {
		truePath = filepath.Join(absPath, progress.VideoPath)
		removeFile(truePath)
	} else if deltype == DELETE_THUMBNAIL {
		truePath = filepath.Join(absPath, progress.ThumbnailPath)
		removeFile(truePath)
	} else if deltype == DELETE_BOTH {
		videoPath := filepath.Join(absPath, progress.VideoPath)
		thumbnailPath := filepath.Join(absPath, progress.ThumbnailPath)
		removeFile(videoPath, thumbnailPath)
	}

	err = errors.New("ERROR OCCURED. file has been stored then deleted for some reason")
	go sendErrorSignal(&err, &progress, progressChan)
}

func removeFile(truePaths ...string) error {
	for _, truePath := range truePaths {
		err := os.Remove(truePath)
		if err != nil {
			return err
		}
	}
	return nil
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

func sendDeleteSignal(deltype string, progress *progress, progressChan chan progress) {
	progress.Status = deltype
	progressChan <- *progress
}
