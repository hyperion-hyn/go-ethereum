package ethapi

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"testing"
)

func TestHeader(t *testing.T) {
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
