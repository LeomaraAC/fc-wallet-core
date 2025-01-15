package main

import (
	"database/sql"
	"fmt"
	"context"

	"github.com.br/leomaraAC/fs-ms-wallet/pkg/events"
	"github.com.br/leomaraAC/fs-ms-wallet/internal/database"
	"github.com.br/leomaraAC/fs-ms-wallet/internal/event"
	"github.com.br/leomaraAC/fs-ms-wallet/internal/usecase/create_account"
	"github.com.br/leomaraAC/fs-ms-wallet/internal/usecase/create_client"
	"github.com.br/leomaraAC/fs-ms-wallet/internal/usecase/create_transaction"
	"github.com.br/leomaraAC/fs-ms-wallet/pkg/uow"
	"github.com.br/leomaraAC/fs-ms-wallet/internal/web"
	"github.com.br/leomaraAC/fs-ms-wallet/internal/web/webserver"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", "root", "root", "mysql", "3306", "wallet"))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	eventDispatcher := events.NewEventDispatcher()
	TransactionCreatedEvent := event.NewTransactionCreated()
	// eventDispatcher.Register("TransactionCreated", handler)

	clientDb := database.NewClientDB(db)
	accountDb := database.NewAccountDB(db)
	// transactionDb := database.NewTransactionDB(db)
	ctx := context.Background()
	uow := uow.NewUow(ctx, db)

	uow.Register("AccountDB", func(tx *sql.Tx) interface{} {
		return database.NewAccountDB(db)
	})

	uow.Register("TransactionDB", func(tx *sql.Tx) interface{} {
		return database.NewTransactionDB(db)
	})

	createClientUseCase := create_client.NewCreateClientUseCase(clientDb)
	createAccountUseCase := create_account.NewCreateAccountUseCase(accountDb, clientDb)
	createTransactionUseCase := create_transaction.NewCreateTransactionUseCase(uow, eventDispatcher, TransactionCreatedEvent)

	webserver := webserver.NewWebServer(":3000")

	clientHandler := web.NewWebClientHandler(*createClientUseCase)
	accountHandler := web.NewWebAccountHandler(*createAccountUseCase)
	transactionHandler := web.NewWebTransactionHandler(*createTransactionUseCase)

	webserver.AddHandler("/clients", clientHandler.CreateClient)
	webserver.AddHandler("/accounts", accountHandler.CreateAccount)
	webserver.AddHandler("/transactions", transactionHandler.CreateTransaction)

	webserver.Start()

}
