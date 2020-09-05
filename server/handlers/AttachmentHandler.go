package handlers

import (
	"air-sync/models"
	repos "air-sync/repositories"
	"air-sync/storages"
	"air-sync/util"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

const attachmentMaxSize int64 = 2 << 20

var ErrUploadFileTooLarge = errors.New("Uploaded file too large")

var ResAttachmentNotFound = &util.Response{
	StatusCode:  404,
	ContentType: "text/plain",
	Body:        []byte("Attachment not found"),
}

type AttachmentHandler struct {
	repo    repos.AttachmentRepository
	storage storages.Storage
}

func NewAttachmentHandler(repo repos.AttachmentRepository, storage storages.Storage) *AttachmentHandler {
	return &AttachmentHandler{repo, storage}
}

func (h *AttachmentHandler) RegisterRoutes(r *mux.Router) {
	s := r.PathPrefix("/attachments").Subrouter()
	s.HandleFunc("/upload", util.WrapRestHandlerFunc(h.UploadAttachment)).Methods("POST")
	s.HandleFunc("/{id}", util.WrapHandlerFunc(h.DownloadAttachment)).Methods("GET")
}

func (h *AttachmentHandler) UploadAttachment(req *http.Request) (*util.RestResponse, error) {
	if err := req.ParseMultipartForm(attachmentMaxSize); err != nil {
		return h.requestError(req, err)
	}
	file, header, err := req.FormFile("file")
	if err != nil {
		return h.requestError(req, err)
	}
	defer file.Close()
	if header.Size > attachmentMaxSize {
		return h.requestError(req, ErrUploadFileTooLarge)
	}

	buf := make([]byte, 512)
	n, err := io.ReadFull(file, buf)
	if err != nil {
		return h.requestError(req, err)
	}

	filename := header.Filename
	mime := http.DetectContentType(buf)
	typ := req.URL.Query().Get("type")
	attachment, err := h.repo.Create(models.NewCreateAttachment(filename, typ, mime))
	if err != nil {
		return nil, err
	}

	logger := util.RequestLogger(req)
	w, err := h.storage.Write(attachment.ID)
	if err != nil {
		return nil, err
	}
	defer w.Close()

	if _, err := w.Write(buf[:n]); err != nil {
		return nil, err
	}
	if _, err := io.Copy(w, file); err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}

	logger.WithField("attachment_id", attachment.ID).Info("Attachment uploaded")
	return &util.RestResponse{
		Message: "Attachment uploaded",
		Data:    attachment,
	}, nil
}

func (h *AttachmentHandler) DownloadAttachment(req *http.Request) (*util.Response, error) {
	id := mux.Vars(req)["id"]
	attachment, err := h.repo.Find(id)
	if err != nil {
		if errors.Is(err, repos.ErrAttachmentNotFound) {
			return ResAttachmentNotFound, nil
		}
		return nil, err
	}
	exists, err := h.storage.Exists(attachment.ID)
	if err != nil {
		return nil, err
	} else if !exists {
		return ResAttachmentNotFound, nil
	}
	r, err := h.storage.Read(attachment.ID)
	if err != nil {
		return nil, err
	}
	header := make(http.Header)
	header.Set("Content-Type", attachment.Mime)
	if attachment.Type == "file" {
		header.Set(
			"Content-Disposition",
			fmt.Sprintf("attachment; filename=\"%s\"", attachment.Name),
		)
	}
	return &util.Response{
		Header:     header,
		BodyStream: r,
	}, nil
}

func (h *AttachmentHandler) requestError(req *http.Request, err error) (*util.RestResponse, error) {
	util.RequestLogger(req).Error(err)
	return &util.RestResponse{
		StatusCode: 400,
		Message:    "Request malformed",
		Error:      err.Error(),
	}, nil
}
