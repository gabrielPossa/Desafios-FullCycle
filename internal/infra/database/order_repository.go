package database

import (
	"database/sql"

	"github.com/gabrielPossa/Desafios-FullCycle/cleanArch/internal/entity"
)

type OrderRepository struct {
	Db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{Db: db}
}

func (r *OrderRepository) Save(order *entity.Order) error {
	stmt, err := r.Db.Prepare("INSERT INTO orders (id, price, tax, final_price) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(order.ID, order.Price, order.Tax, order.FinalPrice)
	if err != nil {
		return err
	}
	return nil
}

func (r *OrderRepository) GetTotal() (int, error) {
	var total int
	err := r.Db.QueryRow("Select count(*) from orders").Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (r *OrderRepository) GetAllOrders() ([]entity.Order, error) {
	stmt, err := r.Db.Prepare("SELECT id, price, tax, final_price FROM orders")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []entity.Order
	for rows.Next() {
		var id string
		var price, tax, finalPrice float64
		if err := rows.Scan(&id, &price, &tax, &finalPrice); err != nil {
			return nil, err
		}

		orders = append(orders, entity.Order{
			ID:         id,
			Price:      price,
			Tax:        tax,
			FinalPrice: finalPrice,
		})
	}
	return orders, nil
}
