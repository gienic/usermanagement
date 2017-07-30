package server

type (
	// EncryptData : general struct for end-to-end encryption
	EncryptData struct {
		Public string
	}

	// AuthData : general authentication data struct
	AuthData struct {
		Token string
	}

	// SignupData : general signup data struct
	SignupData struct {
		Firstname string
		Lastname  string
		Email     string
		SlackData
		TelegramData
		MessengerData
	}

	// LoginData : general login data struct
	LoginData struct {
		User string
		Pass string
	}

	// SlackData : contains slack specific user data
	SlackData struct{}
	// TelegramData : contains telegram specific user data
	TelegramData struct{}
	// MessengerData : contains messenger specific user data
	MessengerData struct{}
)
