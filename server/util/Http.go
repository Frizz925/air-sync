package util

import (
	"air-sync/util/logging"
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

type JsonResponse struct {
	StatusCode int
	Result     interface{}
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
		StatusCode: http.StatusOK,
		Status:     "success",
	}
	SuccessJsonResponse = &JsonResponse{
		StatusCode: http.StatusOK,
	}
	SuccessResponse = &Response{
		StatusCode: http.StatusOK,
	}
)

type RequestHandlerFunc func(req *http.Request) (*Response, error)
type JsonHandlerFunc func(req *http.Request) (*JsonResponse, error)
type RestHandlerFunc func(req *http.Request) (*RestResponse, error)

func WrapHandlerFunc(handler RequestHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		req = DecorateRequest(req)
		res, err := handler(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if res == nil {
			res = SuccessResponse
		}
		if res.StatusCode <= 0 {
			res.StatusCode = SuccessResponse.StatusCode
		}
		WriteResponse(w, req, res)
	}
}

func WrapJsonHandlerFunc(handler JsonHandlerFunc) http.HandlerFunc {
	return WrapHandlerFunc(func(req *http.Request) (*Response, error) {
		res, err := handler(req)
		if err != nil {
			return nil, err
		}
		if res == nil {
			res = SuccessJsonResponse
		}
		if res.StatusCode <= 0 {
			res.StatusCode = SuccessJsonResponse.StatusCode
		}
		body := []byte{}
		if res.Result != nil {
			b, err := json.Marshal(res.Result)
			if err != nil {
				return nil, err
			}
			body = b
		}
		return &Response{
			StatusCode:  res.StatusCode,
			ContentType: "application/json;charset=utf-8",
			Body:        body,
		}, nil
	})
}

func WrapRestHandlerFunc(handler RestHandlerFunc) http.HandlerFunc {
	return WrapJsonHandlerFunc(func(req *http.Request) (*JsonResponse, error) {
		res, err := handler(req)
		if err != nil {
			res = &RestResponse{
				StatusCode: http.StatusInternalServerError,
				Error:      err.Error(),
			}
		}
		if res == nil {
			res = SuccessRestResponse
		}
		if res.StatusCode <= 0 {
			if res.Error != "" {
				res.StatusCode = http.StatusInternalServerError
			} else {
				res.StatusCode = http.StatusOK
			}
		}
		if res.Status == "" {
			if res.StatusCode >= 200 && res.StatusCode < 400 {
				res.Status = "success"
			} else {
				res.Status = "error"
			}
		}
		return &JsonResponse{
			StatusCode: res.StatusCode,
			Result:     res,
		}, nil
	})
}

func DecorateRequest(req *http.Request) *http.Request {
	logger := CreateRequestLogger(req)
	ctx := context.WithValue(req.Context(), requestLoggerKey, logger)
	return req.WithContext(ctx)
}

func CreateRequestLogger(req *http.Request) *log.Logger {
	logger := log.New()
	logger.Formatter = logging.NewRequestLogFormatter(DefaultTextFormatter, req)
	return logger
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

func WriteResponse(w http.ResponseWriter, req *http.Request, res *Response) {
	if res.ContentType != "" {
		w.Header().Set("Content-Type", res.ContentType)
	}
	w.WriteHeader(res.StatusCode)
	if _, err := w.Write(res.Body); err != nil {
		switch err {
		case http.ErrBodyNotAllowed:
		default:
			RequestLogger(req).Error(err)
		}
	}
}
