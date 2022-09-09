package consts

const (
	ErrInvalidEmailOrPassword  = "invalid_email_or_password"
	ErrInvalidUserCode         = "invalid_user_code"
	ErrUserAlreadyExist        = "user_already_exist"
	ErrUserNotFound            = "user_not_found"
	ErrInternalServer          = "internal_transport_error"
	ErrUserNotActive           = "user_is_not_active"
	ErrInvalidToken            = "token_is_invalid"
	ErrAlreadyConfirmed        = "user_already_confirmed"
	ErrUnauthorized            = "unauthorized"
	ErrRecoveryRequestNotFound = "recovery_request_not_found"

	ErrFieldRequired        = "must_be_required"
	ErrFieldIncorrectFormat = "incorrect_format"
)
