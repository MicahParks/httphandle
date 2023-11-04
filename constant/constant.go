// Package constant contains constants for the application.
package constant

const (
	// HeaderAcceptEncoding is the header key for the accepted encodings.
	HeaderAcceptEncoding = "Accept-Encoding"
	// HeaderCacheControl is the header key for the cache control.
	HeaderCacheControl = "Cache-Control"
	// HeaderContentEncoding is the header key for the content encoding.
	HeaderContentEncoding = "Content-Encoding"
	// ContentEncodingGzip is the content encoding for gzip.
	ContentEncodingGzip = "gzip"
	// HeaderContentType is the header key for the content type.
	HeaderContentType = "Content-Type"
	// ContentTypeForm is the content type for form data.
	ContentTypeForm = "application/x-www-form-urlencoded"
	// ContentTypeJSON is the content type for JSON data.
	ContentTypeJSON = "application/json"
	// MsgFailTransactionBegin is the log message for a failed transaction start.
	MsgFailTransactionBegin = "Failed to begin transaction."
	// MsgFailTransactionCommit is the log message for a failed transaction commit.
	MsgFailTransactionCommit = "Failed to commit transaction."
	// MsgFailTransactionRollback is the log message for a failed transaction rollback.
	MsgFailTransactionRollback = "Failed to rollback transaction."
	// LogFmt is the format for logging with the built-in logger.
	LogFmt = "%s\nError: %v"
	// LogErr is the key for the error in slog fields.
	LogErr = "error"
	// LogRespCode is the key for the response code in slog fields.
	LogRespCode = "respCode"
	// PathIndex is the path for the index page.
	PathIndex = "/"
	// RespInternalServerError is the response message for an internal server error.
	RespInternalServerError = "Internal server error."
	// StaticDir is the directory for static files.
	StaticDir = "static"
	// TemplateHeaderAddExtension is the extension for extra HTML to add to the header. files.
	TemplateHeaderAddExtension = ".header"
)
