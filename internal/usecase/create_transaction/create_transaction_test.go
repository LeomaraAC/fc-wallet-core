package create_transaction

import (
	"testing"
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com.br/leomaraAC/fs-ms-wallet/internal/entity"
	"github.com.br/leomaraAC/fs-ms-wallet/pkg/events"
	"github.com.br/leomaraAC/fs-ms-wallet/internal/event"
	"github.com.br/leomaraAC/fs-ms-wallet/internal/usecase/mocks"
)

type TransactionGatewayMock struct {
	mock.Mock
}

func (m *TransactionGatewayMock) Create(transaction *entity.Transaction) error {
	args := m.Called(transaction)
	return args.Error(0)
}


type AccountGatewayMock struct {
	mock.Mock
}

func (m *AccountGatewayMock) Save(account *entity.Account) error {
	args := m.Called(account)
	return args.Error(0)
}

func (m *AccountGatewayMock) FindById(id string) (*entity.Account, error) {
	args := m.Called(id)
	return args.Get(0).(*entity.Account), args.Error(1)
}

func (m *AccountGatewayMock) UpdateBalance(account *entity.Account) error {
	args := m.Called(account)
	return args.Error(0)
}

func TestCreateTransactionUseCasr_Execute(t *testing.T) {
	client1, _ := entity.NewClient("client1", "j@j1.com")
	account1 := entity.NewAccount(client1)
	account1.Credit(1000)

	client2, _ := entity.NewClient("client2", "j@j2.com")
	account2 := entity.NewAccount(client2)
	account2.Credit(1000)

	mockUow := &mocks.UowMock{}
	mockUow.On("Do", mock.Anything, mock.Anything).Return(nil)

	// accountMock := &AccountGatewayMock{}
	// accountMock.On("FindById", account1.ID).Return(account1, nil)
	// accountMock.On("FindById", account2.ID).Return(account2, nil)

	// transactionMock := &TransactionGatewayMock{}
	// transactionMock.On("Create", mock.Anything).Return(nil)

	inputDTO := CreateTransactionInputDTO {
		AccountIDFrom: account1.ID,
		AccountIDTo: account2.ID,
		Amount: 100,
	}

	dispatcher := events.NewEventDispatcher()
	event := event.NewTransactionCreated()
	ctx := context.Background()

	uc := NewCreateTransactionUseCase(mockUow, dispatcher, event)

	output, err := uc.Execute(ctx, inputDTO)
	assert.Nil(t, err)
	assert.NotNil(t, output.ID)
	mockUow.AssertExpectations(t)
	mockUow.AssertNumberOfCalls(t, "Do", 1)
	// transactionMock.AssertExpectations(t)
	// accountMock.AssertExpectations(t)
	// transactionMock.AssertNumberOfCalls(t, "Create", 1)
	// accountMock.AssertNumberOfCalls(t, "FindById", 2)
}
