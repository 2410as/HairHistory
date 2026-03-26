package response

type CreateUser struct {
	UserID string `json:"userId"`
}

func NewCreateUser(userID string) *CreateUser {
	return &CreateUser{UserID: userID}
}

