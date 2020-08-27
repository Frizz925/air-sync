package handlers

import (
	"air-sync/util"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	qrcode "github.com/skip2/go-qrcode"
)

type QrRestHandler int

var _ RouteHandler = (*QrRestHandler)(nil)

func (h QrRestHandler) RegisterRoutes(r *mux.Router) {
	s := r.PathPrefix("/qr").Subrouter()
	s.HandleFunc("/generate", util.WrapHandlerFunc(h.GenerateQR)).Methods("POST")
}

func (QrRestHandler) GenerateQR(req *http.Request) (*util.Response, error) {
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	q, err := qrcode.Encode(string(b), qrcode.Medium, 256)
	if err != nil {
		return nil, err
	}
	return &util.Response{
		ContentType: "image/png",
		Body:        q,
	}, nil
}
