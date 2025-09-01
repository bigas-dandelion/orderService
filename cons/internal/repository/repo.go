package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"l0/cons/internal/cache"
	"l0/cons/internal/models"
	"log"
)

type Repository struct {
	db    *sql.DB
	cache *cache.Cache
}

func NewRepository(db *sql.DB, cache *cache.Cache) *Repository {
	repo := &Repository{
		db:    db,
		cache: cache,
	}

	repo.loadToCache()
	return repo
}

func (r *Repository) StoreOrder(order *models.Order) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		INSERT INTO orders (order_uid, track_number, entry, locale, 
		internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`, order.OrderUID, order.TrackNumber, order.Entry, order.Locale,
		order.InternalSignature, order.CustomerID, order.DeliveryService,
		order.ShardKey, order.SmID, order.DateCreated, order.OofShard)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("не удалось вставить заказ: %w", err)
	}

	_, err = tx.Exec(`
		INSERT INTO delivery (order_uid, name, phone, zip, city, address, region, email)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, order.OrderUID, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip,
		order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("не удалось вставить данные о доставке: %w", err)
	}

	_, err = tx.Exec(`
		INSERT INTO payment (order_uid, transaction, request_id, currency, provider, 
		amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`, order.OrderUID, order.Payment.Transaction, order.Payment.RequestID, order.Payment.Currency,
		order.Payment.Provider, order.Payment.Amount, order.Payment.PaymentDT, order.Payment.Bank,
		order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("не удалось вставить данные об оплате: %w", err)
	}

	for _, item := range order.Items {
		_, err = tx.Exec(`
			INSERT INTO items (chrt_id, track_number, price, rid, name, sale, size, 
			total_price, nm_id, brand, status, order_uid)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		`, item.ChrtID, item.TrackNumber, item.Price, item.Rid, item.Name, item.Sale,
			item.Size, item.TotalPrice, item.NmID, item.Brand, item.Status, order.OrderUID)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("не удалось вставить товар: %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		log.Println("Failed to commit transaction", "error", err)
		tx.Rollback()
		return err
	}

	r.cache.Set(order.OrderUID, order)

	return nil
}

func (r *Repository) GetOrder(orderUID string) (*models.Order, error) {
	if val, ok := r.cache.Get(orderUID); ok {
		log.Println("Заказ найден в кеше")
		return val, nil
	}

	order := &models.Order{}

	err := r.db.QueryRow(`
		SELECT order_uid, track_number, entry, locale, internal_signature, 
		customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
		FROM orders WHERE order_uid = $1
	`, orderUID).Scan(
		&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale,
		&order.InternalSignature, &order.CustomerID,
		&order.DeliveryService, &order.ShardKey, &order.SmID, &order.DateCreated, &order.OofShard,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("не удалось получить заказ из таблицы orders: %w", err)
	}

	delivery := models.Delivery{}
	err = r.db.QueryRow(`
		SELECT name, phone, zip, city, address, region, email
		FROM delivery WHERE order_uid = $1
	`, orderUID).Scan(
		&delivery.Name, &delivery.Phone, &delivery.Zip, &delivery.City,
		&delivery.Address, &delivery.Region, &delivery.Email,
	)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Printf("Предупреждение: не удалось получить информацию о доставке для заказа %s: %v",
			orderUID, err)
	}
	order.Delivery = delivery

	payment := models.Payment{}
	err = r.db.QueryRow(`
		SELECT transaction, request_id, currency, provider, amount, payment_dt, bank, 
		delivery_cost, goods_total, custom_fee
		FROM payment WHERE order_uid = $1
	`, orderUID).Scan(
		&payment.Transaction, &payment.RequestID, &payment.Currency, &payment.Provider,
		&payment.Amount, &payment.PaymentDT,
		&payment.Bank, &payment.DeliveryCost, &payment.GoodsTotal, &payment.CustomFee,
	)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Printf("Предупреждение: не удалось получить информацию об оплате для заказа %s: %v",
			orderUID, err)
	}
	order.Payment = payment

	rows, err := r.db.Query(`
		SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status
		FROM items WHERE order_uid = $1
	`, orderUID)

	if err != nil {
		log.Printf("Предупреждение: не удалось получить товары для заказа %s: %v", orderUID, err)
		return order, nil
	}

	defer rows.Close()

	for rows.Next() {
		item := models.Item{}

		err := rows.Scan(
			&item.ChrtID, &item.TrackNumber, &item.Price, &item.Rid, &item.Name, &item.Sale,
			&item.Size, &item.TotalPrice, &item.NmID, &item.Brand, &item.Status,
		)

		if err != nil {
			log.Printf("Предупреждение: не удалось сканировать товар для заказа %s: %v", orderUID, err)
			continue
		}

		order.Items = append(order.Items, item)
	}

	r.cache.Set(orderUID, order)

	return order, nil
}

func (r *Repository) loadToCache() {
	rows, err := r.db.Query("SELECT order_uid FROM orders")
	if err != nil {
		log.Printf("не удалось получить все order : %v", err)
		return
	}

	defer rows.Close()

	for rows.Next() {
		var uid string
		if err := rows.Scan(&uid); err != nil {
			log.Printf("Предупреждение: не удалось сканировать order : %v", err)
			continue
		}

		order, err := r.GetOrder(uid)
		if err != nil || order == nil {
			log.Printf("Предупреждение: не удалось получить заказ %s: %v", uid, err)
			continue
		}

		r.cache.Set(uid, order)
	}
}
