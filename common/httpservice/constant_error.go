package httpservice

const (
	ERR_BAD_REQUEST            = "bad request payload"
	ERR_CONFLICT               = "conflict error"
	ERR_FAILED_DEPENDENCY      = "failed dependency"
	ERR_FORBIDDEN              = "forbidden access"
	ERR_INVALID_CREDENTIALS    = "invalid credentials"
	ERR_INVALID_REQUEST        = "invalid request parameters"
	ERR_METHOD_NOT_ALLOWED     = "method not allowed"
	ERR_MISSING_HEADER_DATA    = "missing header data"
	ERR_NO_RESULT_DATA         = "no result data"
	ERR_NOT_FOUND              = "resource not found"
	ERR_PERMISSION_DENIED      = "permission denied"
	ERR_PRECONDITION_FAILED    = "precondition failed"
	ERR_RATE_LIMIT_EXCEEDED    = "rate limit exceeded"
	ERR_RESOURCE_EXISTS        = "resource already exists"
	ERR_RESOURCE_EXHAUSTED     = "resource exhausted"
	ERR_RESOURCE_LOCKED        = "resource locked"
	ERR_REQUEST_TIMEOUT        = "request timeout"
	ERR_TOO_MANY_REQUESTS      = "too many requests"
	ERR_UNAUTHORIZED           = "unauthorized access"
	ERR_UNPROCESSABLE_ENTITY   = "unprocessable entity"
	ERR_UNSUPPORTED_MEDIA_TYPE = "unsupported media type"
	ERR_BAD_GATEWAY            = "bad gateway"
	ERR_GATEWAY_TIMEOUT        = "gateway timeout"
	ERR_INTERNAL_SERVER_ERROR  = "internal server error"
	ERR_NOT_IMPLEMENTED        = "not implemented"
	ERR_SERVICE_TIMEOUT        = "service timeout"
	ERR_SERVICE_UNAVAILABLE    = "service unavailable"
	ERR_UNKNOWN_SOURCE         = "unknown error"
	ERR_ACCESS_TOKEN_INVALID   = "access token invalid"
	ERR_ACCESS_TOKEN_EXPIRED   = "access token expired"
)

var ErrorResponses = map[string]int{
	ERR_BAD_REQUEST:            400,
	ERR_CONFLICT:               409,
	ERR_FAILED_DEPENDENCY:      424,
	ERR_FORBIDDEN:              403,
	ERR_INVALID_CREDENTIALS:    401,
	ERR_INVALID_REQUEST:        400,
	ERR_METHOD_NOT_ALLOWED:     405,
	ERR_MISSING_HEADER_DATA:    400,
	ERR_NO_RESULT_DATA:         404,
	ERR_NOT_FOUND:              404,
	ERR_PERMISSION_DENIED:      403,
	ERR_PRECONDITION_FAILED:    412,
	ERR_RATE_LIMIT_EXCEEDED:    429,
	ERR_RESOURCE_EXISTS:        409,
	ERR_RESOURCE_EXHAUSTED:     413,
	ERR_RESOURCE_LOCKED:        423,
	ERR_REQUEST_TIMEOUT:        408,
	ERR_TOO_MANY_REQUESTS:      429,
	ERR_UNAUTHORIZED:           401,
	ERR_UNPROCESSABLE_ENTITY:   422,
	ERR_UNSUPPORTED_MEDIA_TYPE: 415,
	ERR_BAD_GATEWAY:            502,
	ERR_GATEWAY_TIMEOUT:        504,
	ERR_INTERNAL_SERVER_ERROR:  500,
	ERR_NOT_IMPLEMENTED:        501,
	ERR_SERVICE_TIMEOUT:        503,
	ERR_SERVICE_UNAVAILABLE:    503,
	ERR_UNKNOWN_SOURCE:         520,
	ERR_ACCESS_TOKEN_INVALID:   401,
	ERR_ACCESS_TOKEN_EXPIRED:   401,
}
