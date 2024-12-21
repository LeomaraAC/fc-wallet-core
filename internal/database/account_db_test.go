package database

import (
	"testing"
	"database/sql"
	"github.com/stretchr/testify/suite"
	"github.com.br/leomaraAC/fs-ms-wallet/internal/entity"
	_"github.com/mattn/go-sqlite3"
)


type AccountDBTestSuit struct {
	suite.Suite
	db *sql.DB
	accountDB *AccountDB 
	client *entity.Client
}

func (s *AccountDBTestSuit) SetupSuite() {
	db, err := sql.Open("sqlite3", ":memory:")
	s.Nil(err)
	s.db = db
	db.Exec("CREATE TABLE clients(id VARCHAR(255), name VARCHAR(255), email VARCHAR(255), created_at date)")
	db.Exec("CREATE TABLE accounts(id VARCHAR(255), client_id VARCHAR(255), balance int, created_at date)")
	s.accountDB = NewAccountDB(db)
	s.client, _ = entity.NewClient("John", "j@j")
}

func (s *AccountDBTestSuit) TearDownSuite() {
	defer s.db.Close()
	s.db.Exec("DROP TABLE accounts")
	s.db.Exec("DROP TABLE clients")
}

func TestAccountDBTestSuite(t *testing.T) {
	suite.Run(t, new(AccountDBTestSuit))
}

func (s *AccountDBTestSuit) TestSave() {
	account := entity.NewAccount(s.client)
	err := s.accountDB.Save(account)
	s.Nil(err)
}

func (s *AccountDBTestSuit) TestFindById() {
	s.db.Exec("INSERT INTO clients (id, name, email, created_at) VALUES (?, ?, ?, ?)", s.client.ID, s.client.Name, s.client.Email, s.client.CreatedAt)


	account := entity.NewAccount(s.client)
	s.accountDB.Save(account)

	accountDB, err := s.accountDB.FindById(account.ID)
	s.Nil(err)
	s.Equal(account.ID, accountDB.ID)
	s.Equal(account.Client.ID, accountDB.Client.ID)
	s.Equal(account.Balance, accountDB.Balance)
	s.Equal(account.Client.Name, accountDB.Client.Name)
	s.Equal(account.Client.Email, accountDB.Client.Email)
}

func (s *AccountDBTestSuit) TestUpdateBalance() {
	account := entity.NewAccount(s.client)
	s.accountDB.Save(account)

	account.Balance = 100
	err := s.accountDB.UpdateBalance(account)
	s.Nil(err)

	accountDB, _ := s.accountDB.FindById(account.ID)
	s.Equal(account.Balance, accountDB.Balance)

	account.Balance = 80
	err = s.accountDB.UpdateBalance(account)
	s.Nil(err)

	accountDB, _ = s.accountDB.FindById(account.ID)
	s.Equal(account.Balance, accountDB.Balance)
}
