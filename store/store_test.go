package store

import (
	"github.com/hyperledger/fabric-chaincode-go/shimtest"
	"testing"
)

func TestInit(t *testing.T) {
	scc := new(SmartContract)
	stub := shimtest.NewMockStub("fabcar_cc", scc)
	stub.MockInvoke("0",[][]byte{[]byte("initLedger")})
	stub.MockInvoke("0",[][]byte{[]byte("createCar"), []byte("CAR10"), []byte("MianBaoChe"), []byte("Black"), []byte("Zhang San"), []byte("Li Si")})
	stub.MockInvoke("0",[][]byte{[]byte("queryAllCars")})
	stub.MockInvoke("0",[][]byte{[]byte("queryCar"), []byte("CAR10")})
}
