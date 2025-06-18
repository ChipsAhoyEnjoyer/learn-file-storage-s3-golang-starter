package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/auth"
	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/video"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUploadVideo(w http.ResponseWriter, r *http.Request) {
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
	videoID, err := uuid.Parse(r.PathValue("videoID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID", err)
		return
	}

	const maxMemory = 1 << 30 // 1 GB limit
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxMemory))

	videoDetails, err := cfg.db.GetVideo(videoID)
	if err != nil {
		errCode := http.StatusInternalServerError
		if err == sql.ErrNoRows {
			errCode = http.StatusNotFound
		}
		respondWithError(w, errCode, "Unable to retrieve video", err)
		return
	}
	if videoDetails.UserID != userID {
		respondWithError(w, http.StatusUnauthorized, "User not authorized to upload to this video", nil)
		return
	}

	file, header, err := r.FormFile("video")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to upload video", err)
		return
	}
	defer file.Close()

	mediaType, _, err := mime.ParseMediaType(header.Header.Get("Content-Type"))
	if err != nil || mediaType != "video/mp4" {
		respondWithError(w, http.StatusBadRequest, "Incorrect content type", err)
		return
	}

	temp, err := os.CreateTemp("", "temp.mp4")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to upload video to memory", err)
		return
	}
	defer os.Remove(temp.Name())
	defer temp.Close()

	_, err = io.Copy(temp, file)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to copy video to temporary file", err)
		return
	}

	_, err = temp.Seek(0, io.SeekStart)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to read uploaded video file", err)
		return
	}

	ext, err := getAssetExtension(mediaType)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Incorrect media type", err)
		return
	}

	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error generating random bytes", err)
		return
	}
	ratio, err := video.GetVideoAspectRatio(temp.Name())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to get video aspect ratio", err)
		return
	}
	fileKeyPrefix := ""
	switch ratio {
	case "16:9":
		fileKeyPrefix = "landscape"
	case "9:16":
		fileKeyPrefix = "portrait"
	default:
		fileKeyPrefix = "other"
	}

	processedVideoFilePath, err := video.ProcessVideoForFastStart(temp.Name())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to convert video to fast start", err)
		return
	}

	fileKey := fileKeyPrefix + "/" + base64.RawURLEncoding.EncodeToString(b) + "." + ext

	_, err = cfg.s3Client.PutObject(
		r.Context(),
		&s3.PutObjectInput{
			Bucket:      aws.String(cfg.s3Bucket),
			Key:         aws.String(fileKey),
			Body:        temp,
			ContentType: &mediaType,
		},
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to upload video to storage", err)
		return
	}
	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", cfg.s3Bucket, cfg.s3Region, fileKey)
	videoDetails.VideoURL = &url

	if err = cfg.db.UpdateVideo(videoDetails); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to update the database's video url", err)
		return
	}
	respondWithJSON(
		w,
		http.StatusOK,
		map[string]string{"video_url": url})
}
