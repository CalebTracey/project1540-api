package routes

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"project1540-api/external/models"
	"project1540-api/external/models/s3"
	"project1540-api/graph"
	//"project1540-api/graph/generated"
	"project1540-api/internal/facade"
	"time"
)

type Handler struct {
	Resolver graph.Resolver
	Service  *facade.Service
}

// MiddlewareOption type alias for chi middleware options
type MiddlewareOption func(next http.Handler) http.Handler

func (h *Handler) InitializeRoutes(options ...MiddlewareOption) *chi.Mux {
	r := chi.NewRouter()

	for _, middlewareOption := range options {
		r.Use(middlewareOption)
	}

	h.Resolver = graph.Resolver{
		IFacade: h.Service.S3,
	}

	// graphql endpoints
	//r.Handle("/graphql", handler.NewDefaultServer(
	//	generated.NewExecutableSchema(
	//		generated.Config{Resolvers: &h.Resolver},
	//	),
	//))
	//r.Handle("/", playground.Handler(
	//	"GraphQL playground", "/graphql",
	//))

	// rest endpoints
	r.Post("/put", h.UploadS3Handler())
	r.Post("/get", h.DownloadS3Handler())

	return r
}

func (h *Handler) UploadS3Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		// access the file from http POST request
		if file, header, err := accessFile(r, w); err != nil {
			routeHandlerError(w, err.Error(), http.StatusBadRequest)
			return
		} else if tempFile, fileErr := createTempFile(w); fileErr != nil {
			routeHandlerError(w, fileErr.Error(), http.StatusInternalServerError)
			return
		} else if fileInfo, copyErr := copyFile(tempFile, file); copyErr == nil && fileInfo.IsDir() {
			openUploadFile(w, tempFile)
			// copy file, convert file types, and get file info
			if errLog := h.Service.S3.UploadS3Object(r.Context(),
				s3.NewUploadS3Request(s3.FromFile(header, fileInfo, tempFile, devBucket)),
			); errLog != nil {

				routeHandlerError(w, *errLog, errLog.StatusCode)
				return

			} else {
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
				log.Infof("Time taken: %.2fs", time.Since(start).Seconds())
			}
		} else if err != nil {
			routeHandlerError(w, err.Error(), http.StatusInternalServerError)
			log.Error(err)
		} else {
			log.Infof("file info: %v", fileInfo)
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
			return
		}
		if resp, errLog := h.Service.S3.DownloadS3Object(r.Context(), apiRequest); errLog != nil {
			routeHandlerError(w, *errLog, errLog.StatusCode)
			log.Error(*errLog)
			return
		} else {
			w.Header().Set("Content-Type", mime.TypeByExtension(filepath.Ext(*resp.ContentType)))
			// Write the file content to the response body
			if _, copyErr := io.Copy(w, resp.Body); copyErr != nil {
				routeHandlerError(w, copyErr.Error(), http.StatusInternalServerError)
				log.Errorf("failed to write file content to response body: %s", copyErr.Error())
				return
			}
			log.Infof("Time taken: %.2fs", time.Since(start).Seconds())
		}
	}
}

func accessFile(r *http.Request, w http.ResponseWriter) (multipart.File, *multipart.FileHeader, error) {
	if formFile, header, err := r.FormFile("file"); err != nil {
		routeHandlerError(w, err.Error(), http.StatusBadRequest)
		//log.Errorf("accessFile: form file error: %s", err.Error())
		return formFile, header, fmt.Errorf("accessFile: %w", err)
	} else {
		defer func(file multipart.File) {
			_ = file.Close()
		}(formFile)
		return formFile, header, err
	}
}

func createTempFile(w http.ResponseWriter) (*os.File, error) {
	if tempFile, err := os.CreateTemp("", "temp-*"); err != nil {
		routeHandlerError(w, err.Error(), http.StatusInternalServerError)
		log.Errorf("createTempFile: failed to create temporary file: %s", err.Error())
		return nil, err
	} else {
		defer func(name string) {
			_ = os.Remove(name)
		}(tempFile.Name())
		return tempFile, nil
	}
}

// copyFile copies the file data to the temporary file.
// on success, it returns the file info after opening the file for reading.
// needed for conversion of multipart.File to *os.File
func copyFile(tempFile *os.File, file multipart.File) (os.FileInfo, error) {
	if _, err := io.Copy(tempFile, file); err != nil {
		//routeHandlerError(w, err.Error(), http.StatusInternalServerError)
		//log.Errorf("copyFile: failed to copy temp file: %s", err.Error())
		return nil, fmt.Errorf("copyFile: %w", err)
	} else {
		fileInfo, _ := tempFile.Stat()
		//openUploadFile(w, tempFile)
		return fileInfo, nil
	}
}

// openUploadFile opens the file for reading
func openUploadFile(w http.ResponseWriter, uploadFile *os.File) {
	if _, err := uploadFile.Seek(0, io.SeekStart); err != nil {
		routeHandlerError(w, err.Error(), http.StatusInternalServerError)
		log.Errorf("openUploadFile: failed to open temp file: %s", err.Error())
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

const (
	devBucket = "project1540-dev"
)
