package wallet

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/nekruz08/wallet/pkg/types"
)

var defaultTestAccount = testAccount{
	phone: "+992888844290",
	balance: 10_000_00,
	payments: []struct{
		amount types.Money
		category types.PaymentCategory
	}{
		{amount:1_000_00,category:"auto"},
	},
}

//-----------------------------------------------new
func TestService_FindPaymentByID_success(t *testing.T) {
	// создаём сервис
	s:=newTestService()
	_,payments,err:=s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return	
	}

	// попробуем найти платеж
	payment:=payments[0]
	got, err:=s.FindPaymentByID(payment.ID)
	if err != nil {
		t.Errorf("FindPaymentByID(): error = %v",err)
		return
	}

	// сравниваем платежи
	if !reflect.DeepEqual(payment, got){
		t.Errorf("FindPaymentByID(): wrong payment returned = %v",err)
		return
	}
}


//-----------------------------------------------new

func TestService_FindPaymentByID_fail(t *testing.T) {
	// создаем сервис
	s:=newTestService()
	_,_,err:=s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	// пробуем найти несуществующий платеж
	_, err = s.FindPaymentByID(uuid.New().String())
	if err == nil {
		t.Error("FindPaymentByID(): must return error, returned nil")
		return
	}

	if err != ErrPaymentNotFound {
		t.Errorf("FindPaymentByID(): must return ErrPaymentNotFound, returned = %v", err)
		return
	}
}

//-----------------------------------------------new

func TestService_Reject_success(t *testing.T) {
	// создаем сервис
	s:=newTestService()
	_, payments, err:=s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	// пробуем отменить платеж
	payment:=payments[0]
	err=s.Reject(payment.ID)
	if err != nil {
		t.Errorf("Reject(): error = %v", err)
		return
	}

	savedPayment, err:=s.FindPaymentByID(payment.ID)
	if err != nil {
		t.Errorf("Reject(): can't find payment by id, error = %v", err)
		return
	}
	if savedPayment.Status!=types.PaymentStatusFail{
		t.Errorf("Reject(): status didn't changd, payment %v", savedPayment)
		return
	}

	savedAccount, err:=s.FindAccountByID(payment.AccountID)
	if err != nil {
		t.Errorf("Reject(): can't find account by id, error = %v", err)
		return
	}
	if savedAccount.Balance!=defaultTestAccount.balance{
		t.Errorf("Reject(): balance didn't changed, account = %v", savedAccount)
		return	
	}
}

//-----------------------------------------------
func TestService_Repeat_success(t *testing.T) {
	// создаем сервис
	s:=newTestService()
	_, payments, err:=s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	// пробуем повторит платеж
	payment:=payments[0]
	newPayment,err:=s.Repeat(payment.ID)
	if err != nil {
		t.Errorf("Repeat(): error = %v", err)
		return
	}

	// сравниваем платежи
	if !reflect.DeepEqual(payment.AccountID, newPayment.AccountID){
		t.Error("FindPaymentByID(): payment.AccountID and newPayment.AccountID should not be equal")
		return
	}
}