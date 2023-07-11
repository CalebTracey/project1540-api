package routes

import (
	"encoding/json"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"project1540-api/external/models"
	"project1540-api/graph"
	"project1540-api/graph/generated"
	"project1540-api/internal/facade"
	"strconv"
)

type Handler struct {
	Resolver graph.Resolver
	Service  facade.IFacade
}

type MiddlewareOption func(next http.Handler) http.Handler

func (h *Handler) InitializeRoutes(options ...MiddlewareOption) *chi.Mux {
	r := chi.NewRouter()

	for _, middlewareOption := range options {
		r.Use(middlewareOption)
	}

	h.Resolver.IFacade = h.Service

	// graphql endpoints
	r.Handle("/graphql", handler.NewDefaultServer(
		generated.NewExecutableSchema(
			generated.Config{Resolvers: &h.Resolver},
		),
	))
	r.Handle("/", playground.Handler(
		"GraphQL playground", "/graphql",
	))

	// rest endpoints
	r.Post("/put", h.UploadS3Handler())

	return r
}

func (h *Handler) UploadS3Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var file multipart.File
		var tempFile *os.File
		var header *multipart.FileHeader
		var err error

		//Access the file key
		if file, header, err = r.FormFile("file"); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Errorf("UploadS3Handler: form file error: %s", err.Error())
			return
		}
		defer func(file multipart.File) {
			if err = file.Close(); err != nil {
				log.Errorf("UploadS3Handler: failed to closed file: %s", err.Error())
			}
		}(file)

		if tempFile, err = os.CreateTemp("", "temp-*"); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Errorf("UploadS3Handler: failed to create temp file: %s", err.Error())
			return
		}
		defer func(name string) {
			if err := os.Remove(name); err != nil {
				log.Error(err)
			}
		}(tempFile.Name())

		log.Printf("UploadS3Handler: file header name: %v", header.Filename)

		// Copy the file data to the temporary file
		if _, err = io.Copy(tempFile, file); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Errorf("UploadS3Handler: failed to copy temp file: %s", err.Error())
			return
		}

		// Get information about the uploaded file
		fileInfo, _ := tempFile.Stat()
		fileSize := fileInfo.Size()
		fileName := header.Filename

		// Open the temporary file for reading
		if _, err = tempFile.Seek(0, io.SeekStart); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Errorf("UploadS3Handler: failed to open temp file: %s", err.Error())
			return
		}

		// Create an input model for uploading to S3
		input := models.InputFile{
			Name: fileName,
			Size: strconv.Itoa(int(fileSize)),
			Type: header.Header.Get("Content-Type"),
			File: tempFile,
		}

		if errLog := h.Service.UploadS3(r.Context(), input); errLog != nil {
			if _, jsonErr := json.Marshal(&w); jsonErr != nil {
				log.Error("UploadS3Handler: %s", jsonErr.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else {
			log.Infoln("s3 upload successful!")
			w.WriteHeader(http.StatusOK)
		}
	}
}

const devAwsS3Bucket = "project1540-dev"
