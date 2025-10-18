package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type User struct {
	ID      int     `db:"id"`
	Name    string  `db:"name"`
	Email   string  `db:"email"`
	Balance float64 `db:"balance"`
}

func ConnectDB() (*sqlx.DB, error) {
	dsn := "user=postgres password=1234 dbname=usersdb sslmode=disable host=localhost port=5432"
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func InsertUser(db *sqlx.DB, user User) error {
	query := `INSERT INTO users (name, email, balance) VALUES (:name, :email, :balance)`
	_, err := db.NamedExec(query, user)
	return err
}

func GetAllUsers(db *sqlx.DB) ([]User, error) {
	var users []User
	query := `SELECT * FROM users ORDER BY id`
	err := db.Select(&users, query)
	return users, err
}

func GetUserByID(db *sqlx.DB, id int) (User, error) {
	var user User
	query := `SELECT * FROM users WHERE id=$1`
	err := db.Get(&user, query, id)
	return user, err
}
func TransferBalance(db *sqlx.DB, fromID int, toID int, amount float64) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	var fromUser User
	if err := tx.Get(&fromUser, "SELECT * FROM users WHERE id=$1", fromID); err != nil {
		tx.Rollback()
		return fmt.Errorf("sender not found: %v", err)
	}

	if fromUser.Balance < amount {
		tx.Rollback()
		return fmt.Errorf("insufficient funds for user %d", fromID)
	}

	_, err = tx.Exec("UPDATE users SET balance = balance - $1 WHERE id=$2", amount, fromID)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("UPDATE users SET balance = balance + $1 WHERE id=$2", amount, toID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("receiver update failed: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func main() {
	db, err := ConnectDB()
	if err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}
	defer db.Close()

	fmt.Println("✅ Connected to PostgreSQL")

	newUser := User{Name: "David", Email: "david@example.com", Balance: 120.0}
	if err := InsertUser(db, newUser); err != nil {
		log.Println("InsertUser error:", err)
	} else {
		log.Println("✅ New user inserted: David")
	}

	users, err := GetAllUsers(db)
	if err != nil {
		log.Println("GetAllUsers error:", err)
	} else {
		fmt.Println("\n📋 Users list:")
		for _, u := range users {
			fmt.Printf("%d | %s | %.2f\n", u.ID, u.Name, u.Balance)
		}
	}

	fmt.Println("\n💸 Transferring 20.00 from Bob (ID 2) → Charlie (ID 3)")
	if err := TransferBalance(db, 2, 3, 20.0); err != nil {
		log.Println("Transfer error:", err)
	} else {
		log.Println("✅ Transfer successful")
	}

	users, _ = GetAllUsers(db)
	fmt.Println("\n💰 Updated Balances:")
	for _, u := range users {
		fmt.Printf("%d | %s | %.2f\n", u.ID, u.Name, u.Balance)
	}
}
