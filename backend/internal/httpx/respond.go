package httpx

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

const maxRequestBodyBytes = 1 << 20

type errorBody struct {
	Code    ErrorCode         `json:"code"`
	Message string            `json:"message"`
	Details map[string]string `json:"details,omitempty"`
}

type errorEnvelope struct {
	Error errorBody `json:"error"`
}

func JSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if payload == nil {
		return
	}
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		slog.Error("encode response body failed", "error", err)
	}
}

func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

func Error(w http.ResponseWriter, r *http.Request, err error) {
	apiErr := AsAPIError(err)

	logger := slog.With("request_id", RequestIDFromContext(r.Context()), "method", r.Method, "path", r.URL.Path, "code", string(apiErr.Code))
	if apiErr.Code == CodeInternal {
		logger.Error("request failed", "error", apiErr.Error())
	} else {
		logger.Info("request rejected", "message", apiErr.Message)
	}

	JSON(w, apiErr.HTTPStatus(), errorEnvelope{Error: errorBody{
		Code:    apiErr.Code,
		Message: apiErr.Message,
		Details: apiErr.Details,
	}})
}

func DecodeJSON(r *http.Request, dst any) error {
	decoder := json.NewDecoder(io.LimitReader(r.Body, maxRequestBodyBytes))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(dst); err != nil {
		if errors.Is(err, io.EOF) {
			return InvalidArgument("request body is required", nil)
		}
		var typeErr *json.UnmarshalTypeError
		if errors.As(err, &typeErr) {
			return InvalidArgument(fmt.Sprintf("field %q has an invalid type", typeErr.Field), nil)
		}
		return InvalidArgument("request body is not valid JSON", nil)
	}

	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return InvalidArgument("request body must contain a single JSON object", nil)
	}
	return nil
}
