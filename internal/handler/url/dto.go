package url

import "bookmark-management/internal/handler/response"

// CreateRequest is the body of a create-short-URL request; Expr is in minutes (0 = never expires).
type CreateRequest struct {
	Url string `json:"url" binding:"required,url"`
	Exp int    `json:"exp" binding:"gte=0,lte=84600"`
}

// CreateResponse is the body returned after a short URL is created.
type CreateResponse struct {
	response.Message
	Code string `json:"code"`
}

// CodeParam is the path parameter of a redirect request.
type CodeParam struct {
	Code string `uri:"code" binding:"required"`
}
