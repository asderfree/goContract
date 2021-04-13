package Shequ

import (
	"github.com/hyperledger/fabric-chaincode-go/shimtest"
	"testing"
)

func TestAll(t *testing.T) {
	scc := new(ShequContract)
	stub := shimtest.NewMockStub("shequ_cc", scc)
	stub.MockInvoke("0", [][]byte{[]byte("InitLedger")})
	stub.MockInvoke("0", [][]byte{[]byte("GetSQByToken"), []byte("token1")})
	stub.MockInvoke("0", [][]byte{[]byte("QueryAllSQ"), []byte("token0"), []byte("token10")})
	stub.MockInvoke("0", [][]byte{[]byte("AddNewSheQu"), []byte("token6"), []byte("affjhafhf")})
	stub.MockInvoke("0", [][]byte{[]byte("QueryAllSQ"), []byte("token0"), []byte("token10")})
}
