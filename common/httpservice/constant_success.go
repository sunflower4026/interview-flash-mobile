package httpservice

const (
	SUCCESS_OK               = "OK"
	SUCCESS_CREATED          = "Resource created successfully"
	SUCCESS_ACCEPTED         = "Request accepted"
	SUCCESS_NO_CONTENT       = "No content to return"
	SUCCESS_RESET_CONTENT    = "Content reset successfully"
	SUCCESS_PARTIAL_CONTENT  = "Partial content delivered"
	SUCCESS_MULTI_STATUS     = "Multiple statuses available"
	SUCCESS_ALREADY_REPORTED = "Already reported"
	SUCCESS_IM_USED          = "Instance manipulations used"
)

var SuccessResponses = map[string]int{
	SUCCESS_OK:               200,
	SUCCESS_CREATED:          201,
	SUCCESS_ACCEPTED:         202,
	SUCCESS_NO_CONTENT:       204,
	SUCCESS_RESET_CONTENT:    205,
	SUCCESS_PARTIAL_CONTENT:  206,
	SUCCESS_MULTI_STATUS:     207,
	SUCCESS_ALREADY_REPORTED: 208,
	SUCCESS_IM_USED:          226,
}
