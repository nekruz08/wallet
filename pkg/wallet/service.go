package wallet

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/nekruz08/wallet/pkg/types"
)

var ErrPhoneRegistered=errors.New("phone aready registered")
var ErrAmountMustBePositive=errors.New("amount must be greater than zero")
var ErrAccountNotFound=errors.New("account not found")
var ErrNotEnoughBalance=errors.New("not enough money")
var ErrPaymentNotFound=errors.New("payment not found")

type Service struct {
	nextAccountID int64
	accounts      []*types.Account
	payments      []*types.Payment
}

//-----------------------------------------------

func (s *Service) RegisterAccount(phone types.Phone) (*types.Account, error) {
	for _, account := range s.accounts {
		if account.Phone == phone {
			return nil, ErrPhoneRegistered
		}
	}
	s.nextAccountID++
	account := &types.Account{
		ID:      s.nextAccountID,
		Phone:   phone,
		Balance: 0,
	}
	s.accounts = append(s.accounts, account)
	return account, nil
}

//-----------------------------------------------

func (s *Service) Deposit(accountID int64, amount types.Money) error {
	if amount <= 0 {
		return ErrAmountMustBePositive
	}

	var account *types.Account
	for _, acc := range s.accounts {
		if acc.ID == accountID {
			account = acc
			break
		}
	}

	if account == nil {
		return ErrAccountNotFound
	}

	// зачисление средств пока не рассматриваем как платеж
	account.Balance += amount
	return nil
}

//-----------------------------------------------

func (s *Service) Pay(accountID int64,amount types.Money, category types.PaymentCategory) (*types.Payment,error) {
	if amount<=0{
		return nil, ErrAmountMustBePositive
	}

	var account *types.Account
	for _, acc := range s.accounts {
		if acc.ID==accountID{
			account=acc
			break
		}
	}

	if account==nil{
		return nil,ErrAccountNotFound
	}

	if account.Balance<amount{
		return nil,ErrNotEnoughBalance
	}

	account.Balance-=amount
	paymentID:=uuid.New().String()
	payment:=&types.Payment{
		ID: paymentID,
		AccountID: accountID,
		Amount: amount,
		Category: category,
		Status: types.PaymentStatusInProgress,
	}
	s.payments=append(s.payments, payment)
	return payment,nil
}

//-----------------------------------------------

func (s *Service) FindAccountByID(accountID int64) (*types.Account, error) {
	for _, account := range s.accounts {
		if account.ID == accountID {
			return account,nil
		}
	}
	return nil,ErrAccountNotFound
}

//-----------------------------------------------new

func (s *Service) Reject(paymentID string) error {
	payment,err:=s.FindPaymentByID(paymentID)
	if err != nil {
		return err
	}

	account,err:=s.FindAccountByID(payment.AccountID)
	if err != nil {
		return err
	}

	payment.Status=types.PaymentStatusFail
	account.Balance+=payment.Amount
	return nil
}

//-----------------------------------------------new

func (s *Service) FindPaymentByID(paymentID string) (*types.Payment,error) {
	for _, payment := range s.payments {
		if payment.ID==paymentID{
			return payment,nil
		}
	}
	return nil, ErrPaymentNotFound
}

//-----------------------------------------------new
type testService struct{
	//(встраивание)
	*Service					
}

//-----------------------------------------------new

//(функция конструктор)
func newTestService() *testService{	
	return &testService{Service: &Service{}}
}

//-----------------------------------------------new

type testAccount struct{
	phone types.Phone
	balance types.Money
	payments []struct{
		amount types.Money
		category types.PaymentCategory
	}
}

//-----------------------------------------------new


func (s *testService) addAccount(data testAccount) (*types.Account, []*types.Payment,error) {
	// регистрируем там пользователья
	account, err:=s.RegisterAccount(data.phone)
	if err != nil {
		return nil, nil,fmt.Errorf("can't register account, error = %v", err)
	}

	// пополняем его счёт
	err=s.Deposit(account.ID,data.balance)
	if err != nil {
		return nil, nil,fmt.Errorf("can't deposity account, error = %v", err)
	}

	// выполняем платежи
	// можем создать слайс сразу нужной длины, поскольку знаем размер
	payments:=make([]*types.Payment,len(data.payments))
	for i, payment := range data.payments {
		// тогда здесь работаем просто через index, а не через append
		payments[i],err=s.Pay(account.ID,payment.amount,payment.category)
		if err != nil {
			return nil, nil, fmt.Errorf("can't make payment, error = %v",err)
		}
	}

	return account, payments, nil
}

//-----------------------------------------------new
func(s *Service)Repeat(paymentID string) (*types.Payment, error){
	payment,err:=s.FindPaymentByID(paymentID)
	if err != nil {
		return nil,err
	}

	newPayment,err:=s.Pay(payment.AccountID,payment.Amount,payment.Category)
	if err != nil {
		return nil,err
	}
	return newPayment,nil
}
