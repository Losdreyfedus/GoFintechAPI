package balance

type Service interface {
	UpdateBalance(userID int, amount float64) error
	GetCurrentBalance(userID int) (float64, error)
	GetHistoricalBalance(userID int, atTime string) (float64, error)
}
