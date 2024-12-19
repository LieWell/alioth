package auth

const (
	LackOfParamCode    = "100100"
	LackOfParamMessage = "lack of parameters!"

	JWTErrorCode    = "100401"
	JWTErrorMessage = "unauthorized!"

	DumplicateUserNameCode    = "100409"
	DumplicateUserNameMessage = "user already exists!"

	RegisterFailedCode    = "100102"
	RegisterFailedMessage = "register failed, please try again later!"

	UserNameOrPasswordIncorrectCode    = "100103"
	UserNameOrPasswordIncorrectMessage = "username or password incorrect!"
)
