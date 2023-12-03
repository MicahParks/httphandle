package httphandle

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	jt "github.com/MicahParks/jsontype"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	hhconst "github.com/MicahParks/httphandle/constant"
	"github.com/MicahParks/httphandle/middleware/ctxkey"
)

type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewAPIError(ctx context.Context, code int, message string) APIResponse[APIError] {
	apiError := APIError{
		Code:    code,
		Message: message,
	}
	meta := APIMetadata{
		RequestUUID: ctx.Value(ctxkey.ReqUUID).(uuid.UUID),
	}
	return APIResponse[APIError]{
		Data:     apiError,
		Metadata: meta,
	}
}

type APIMetadata struct {
	RequestUUID uuid.UUID `json:"requestMetadata"`
}

type APIResponse[Data any] struct {
	Data     Data        `json:"data,omitempty"`
	Metadata APIMetadata `json:"metadata"`
}

func APICommitTx(ctx context.Context, responseCode int) (code int, body []byte, err error) {
	tx := ctx.Value(ctxkey.Tx).(pgx.Tx)
	err = tx.Commit(ctx)
	if err != nil {
		l := ctx.Value(ctxkey.Logger).(*slog.Logger)
		l.ErrorContext(ctx, "Failed to commit transaction.",
			hhconst.LogErr, err,
		)
		return APIErrorResponse(ctx, http.StatusInternalServerError, hhconst.RespInternalServerError)
	}
	return APIJSON(ctx, responseCode, APIResponse[any]{})
}

func APIErrorResponse(ctx context.Context, code int, message string) (int, []byte, error) {
	data, err := errorBody(ctx, code, message)
	if err != nil {
		return 0, nil, err
	}
	return code, data, nil
}

func APIJSONBody[ReqData jt.Defaulter[ReqData]](r *http.Request) (reqData ReqData, ctx context.Context, code int, body []byte, err error) {
	//goland:noinspection GoUnhandledErrorResult
	defer r.Body.Close()

	b, err := io.ReadAll(r.Body)
	if err != nil {
		code, body, _ = APIErrorResponse(ctx, http.StatusBadRequest, "Failed to read request body.")
		return reqData, ctx, code, body, err
	}

	err = json.Unmarshal(b, &reqData)
	if err != nil {
		code, body, _ = APIErrorResponse(ctx, http.StatusUnsupportedMediaType, "Failed to JSON parse request body.")
		return reqData, ctx, code, body, err
	}

	reqData, err = reqData.DefaultsAndValidate()
	if err != nil {
		code, body, _ = APIErrorResponse(ctx, http.StatusUnprocessableEntity, "Failed to validate request body.")
		return reqData, ctx, code, body, err
	}

	return reqData, ctx, http.StatusOK, nil, nil
}

func APIJSON(ctx context.Context, code int, r APIResponse[any]) (int, []byte, error) {
	meta := APIMetadata{
		RequestUUID: ctx.Value(ctxkey.ReqUUID).(uuid.UUID),
	}
	r.Metadata = meta
	data, err := json.Marshal(r)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to JSON marshal response: %w", err)
	}
	return code, data, nil
}

func errorBody(ctx context.Context, code int, message string) ([]byte, error) {
	data, err := json.Marshal(NewAPIError(ctx, code, message))
	if err != nil {
		return nil, fmt.Errorf("failed to JSON marshal error response: %w", err)
	}
	return data, nil
}
