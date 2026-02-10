package types

type FormFieldValidation struct {
	Pattern   string `json:"pattern,optional"`   // 正则表达式
	MinLength int    `json:"minLength,optional"` // 最小长度
	MaxLength int    `json:"maxLength,optional"` // 最大长度
	Message   string `json:"message,optional"`   // 错误提示
}
