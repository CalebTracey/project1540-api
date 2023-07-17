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
	"project1540-api/external/models/postgres"
	"project1540-api/external/models/s3"
	"project1540-api/internal/facade"
	"strconv"
	"time"
)

type Handler struct {
	Service *facade.Service
}

// MiddlewareOption type alias for chi middleware options
type MiddlewareOption func(next http.Handler) http.Handler

func (h *Handler) InitializeRoutes(options ...MiddlewareOption) *chi.Mux {
	r := chi.NewRouter()

	// REST endpoints
	r.Route("/api", func(r chi.Router) {
		for _, middlewareOption := range options {
			r.Use(middlewareOption)
		}

		r.Post("/put", h.UploadS3Handler())
		r.Post("/get", h.DownloadS3Handler())
		r.Post("/newFile", h.InsertNewFileHandler())
		r.Post("/update", h.UpdateDatabaseWithS3Data())
		r.Post("/search", h.SearchFilesByTagHandler())
	})

	return r
}

func (h *Handler) SearchFilesByTagHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		fileRequest := new(postgres.SearchFilesByTagRequest)
		if err := fileRequest.FromJSON(r); err != nil {
			log.Errorf("SearchFilesByTagHandler: %v", err)
			routeHandlerError(w, err.Error(), http.StatusBadRequest)
		}

		if response := h.Service.PostgresQL.SearchFilesByTag(
			r.Context(), fileRequest.Tags,
		); response.Message.ErrorLogs == nil {

			hostname, _ := os.Hostname()
			response.Message.Hostname = hostname
			response.Message.Time = time.Since(start).String()
			response.Message.Status = http.StatusText(http.StatusOK)
			response.Message.Count = strconv.Itoa(len(response.Files))

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			if err := response.ToJSON(w); err != nil {
				//writeResponseFromResults(w, start, response)
				log.Errorf("SearchFilesByTagHandler: %v", err)
				routeHandlerError(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			log.Errorf("SearchFilesByTagHandler: %v", response.Message)
			routeHandlerError(w, response.Message.ErrorLogs[0], http.StatusInternalServerError)
		}
	}
}

func (h *Handler) UpdateDatabaseWithS3Data() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		if serviceErr := h.Service.UpdateDatabaseFromS3Bucket(
			r.Context(), devBucket,
		); serviceErr == nil {
			writeResponse(w, start)
		} else {
			log.Errorf("UpdateDatabaseWithS3Data: %v", *serviceErr)
			routeHandlerError(w, *serviceErr, http.StatusInternalServerError)
		}
	}
}

func (h *Handler) InsertNewFileHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		var fileRequest *postgres.NewFileRequest
		if err := fileRequest.FromJSON(r); err != nil {
			log.Errorf("InsertNewFileHandler: %v", err)
			routeHandlerError(w, err.Error(), http.StatusBadRequest)
		}

		if serviceErr := h.Service.InsertNewFileByS3Bucket(
			r.Context(), *fileRequest,
		); serviceErr == nil {
			writeResponse(w, start)
		} else {
			log.Errorf("InsertNewFileHandler: %v", *serviceErr)
			routeHandlerError(w, *serviceErr, http.StatusInternalServerError)
		}
	}
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
