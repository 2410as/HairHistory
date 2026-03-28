package controller

type Health struct {
	db DBPing
}

func NewHealth(db DBPing) *Health {
	return &Health{db: db}
}
