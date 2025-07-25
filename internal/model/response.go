package model

// Response 统一响应格式
type Response struct {
	Status    string      `json:"status"`               // 状态：success/error
	Code      string      `json:"code,omitempty"`       // 错误代码
	Message   string      `json:"message"`              // 响应消息
	Data      interface{} `json:"data,omitempty"`       // 响应数据
	RequestID string      `json:"request_id,omitempty"` // 请求ID
}

// 响应状态常量
const (
	StatusSuccess = "success"
	StatusError   = "error"
)

// 错误代码常量
const (
	// 通用错误
	ErrInvalidParam   = "INVALID_PARAM"   // 无效的参数
	ErrInvalidRequest = "INVALID_REQUEST" // 无效的请求
	ErrInternal       = "INTERNAL_ERROR"  // 内部错误
	ErrUnauthorized   = "UNAUTHORIZED"    // 未授权

	// 业务错误
	// 仅为错误码常量，无敏感信息
	ErrInvalidAPIKey    = "INVALID_API_KEY"    // 无效的API密钥
	ErrTemplateNotFound = "TEMPLATE_NOT_FOUND" // 模板不存在
	ErrGroupNotFound    = "GROUP_NOT_FOUND"    // 通知组不存在
	ErrUserNotFound     = "USER_NOT_FOUND"     // 用户不存在
	ErrSendFailed       = "SEND_FAILED"        // 发送失败
)

// NewSuccessResponse 创建成功响应
func NewSuccessResponse(message string, data interface{}, requestID string) Response {
	return Response{
		Status:    StatusSuccess,
		Message:   message,
		Data:      data,
		RequestID: requestID,
	}
}

// NewErrorResponse 创建错误响应
func NewErrorResponse(code, message string, requestID string) Response {
	return Response{
		Status:    StatusError,
		Code:      code,
		Message:   message,
		RequestID: requestID,
	}
}

// SendResult 发送结果
type SendResult struct {
	SuccessCount int      `json:"success_count"` // 成功数量
	FailedCount  int      `json:"failed_count"`  // 失败数量
	SuccessUsers []string `json:"success_users"` // 成功用户列表
	FailedUsers  []string `json:"failed_users"`  // 失败用户列表
}
