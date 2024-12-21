package database

import (
	"testing"
	"database/sql"
	"github.com/stretchr/testify/suite"
	"github.com.br/leomaraAC/fs-ms-wallet/internal/entity"
	_"github.com/mattn/go-sqlite3"
)


type ClientDBTestSuit struct {
	suite.Suite
	db *sql.DB
	clientDB *ClientDB 
}

func (s *ClientDBTestSuit) SetupSuite() {
	db, err := sql.Open("sqlite3", ":memory:")
	s.Nil(err)
	s.db = db
	db.Exec("CREATE TABLE clients(id VARCHAR(255), name VARCHAR(255), email VARCHAR(255), created_at date)")
	s.clientDB = NewClientDB(db)
}

func (s *ClientDBTestSuit) TearDownSuite() {
	defer s.db.Close()
	s.db.Exec("DROP TABLE clients")
}

func TestClientDBTestSuite(t *testing.T) {
	suite.Run(t, new(ClientDBTestSuit))
}

func (s *ClientDBTestSuit) TestSave() {
	client := &entity.Client{ID: "1", Name: "Test", Email: "j@j.com"}
	err := s.clientDB.Save(client)
	s.Nil(err)
}

func (s *ClientDBTestSuit) TestGet() {
	client, _ := entity.NewClient("John", "j@j.com")
	s.clientDB.Save(client)

	clientDB, err := s.clientDB.Get(client.ID)
	s.Nil(err)
	s.Equal(client.ID, clientDB.ID)
	s.Equal(client.Name, clientDB.Name)
	s.Equal(client.Email, clientDB.Email)
}
