package models

type EmailTask struct {
	Email string `json:"email"`
	Code  string `json:"code"`
	Type  string `json:"type"`
}

func NewEmailTask(email string, code string, typeTask string) *EmailTask {
	return &EmailTask{
		Email: email,
		Code:  code,
		Type:  typeTask,
	}
}
