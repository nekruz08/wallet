package wallet

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/nekruz08/wallet/pkg/types"
)

func TestService_FindAccountByID_success(t *testing.T) {
	svc := &Service{}
	account,err:=svc.RegisterAccount("+992888844290")
	if err!=nil{
		fmt.Println(err)
		return
	}

	account1,err:=svc.FindAccountByID(account.ID)
	if err!=nil{
		switch err{
		case ErrAccountNotFound:
			fmt.Println("Аккаунт пользователья не найден")	
		}
		return
	}

	if !reflect.DeepEqual(account,account1){
		t.Errorf("invalid result, expected: %v, actual: %v",account,account1)
	}
}

//-----------------------------------------------

func TestService_FindAccountByID_notExist(t *testing.T) {
	svc := &Service{}
	account:=int64(2)
	account1,err:=svc.FindAccountByID(account)
	if err!=nil{
		switch err{
		case ErrAccountNotFound:
			fmt.Println("Аккаунт пользователья не найден")	
		}
		return
	}

	if reflect.DeepEqual(account,account1){
		t.Errorf("invalid result, expected: %v, actual: %v",account,account1)
	}
}

//-----------------------------------------------

func TestSevice_Reject_success(t *testing.T) {
	svc := &Service{}
	account,err:=svc.RegisterAccount("+992888844290")
	if err!=nil{
		fmt.Println(err)
		return
	}

	err=svc.Deposit(account.ID,10)
	if err!=nil{
		switch err{
		case ErrAmountMustBePositive:
			fmt.Println("Сумма должна быть положительной")
		case ErrAccountNotFound:
			fmt.Println("Аккаунт пользователья не найден")	
		}
		return
	}

	payment,err:=svc.Pay(account.ID,5,"sadaqa")
	if err!=nil{
		switch err{
		case ErrAmountMustBePositive:
			fmt.Println("Сумма должна быть положительной")
		case ErrAccountNotFound:
			fmt.Println("Аккаунт пользователья не найден")	
		}
		return
	}

	err=svc.Reject(payment.ID)
	if err!=nil{
		switch err{
		case ErrAmountMustBePositive:
			fmt.Println("Сумма должна быть положительной")
		case ErrAccountNotFound:
			fmt.Println("Аккаунт пользователья не найден")	
		}
		return
	}

	acc,err:=svc.FindAccountByID(payment.AccountID)
	if err!=nil{
		switch err{
		case ErrAmountMustBePositive:
			fmt.Println("Сумма должна быть положительной")
		case ErrAccountNotFound:
			fmt.Println("Аккаунт пользователья не найден")	
		}
		return
	}

	if !reflect.DeepEqual(account,acc){
		t.Errorf("invalid result, expected: %v, actual: %v",account,acc)
	}
}

//-----------------------------------------------
func TestSevice_Reject_notFound(t *testing.T) {
	svc := &Service{}
	account,err:=svc.RegisterAccount("+992888844290")
	if err!=nil{
		fmt.Println(err)
		return
	}

	err=svc.Deposit(account.ID,10)
	if err!=nil{
		switch err{
		case ErrAmountMustBePositive:
			fmt.Println("Сумма должна быть положительной")
		case ErrAccountNotFound:
			fmt.Println("Аккаунт пользователья не найден")	
		}
		return
	}

	payment,err:=svc.Pay(account.ID,5,"sadaqa")
	if err!=nil{
		switch err{
		case ErrAmountMustBePositive:
			fmt.Println("Сумма должна быть положительной")
		case ErrAccountNotFound:
			fmt.Println("Аккаунт пользователья не найден")	
		}
		return
	}
	
	wrongPayment:=&types.Payment{
		ID: "helloWorld",
	}
	err=svc.Reject(wrongPayment.ID)
	if err!=nil{
		switch err{
		case ErrPaymentNotFound:
			fmt.Println("ErrPaymentNotFound")
		case ErrAccountNotFound:
			fmt.Println("Аккаунт пользователья не найден")	
		}
		return
	}

	acc,err:=svc.FindAccountByID(payment.AccountID)
	if err!=nil{
		switch err{
		case ErrAmountMustBePositive:
			fmt.Println("Сумма должна быть положительной")
		case ErrAccountNotFound:
			fmt.Println("Аккаунт пользователья не найден")	
		}
		return
	}

	if reflect.DeepEqual(wrongPayment,acc){
		t.Errorf("invalid result, expected: %v, actual: %v",wrongPayment,acc)
	}
}