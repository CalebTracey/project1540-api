package routes

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"project1540-api/external/models"
	"project1540-api/external/models/s3"
	"project1540-api/internal/facade"
	"time"
)

type Handler struct {
	Service *facade.Service
}

// MiddlewareOption type alias for chi middleware options
type MiddlewareOption func(next http.Handler) http.Handler

func (h *Handler) InitializeRoutes(options ...MiddlewareOption) *chi.Mux {
	r := chi.NewRouter()

	for _, middlewareOption := range options {
		r.Use(middlewareOption)
	}

	// REST endpoints
	r.Post("/put", h.UploadS3Handler())
	r.Post("/get", h.DownloadS3Handler())

	return r
}

func (h *Handler) UploadS3Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()
		tempFile, _ := createTempFile()

		if file, header, accessErr := accessFile(r); accessErr == nil {
			// copy file, convert file types, and get file info
			if fileInfo, copyErr := copyFile(tempFile, file); copyErr == nil && fileInfo.IsDir() {
				_ = openUploadFile(tempFile)

				if errLog := h.Service.S3.UploadS3Object(
					r.Context(),
					s3.NewUploadS3Request(s3.FromFile(header, fileInfo, tempFile, devBucket)),
				); errLog == nil {

					writeResponse(w, start)
					log.Infof("Time taken: %.2fs", time.Since(start).Seconds())

				} else {
					log.Error(*errLog)
					routeHandlerError(w, *errLog, errLog.StatusCode)
				}
			} else if copyErr != nil {
				routeHandlerError(w, copyErr.Error(), http.StatusInternalServerError)
			} else {
				log.Infof("fileInfo: %v", fileInfo)
			}
		} else {
			routeHandlerError(w, accessErr.Error(), http.StatusBadRequest)
		}
	}
}

func (h *Handler) DownloadS3Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()
		var apiRequest s3.DownloadS3Request

		if err := json.NewDecoder(r.Body).Decode(&apiRequest); err != nil {
			routeHandlerError(w, err.Error(), http.StatusBadRequest)
			log.Error(err.Error())
		}

		if resp, errLog := h.Service.S3.DownloadS3Object(
			r.Context(), new(s3.DownloadS3Request).FromJSON(r),
		); errLog == nil {
			w.Header().Set("Content-Type", mime.TypeByExtension(filepath.Ext(*resp.ContentType)))
			// Write the file content to the response body
			if _, copyErr := io.Copy(w, resp.Body); copyErr != nil {
				routeHandlerError(w, copyErr.Error(), http.StatusInternalServerError)
				log.Errorf("failed to write file content to response body: %s", copyErr.Error())
			}
			log.Infof("Time taken: %.2fs", time.Since(start).Seconds())

		} else {

			routeHandlerError(w, *errLog, errLog.StatusCode)
			log.Error(*errLog)
		}
	}
}

func routeHandlerError[E string | models.ErrorLog](w http.ResponseWriter, err E, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	errorWrapper := struct {
		Error E `json:"error"`
	}{Error: err}

	if jsonErr := json.NewEncoder(w).Encode(&errorWrapper); jsonErr != nil {
		log.Error("routeHandlerError: failed encoding response to json: %s", jsonErr.Error())
	}
}

func writeResponse(w http.ResponseWriter, start time.Time) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	hostname, _ := os.Hostname()
	_ = json.NewEncoder(w).Encode(
		models.Message{
			Hostname: hostname,
			Time:     fmt.Sprintf("%.2fs", time.Since(start).Seconds()),
			Status:   http.StatusText(http.StatusOK),
		},
	)
}

const (
	devBucket = "project1540-dev"
)
