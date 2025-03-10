package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAccounts() ([]*Account, error)
	GetAccountById(int) (*Account, error)
	GetAccountByNumber(int) (*Account, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	connStr := "user=postgres dbname=postgres password=Admin123 sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{db: db}, nil
}

func (s *PostgresStore) Init() error {
	return s.createAccountTable()
}
func (s *PostgresStore) createAccountTable() error {
	query := `CREATE TABLE IF NOT EXISTS account (
	    id SERIAL PRIMARY KEY,
        first_name VARCHAR(50),
        last_name VARCHAR(50),
		number serial,
        balance DECIMAL(10, 2) DEFAULT 0,
		encrypted_password VARCHAR(255),
		created_at TIMESTAMP
		)
		`
	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) CreateAccount(acc *Account) error {
	query := `
	INSERT INTO 
		account 
	(first_name, last_name, number, balance, encrypted_password, created_at)
		VALUES
	($1, $2, $3, $4, $5, $6)
	`
	resp, err := s.db.Query(query, acc.FirstName, acc.LastName, acc.Number, acc.Balance, acc.EncryptedPassword, acc.CreateAt)
	if err != nil {
		return err
	}

	fmt.Printf("%v\n", resp)
	return nil
}
func (s *PostgresStore) UpdateAccount(account *Account) error {
	return nil
}

func (s *PostgresStore) DeleteAccount(id int) error {
	_, err := s.db.Query("DELETE FROM account where id = $1", id)
	return err
}

func (s *PostgresStore) GetAccounts() ([]*Account, error) {
	query := `SELECT * FROM account`

	rows, err := s.db.Query(query)

	if err != nil {
		return nil, err
	}

	accounts := []*Account{}
	for rows.Next() {
		account, err := scanIntoAccount(rows)
		if err != nil {
			return nil, err
		}

		accounts = append(accounts, account)

	}

	return accounts, nil
}

func (s *PostgresStore) GetAccountById(id int) (*Account, error) {
	rows, err := s.db.Query("SELECT * FROM account where id = $1", id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoAccount(rows)
	}
	return nil, fmt.Errorf("account %d not found", id)
}

func (s *PostgresStore) GetAccountByNumber(number int) (*Account, error) {
	rows, err := s.db.Query("SELECT * FROM account where number = $1", number)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoAccount(rows)
	}
	return nil, fmt.Errorf("account %d not found", number)
}

func scanIntoAccount(rows *sql.Rows) (*Account, error) {
	account := new(Account)
	err := rows.Scan(
		&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.Number,
		&account.Balance,
		&account.EncryptedPassword,
		&account.CreateAt,
	)

	if err != nil {
		return nil, err
	}
	return account, nil
}
