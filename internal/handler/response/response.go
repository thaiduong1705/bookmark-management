// Package response holds the response templates shared by all handlers.
//
// Endpoint-specific responses embed these templates and add their own fields:
//
//	type CreateResponse struct {
//	    response.Message
//	    Code string `json:"code"`
//	}
//
// Embedded structs are flattened in JSON and understood by swag.
package response

// Message is the common part of every response body, success or error.
type Message struct {
	Message string `json:"message"`
}

// NewMessage builds a Message.
func NewMessage(message string) Message {
	return Message{Message: message}
}
