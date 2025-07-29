package transaction

type Repository interface {
	Create(tx *Transaction) error
	GetByID(id int) (*Transaction, error)
	GetByUser(userID int) ([]*Transaction, error)
	Update(tx *Transaction) error
}
