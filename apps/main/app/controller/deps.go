package controller

import "github.com/annasakai/hairhistorymemo/apps/main/app/usecase"

type Deps struct {
	CreateUser     *usecase.CreateUserUsecase
	ListHistories  *usecase.ListHistoriesUsecase
	CreateHistory  *usecase.CreateHistoryUsecase
	UpdateHistory  *usecase.UpdateHistoryUsecase
	DeleteHistory  *usecase.DeleteHistoryUsecase
}

