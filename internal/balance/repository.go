package balance

type Repository interface {
	GetByUserID(userID int) (*Balance, error)
	Update(balance *Balance) error
}
