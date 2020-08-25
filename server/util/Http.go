package util

import (
	"context"
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type contextKey int

const requestLoggerKey contextKey = iota

type Response struct {
	StatusCode  int
	ContentType string
	Body        []byte
}

type RestResponse struct {
	Status     string      `json:"status"`
	StatusCode int         `json:"-"`
	Message    string      `json:"message,omitempty"`
	Data       interface{} `json:"data,omitempty"`
	Error      string      `json:"error,omitempty"`
}

type RequestContext struct {
	Logger *log.Logger
	Vars   map[string]string
}

var (
	SuccessRestResponse = &RestResponse{
		Status:     "success",
		StatusCode: 200,
	}
	SuccessResponse = &Response{
		StatusCode:  200,
		ContentType: "text/plain",
		Body:        []byte("Success"),
	}
)

type RequestHandlerFunc func(req *http.Request) (*Response, error)
type RestHandlerFunc func(req *http.Request) (*RestResponse, error)

func WrapHandlerFunc(handler RequestHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		req = DecorateRequest(req)
		res, err := handler(req)
		if err != nil {
			res = &Response{
				StatusCode:  500,
				ContentType: "text/plain",
				Body:        []byte(err.Error()),
			}
		}
		if res == nil {
			res = SuccessResponse
		}
		if res.StatusCode <= 0 {
			res.StatusCode = SuccessResponse.StatusCode
		}
		if res.ContentType == "" {
			res.ContentType = SuccessResponse.ContentType
		}
		if res.Body == nil {
			res.Body = SuccessResponse.Body
		}
		WriteResponse(w, req, res)
	}
}

func WrapRestHandlerFunc(handler RestHandlerFunc) http.HandlerFunc {
	return WrapHandlerFunc(func(req *http.Request) (*Response, error) {
		res, err := handler(req)
		if err != nil {
			res = &RestResponse{
				StatusCode: 500,
				Error:      err.Error(),
			}
		}
		if res == nil {
			res = SuccessRestResponse
		}
		if res.StatusCode <= 0 {
			if res.Error != "" {
				res.StatusCode = 500
			} else {
				res.StatusCode = 200
			}
		}
		if res.Status == "" {
			if res.StatusCode >= 200 && res.StatusCode < 400 {
				res.Status = "success"
			} else {
				res.Status = "error"
			}
		}
		b, err := json.Marshal(res)
		if err != nil {
			return nil, err
		}
		return &Response{
			StatusCode:  res.StatusCode,
			ContentType: "application/json;charset=utf-8",
			Body:        b,
		}, nil
	})
}

func DecorateRequest(req *http.Request) *http.Request {
	logger := log.New()
	logger.Formatter = NewRequestLogFormatter(DefaultTextFormatter, req)
	ctx := context.WithValue(req.Context(), requestLoggerKey, logger)
	return req.WithContext(ctx)
}

func RequestLogger(req *http.Request) *log.Logger {
	if v := req.Context().Value(requestLoggerKey); v != nil {
		return v.(*log.Logger)
	}
	return nil
}

func CreateRestResponse(data interface{}) *RestResponse {
	return &RestResponse{
		Data: data,
	}
}

func WriteHttpError(w http.ResponseWriter, req *http.Request, err error) {
	WriteResponse(w, req, &Response{
		StatusCode:  500,
		ContentType: "text/plain",
		Body:        []byte(err.Error()),
	})
}

func WriteResponse(w http.ResponseWriter, req *http.Request, res *Response) {
	w.Header().Set("Content-Type", res.ContentType)
	w.WriteHeader(res.StatusCode)
	if _, err := w.Write(res.Body); err != nil {
		RequestLogger(req).Error(err)
	}
}
