package create_transaction

import (
	"context"
	"github.com.br/leomaraAC/fs-ms-wallet/internal/gateway"
	"github.com.br/leomaraAC/fs-ms-wallet/internal/entity"
	"github.com.br/leomaraAC/fs-ms-wallet/pkg/events"
	"github.com.br/leomaraAC/fs-ms-wallet/pkg/uow"
)

type CreateTransactionInputDTO struct {
	AccountIDFrom string `json:"account_id_from"`
	AccountIDTo string `json:"account_id_to"`
	Amount float64 `json:"amount"`
}

type CreateTransactionOutputDTO struct {
	ID string `json:"id"`	
	AccountIDFrom string `json:"account_id_from"`
	AccountIDTo string `json:"account_id_to"`
	Amount float64 `json:"amount"`
}

type CreateTransactionUseCase struct {
	Uow uow.UowInterface
	EventDispatcher events.EventDispatcherInterface
	TransactionCreated events.EventInterface
}

func NewCreateTransactionUseCase(
	t gateway.TransactionGateway,
	a gateway.AccountGateway,
	eventDispatcher events.EventDispatcherInterface,
	transactionCreated events.EventInterface,
) *CreateTransactionUseCase {
	return &CreateTransactionUseCase{
		TransactionGateway: t,
		AccountGateway: a,
		EventDispatcher: eventDispatcher,
		TransactionCreated: transactionCreated,
	}
}

func (uc *CreateTransactionUseCase) Execute(ctx context.Context, input CreateTransactionInputDTO) (*CreateTransactionOutputDTO, error) {
	output := &CreateTransactionOutputDTO{}
	err := uc.Uow.Do(ctx, func(_ *uow.Uow) error {
		accountRepository := uc.getAccountRepository(ctx)
		transactionRepository := uc.getTransactionRepository(ctx)
		accountFrom, err := accountRepository.FindById(input.AccountIDFrom)
		if err != nil {
			return nil, err
		}

		accountTo, err := accountRepository.FindById(input.AccountIDTo)
		if err != nil {
			return nil, err
		}

		transaction, err := entity.NewTransaction(accountFrom, accountTo, input.Amount)
		if err != nil {
			return nil, err
		}

		err = accountRepository.UpdateBalance(accountFrom)
		if err != nil {
			return nil, err
		}

		err = accountRepository.UpdateBalance(accountTo)
		if err != nil {
			return nil, err
		}

		err = transactionRepository.Create(transaction)
		if err != nil {
			return nil, err
		}

		output.ID = transaction.ID
		output.AccountIDFrom = input.AccountIDFrom
		output.AccountIDTo = input.AccountIDTo
		output.Amount = input.Amount
		return nil
	})

	if err != nil {
		return nil, err
	}

	uc.TransactionCreated.SetPayload(output)
	uc.EventDispatcher.Dispatch(uc.TransactionCreated)

	return output, nil
}

func (uc *CreateTransactionUseCase) getAccountRepository(ctx context.Context) gateway.AccountRepository {
	repo, err := uc.Uow.GetRepository(ctx, "AccountDB")
	if err != nil {
		panic(err)
	}
	return repo.(gateway.AccountGateway)
}

func (uc *CreateTransactionUseCase) getTransactionRepository(ctx context.Context) gateway.TransactionRepository {
	repo, err := uc.Uow.GetRepository(ctx, "TransactionDB")
	if err != nil {
		panic(err)
	}
	return repo.(gateway.TransactionGateway)
}
