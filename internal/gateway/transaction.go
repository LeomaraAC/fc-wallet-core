package gateway

import "github.com.br/leomaraAC/fs-ms-wallet/internal/entity"

type TransactionGateway interface {
	Create(transaction *entity.Transaction) error
}