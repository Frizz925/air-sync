package handlers

import (
	"air-sync/models"
	repos "air-sync/repositories"
	"air-sync/storages"
	"air-sync/util"
	"errors"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

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
	req.ParseMultipartForm(2 << 20)
	file, header, err := req.FormFile("file")
	if err != nil {
		return h.requestError(req, err)
	}
	defer file.Close()

	buf := make([]byte, 512)
	n, err := io.ReadFull(file, buf)
	if err != nil {
		return h.requestError(req, err)
	}

	filename := header.Filename
	mime := http.DetectContentType(buf)
	atype := req.URL.Query().Get("type")
	attachment, err := h.repo.Create(models.NewCreateAttachment(filename, atype, mime))
	if err != nil {
		return nil, err
	}

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
	return &util.Response{
		Header: http.Header{
			"Content-Type":        []string{attachment.Mime},
			"Content-Disposition": []string{"attachment", "filename=" + attachment.Filename},
		},
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
