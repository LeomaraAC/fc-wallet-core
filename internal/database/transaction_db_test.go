package database

import (
	"testing"
	"database/sql"
	"github.com/stretchr/testify/suite"
	"github.com.br/leomaraAC/fs-ms-wallet/internal/entity"
	_"github.com/mattn/go-sqlite3"
)


type TransactionDBTestSuite struct {
	suite.Suite
	db *sql.DB
	transactionDB *TransactionDB 
	client *entity.Client
	client2 *entity.Client
	accountFrom *entity.Account
	accountTo *entity.Account
}

func (s *TransactionDBTestSuite) SetupSuite() {
	db, err := sql.Open("sqlite3", ":memory:")
	s.Nil(err)
	s.db = db
	db.Exec("CREATE TABLE clients(id VARCHAR(255), name VARCHAR(255), email VARCHAR(255), created_at date)")
	db.Exec("CREATE TABLE accounts(id VARCHAR(255), client_id VARCHAR(255), balance int, created_at date)")
	db.Exec("CREATE TABLE transactions(id VARCHAR(255), account_id_from VARCHAR(255), account_id_to VARCHAR(255), amount float, created_at date)")
	s.transactionDB = NewTransactionDB(db)
	client, err := entity.NewClient("John", "j@j")
	s.Nil(err)
	s.client = client
	client2, err := entity.NewClient("John2", "j2@j")
	s.Nil(err)
	s.client2 = client2
	accountFrom := entity.NewAccount(s.client)
	accountFrom.Balance = 1000
	s.accountFrom = accountFrom
	accountTo := entity.NewAccount(s.client2)
	accountTo.Balance = 1000
	s.accountTo = accountTo
}

func (s *TransactionDBTestSuite) TearDownSuite() {
	defer s.db.Close()
	s.db.Exec("DROP TABLE accounts")
	s.db.Exec("DROP TABLE clients")
	s.db.Exec("DROP TABLE transactions")
}

func TestTransactionDBTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionDBTestSuite))
}

func (s *TransactionDBTestSuite) TestCreate()  {
	transaction, err := entity.NewTransaction(s.accountFrom, s.accountTo, 100)
	s.Nil(err)
	err = s.transactionDB.Create(transaction)
	s.Nil(err)
}

