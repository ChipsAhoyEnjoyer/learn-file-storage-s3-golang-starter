package main

import (
	"bytes"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"

	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUploadThumbnail(w http.ResponseWriter, r *http.Request) {
	videoIDString := r.PathValue("videoID")
	videoID, err := uuid.Parse(videoIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
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
	data, err := io.ReadAll(file)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to read image", err)
		return
	}

	mediaTypeStr := header.Header.Get("Content-Type")
	if mediaTypeStr == "" {
		respondWithError(w, http.StatusBadRequest, "Missing Content-Type for thumbnail", nil)
		return
	}

	mediaType, _, err := mime.ParseMediaType(mediaTypeStr)
	if (err != nil) || (mediaType != "image/png" && mediaType != "image/jpeg") {
		respondWithError(w, http.StatusBadRequest, "Malformed/incorrect media type", err)
		return
	}

	video, err := cfg.db.GetVideo(videoID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to fetch video data", err)
		return
	}

	if video.UserID != userID {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized user", nil)
		return
	}

	ext, err := getAssetExtension(mediaType)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Incorrect media type", err)
		return
	}
	filename := getAssetName(videoID, ext)
	url := getAssetURL(
		filename,
		cfg.port,
	)

	if err = createThumbnailFile(cfg.assetsRoot, filename, data); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error uploading thumbnail", err)
		return
	}
	video.ThumbnailURL = &url

	if err = cfg.db.UpdateVideo(video); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to update video", err)
		return
	}

	respondWithJSON(w, http.StatusOK, video)
}

func createThumbnailFile(assetsRoot, filename string, data []byte) error {
	path := getAssetPath(filename, assetsRoot)
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	thumbnail := bytes.NewReader(data)
	if _, err := io.Copy(f, thumbnail); err != nil {
		return err
	}
	return nil
}
