package models

import "errors"

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

type VerifyInput struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

func (v VerifyInput) Validate() error {
	if v.Code == "" {
		return errors.New("Код не может быть пустым")
	}
	if v.Email == "" {
		return errors.New("Почта не может быть пустой")
	}

	return nil
}

type GenerateNewCodeInput struct {
	Email string `json:"email"`
}

func (v GenerateNewCodeInput) Validate() error {
	if v.Email == "" {
		return errors.New("Почта не может быть пустой")
	}

	return nil
}
