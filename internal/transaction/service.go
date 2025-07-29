package transaction

type Service interface {
	Credit(userID int, amount float64) error
	Debit(userID int, amount float64) error
	Transfer(fromUserID, toUserID int, amount float64) error
	Rollback(txID int) error
}
