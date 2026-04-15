package apperrors

type ErrEmailNotVerified struct {
	Email string
}

func (e *ErrEmailNotVerified) Error() string {
	return "Для входа в аккаунт необходимо подтвердить почту"
}
