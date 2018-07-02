package services

import "github.com/ZhenlyChen/BugServer/httpServer/models"

// GameService ...
type GameService interface {
	GetNewestVersion() models.Game
	SetNewVersion(data models.Game) error
}

type gameService struct {
	Model   *models.GameModel
	Service *Service
}

func (s *gameService) GetNewestVersion() models.Game {
	return s.Model.GetNewestVersion()
}

func (s *gameService) SetNewVersion(data models.Game) error {
	return s.Model.SetNewVersion(data)
}
