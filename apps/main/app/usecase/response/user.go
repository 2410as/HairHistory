package response

// ---

type CreateUser struct {
	Ent *userEnt `json:"ent"`
}

func NewCreateUser(userID string) *CreateUser {
	return &CreateUser{
		Ent: &userEnt{ID: userID},
	}
}

// ---

type userEnt struct {
	ID string `json:"id"`
}
