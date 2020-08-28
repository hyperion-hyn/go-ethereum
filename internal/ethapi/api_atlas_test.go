package ethapi

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"testing"
)

func TestGetAllValidatorAddresses(t *testing.T) {
	client, err := ethclient.Dial("http://localhost:8545")
	if err != nil {
		fmt.Printf("%v", err)
	}

	address, err := client.GetAllValidatorAddresses(context.Background(), nil)

	if err != nil {
		fmt.Printf("%v", err)
	}
	fmt.Printf("got address :%s", address)

}

func TestGetValidatorInformation(t *testing.T) {
	client, err := ethclient.Dial("http://localhost:8545")
	if err != nil {
		fmt.Printf("%v", err)
	}

	validatorAddress := common.HexToAddress("0xFD58E69Ebe3a2eF59181A87b811440DB3AC4f97a")

	validator, err := client.GetValidatorInformation(context.Background(), validatorAddress, nil)

	if err != nil {
		fmt.Printf("%v", err)
	}
	fmt.Printf("got validator :%v", validator.Validator.ValidatorAddress.Hex())

}
