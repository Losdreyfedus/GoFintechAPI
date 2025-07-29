package transaction

type Service interface {
	Credit(userID int, amount float64) error
	Debit(userID int, amount float64) error
	Transfer(fromUserID, toUserID int, amount float64) error
	Rollback(txID int) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Credit(userID int, amount float64) error {
	// TODO: Implement credit logic
	return nil
}

func (s *service) Debit(userID int, amount float64) error {
	// TODO: Implement debit logic
	return nil
}

func (s *service) Transfer(fromUserID, toUserID int, amount float64) error {
	// TODO: Implement transfer logic
	return nil
}

func (s *service) Rollback(txID int) error {
	// TODO: Implement rollback logic
	return nil
}
