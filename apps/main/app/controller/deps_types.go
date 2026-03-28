package controller

import "github.com/annasakai/hairhistorymemo/apps/main/app/usecase"

// Deps bundles usecase instances for HTTP handlers.
type Deps struct {
	User        *usecase.User
	HairHistory *usecase.HairHistory
}
