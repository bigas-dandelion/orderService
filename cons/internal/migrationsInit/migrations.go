package migrationsinit

import (
	"database/sql"
	"fmt"
	"log"
)

func ApplyMigrations(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS orders (
		order_uid VARCHAR(255) PRIMARY KEY, 
		track_number VARCHAR(255) NOT NULL, 
		entry VARCHAR(255),                 
		locale VARCHAR(255),                
		internal_signature VARCHAR(255),    
		customer_id VARCHAR(255),           
		delivery_service VARCHAR(255),      
		shardkey VARCHAR(255),              
		sm_id INTEGER,                      
		date_created TIMESTAMP WITH TIME ZONE, 
		oof_shard VARCHAR(255)              
	);

	CREATE TABLE IF NOT EXISTS delivery (
		order_uid VARCHAR(255) PRIMARY KEY REFERENCES orders(order_uid), 
		name VARCHAR(255),                  
		phone VARCHAR(255),                 
		zip VARCHAR(255),                  
		city VARCHAR(255),                 
		address VARCHAR(255),              
		region VARCHAR(255),                
		email VARCHAR(255)                  
	);

	CREATE TABLE IF NOT EXISTS payment (
		order_uid VARCHAR(255) PRIMARY KEY REFERENCES orders(order_uid), 
		transaction VARCHAR(255) NOT NULL,  
		request_id VARCHAR(255),
		currency VARCHAR(255),              
		provider VARCHAR(255),              
		amount INTEGER,                     
		payment_dt INTEGER,                 
		bank VARCHAR(255),                  
		delivery_cost INTEGER,              
		goods_total INTEGER,                
		custom_fee INTEGER                  
	);

	CREATE TABLE IF NOT EXISTS items (
		chrt_id INTEGER,                    
		track_number VARCHAR(255),          
		price INTEGER,                      
		rid VARCHAR(255) PRIMARY KEY,      
		name VARCHAR(255),                 
		sale INTEGER,                      
		size VARCHAR(255),                 
		total_price INTEGER,                
		nm_id INTEGER,                     
		brand VARCHAR(255),               
		status INTEGER,                    
		order_uid VARCHAR(255) REFERENCES orders(order_uid) 
	);`)

	if err != nil {
		return fmt.Errorf("ошибка  миграций: %w", err)
	}

	log.Println("Миграции базы данных успешно сделаны")
	
	return nil
}
