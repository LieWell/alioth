package auth

const (
	JWTErrorCode    = "100401"
	JWTErrorMessage = "unauthorized!"

	DumplicateUserNameCode    = "100409"
	DumplicateUserNameMessage = "user already exists!"

	LackOfParamCode    = "100100"
	LackOfParamMessage = "lack of parameters!"

	RegisterDisabledCode    = "100102"
	RegisterDisabledMessage = "register not enabled!"

	RegisterFailedCode    = "100103"
	RegisterFailedMessage = "register failed, please try again later!"

	UserNameOrPasswordIncorrectCode    = "100104"
	UserNameOrPasswordIncorrectMessage = "username or password incorrect!"
)
