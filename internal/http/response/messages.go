package response

const (
	CodeInvalidRequest      = "INVALID_REQUEST"
	CodeUnauthorized        = "UNAUTHORIZED"
	CodeForbidden           = "FORBIDDEN"
	CodeNotFound            = "NOT_FOUND"
	CodeTooManyRequests     = "TOO_MANY_REQUESTS"
	CodeRegisterFailed      = "REGISTER_FAILED"
	CodeInvalidCredentials  = "INVALID_CREDENTIALS"
	CodeInvalidRefreshToken = "INVALID_REFRESH_TOKEN"
	CodeConversationFailed  = "CONVERSATION_OPERATION_FAILED"
	CodeMessageFailed       = "MESSAGE_OPERATION_FAILED"
	CodeInternal            = "INTERNAL_ERROR"
)

const (
	MsgInvalidRequestBody   = "Invalid request payload."
	MsgUnauthorized         = "Authentication is required."
	MsgForbidden            = "You do not have access to this resource."
	MsgUserNotFound         = "User not found."
	MsgInvalidConversation  = "Conversation ID must be a positive integer."
	MsgConversationRequired = "conversationId query parameter is required."
	MsgTooManyRequests      = "Too many requests. Please try again later."
	MsgRegisterFailed       = "Failed to register user."
	MsgInvalidCredentials   = "Invalid email or password."
	MsgInvalidRefreshToken  = "Invalid refresh token."
	MsgCreateConversation   = "Failed to create conversation."
	MsgListConversations    = "Failed to list conversations."
	MsgSendMessage          = "Failed to send message."
	MsgGetMessages          = "Failed to get messages."
	MsgInternalServer       = "Internal server error."
)

const (
	MsgRegistered = "Registered successfully."
	MsgOK         = "OK"
)
