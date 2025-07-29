package balance

type Service interface {
	UpdateBalance(userID int, amount float64) error
	GetCurrentBalance(userID int) (float64, error)
	GetHistoricalBalance(userID int, atTime string) (float64, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) UpdateBalance(userID int, amount float64) error {
	// TODO: Implement thread-safe balance update
	return nil
}

func (s *service) GetCurrentBalance(userID int) (float64, error) {
	// TODO: Implement current balance retrieval
	return 0, nil
}

func (s *service) GetHistoricalBalance(userID int, atTime string) (float64, error) {
	// TODO: Implement historical balance retrieval
	return 0, nil
}
