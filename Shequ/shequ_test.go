package Shequ

import (
	"github.com/hyperledger/fabric-chaincode-go/shimtest"
	"testing"
)

func TestAll(t *testing.T){
	scc := new(ShequContract)
	stub := shimtest.NewMockStub("shequ_cc", scc)
	stub.MockInvoke("0",[][]byte{[]byte("InitLedger")})
	stub.MockInvoke("0",[][]byte{[]byte("GetSQByToken"),[]byte("token1")})
}