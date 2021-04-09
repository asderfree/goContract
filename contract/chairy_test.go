package contract

import (
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-chaincode-go/shimtest"
	"testing"
)
func checkInvoke(t *testing.T, stub *shimtest.MockStub, args [][]byte) {
	res := stub.MockInvoke("0", args)
	if res.Status != shim.OK {
		fmt.Println("Invoke", args, "failed", string(res.Message))
		t.FailNow()
	}
}

func checkDonation(t *testing.T, stub *shimtest.MockStub, name string) {
	res := stub.MockInvoke("0", [][]byte{[]byte("queryUserInfo"), []byte(name)})
	if res.Status != shim.OK {
		fmt.Println("Query", name, "failed", string(res.Message))
		t.FailNow()
	}
	if res.Payload == nil {
		fmt.Println("Query", name, "failed to get value")
		t.FailNow()
	}
	fmt.Println("Query value", "name is ", name, "value is ", res.Payload)
	t.FailNow()
}
func TestA1(t *testing.T)  {
	fmt.Println("Teat")
	scc := new(SmartContract)
	stub := shimtest.NewMockStub("charity", scc)
	checkInvoke(t,stub, [][]byte{[]byte("invoke"), []byte("A"), []byte("B"), []byte("123")})
	checkDonation(t, stub, "charity")
}


