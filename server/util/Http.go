package util

import (
	"air-sync/util/logging"
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

type contextKey int

const requestLoggerKey contextKey = iota

type Response struct {
	StatusCode  int
	ContentType string
	Header      http.Header
	Body        []byte
	BodyStream  io.ReadCloser
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
			RequestLogger(req).Error(err)
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
			RequestLogger(req).Error(err)
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
		if res.Message == "" {
			if res.StatusCode >= 500 {
				res.Message = "Internal server error"
			} else if res.StatusCode >= 400 {
				res.Message = "Invalid request error"
			} else {
				res.Message = "OK"
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
	if res.Header != nil {
		for key, value := range res.Header {
			w.Header()[key] = value
		}
	} else if res.ContentType != "" {
		w.Header().Set("Content-Type", res.ContentType)
	}
	w.WriteHeader(res.StatusCode)
	if res.StatusCode == http.StatusNoContent {
		return
	}
	if res.Body != nil {
		if _, err := w.Write(res.Body); err != nil {
			RequestLogger(req).Error(err)
		}
	} else if res.BodyStream != nil {
		defer res.BodyStream.Close()
		if _, err := io.Copy(w, res.BodyStream); err != nil {
			RequestLogger(req).Error(err)
		}
	}
}

func GetClientIP(req *http.Request) string {
	host := req.Header.Get("X-Real-IP")
	if host == "" {
		host = req.Header.Get("X-Forwarded-For")
	}
	if host != "" {
		parts := strings.Split(host, ",")
		return parts[0]
	}
	host, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		return ""
	}
	return host
}
