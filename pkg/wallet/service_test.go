package wallet

import (
	"testing"
	"fmt"
	"reflect"
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
