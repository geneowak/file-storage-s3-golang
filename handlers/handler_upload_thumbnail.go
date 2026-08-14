package handlers

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"

	"github.com/geneowak/file-storage-s3-golang/internal/auth"
	"github.com/google/uuid"
)

func (cfg *ApiConfig) handlerUploadThumbnail(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.JwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	videoIDString := r.PathValue("videoID")
	videoID, err := uuid.Parse(videoIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID", err)
		return
	}

	video, err := cfg.DB.GetVideo(videoID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't find video", err)
		return
	}

	if video.UserID != userID {
		respondWithError(w, http.StatusUnauthorized, "You don't own this video", err)
		return
	}
	fmt.Println("uploading thumbnail for video", videoID, "by user", userID)

	const maxMemory = 10 << 20
	r.ParseMultipartForm(maxMemory)
	file, header, err := r.FormFile("thumbnail")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to parse form file", err)
		return
	}
	defer file.Close()

	mediaType := header.Header.Get("Content-Type")
	if mediaType == "" {
		respondWithError(w, http.StatusBadRequest, "Missing Content-Type for thumbnail", nil)
		return
	}
	mimeType, _, err := mime.ParseMediaType(mediaType)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to parse mime type", err)
		return
	}

	if mimeType != "image/jpeg" && mimeType != "image/png" {
		respondWithError(w, http.StatusBadRequest, "Invalid file type: only upload 'image/jpeg' and 'image/png'", err)
		return
	}

	fileName := getAssetName(videoID, mediaType)
	fileDiskPath := cfg.getAssertDiskPath(fileName)

	f, err := os.Create(fileDiskPath)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to create file.", err)
		return
	}
	defer f.Close()
	_, err = io.Copy(f, file)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to write to file.", err)
		return
	}
	thumbnailUrl := cfg.getAssertURL(fileDiskPath)
	video.ThumbnailURL = &thumbnailUrl

	err = cfg.DB.UpdateVideo(video)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update video", err)
		return
	}

	respondWithJSON(w, http.StatusOK, video)
}
