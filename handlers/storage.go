package handlers

import (
	"bytes"
	"fileupbackendv2/internal/dirManager"
	"fileupbackendv2/internal/storage"
	"fileupbackendv2/pkg/logging"
	"fileupbackendv2/pkg/response"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func GetFilesHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.Header.Get("userId")
	files, err := storage.GetFiles(userId)
	if err != nil {
		logging.Errorlogger.Error().Msg(err.Error())
		response.InternalServerErrorResponse(w, err.Error())
		return
	}
	response.JSONResponse(w, files)
}

func UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	file, fileHeader, err := r.FormFile("file")
	username := r.Header.Get("username")
	folderPath := r.FormValue("folder")
	if err != nil {
		logging.Errorlogger.Error().Msg(err.Error())
		response.BadRequestResponse(w, err.Error())
		return
	}
	defer file.Close()

	err = storage.SaveFile(file, fileHeader, username, folderPath)
	if err != nil {
		logging.Errorlogger.Error().Msg(err.Error())
		response.InternalServerErrorResponse(w, err.Error())
		return
	}
	response.SuccessResponse(w, "File uploaded successfully")
}

func CreateFolderHandler(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("username")
	folderPath := r.FormValue("folder")
	if folderPath == "" {
		response.BadRequestResponse(w, "Folder path is required")
		return
	}
	err := dirManager.CreateFolder(username, folderPath)
	if err != nil {
		logging.Errorlogger.Error().Msg(err.Error())
		response.InternalServerErrorResponse(w, err.Error())
		return
	}
	response.SuccessResponse(w, "Folder created successfully")
}

func GetSignedUrlHandler(w http.ResponseWriter, r *http.Request) {
	filepath := mux.Vars(r)["filepath"]
	userId := r.Header.Get("userId")
	signedUrl, err := storage.GetSignedUrl(filepath, userId)
	// ownership check please
	if err != nil {
		logging.Errorlogger.Error().Msg(err.Error())
		response.InternalServerErrorResponse(w, err.Error())
		return
	}
	response.SuccessResponse(w, signedUrl)
}

func DownloadFileHandler(w http.ResponseWriter, r *http.Request) {
	compressed := false
	filepath := mux.Vars(r)["filepath"]
	signature := r.URL.Query().Get("signature")
	expiresAtStr := r.URL.Query().Get("expires_at")
	compressedStr := r.URL.Query().Get("compressed")
	if compressedStr == "true" {
		compressed = true
	}
	expiresAt, err := strconv.Atoi(expiresAtStr)
	if err != nil {
		logging.Errorlogger.Error().Msg(err.Error())
		response.BadRequestResponse(w, err.Error())
		return
	}
	if time.Now().After(time.Unix(int64(expiresAt), 0)) {
		response.UnauthorizedResponse(w, "Signature expired")
		return
	}
	f, err := storage.DownloadFile(filepath, signature, expiresAtStr, compressed)
	if err != nil {
		logging.Errorlogger.Error().Msg(err.Error())
		response.InternalServerErrorResponse(w, err.Error())
		return
	}
	http.ServeContent(w, r, filepath, time.Now(), bytes.NewReader(f))
}
