/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package integration

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/hyperledger/fabric-sdk-go/api/apitxn"
	fabricTxn "github.com/hyperledger/fabric-sdk-go/pkg/fabric-txn"
)

func TestChainCodeInvoke(t *testing.T) {

	//testSetup := BaseSetupImpl{
	//	ConfigFile:      "../fixtures/config/config_test.yaml",
	//	ChannelID:       "mychannel",
	//	OrgID:           "peerorg1",
	//	ChannelConfig:   "../fixtures/channel/mychannel.tx",
	//	ConnectEventHub: true,
	//}

	testSetup := BaseSetupImpl{
		ConfigFile:      "../fixtures/config/config-test.yaml",
		ChannelID:       "mychannel",
		OrgID:           org1Name,
		ChannelConfig:   "../fixtures/channel-artifacts/channel.tx",
		ConnectEventHub: true,
	}

	if err := testSetup.Initialize_noCA1(); err != nil {
		t.Fatalf(err.Error())
	}

	//testSetup.ChainCodeID ="mycc"

	if err := testSetup.InstallAndInstantiateExampleCC(); err != nil {
		t.Fatalf("InstallAndInstantiateExampleCC return error: %v", err)
	}

	// Get Query value before invoke
	value, err := testSetup.QueryAsset()
	if err != nil {
		t.Fatalf("getQueryValue return error: %v", err)
	}
	fmt.Println("################FUNCTION:QUERY value:b:")
	fmt.Printf("*** QueryValue b before invoke %s\n", value)

	eventID := "test([a-zA-Z]+)"

	// Register callback for chaincode event
	done, rce := fabricTxn.RegisterCCEvent(testSetup.ChainCodeID, eventID, testSetup.EventHub)

	err = moveFunds(&testSetup)
	fmt.Println("################FUNCTION:INVOKE move 1 from a to b")
	if err != nil {
		t.Fatalf("Move funds return error: %v", err)
	}

	select {
	case <-done:
	case <-time.After(time.Second * 20):
		t.Fatalf("Did NOT receive CC for eventId(%s)\n", eventID)
	}

	testSetup.EventHub.UnregisterChaincodeEvent(rce)

	valueAfterInvoke, err := testSetup.QueryAsset()
	if err != nil {
		t.Errorf("getQueryValue return error: %v", err)
		return
	}
	fmt.Println("################FUNCTION:QUERY value:b:")
	fmt.Printf("*** QueryValue after invoke %s\n", valueAfterInvoke)

	valueInt, _ := strconv.Atoi(value)
	valueInt = valueInt + 1
	valueAfterInvokeInt, _ := strconv.Atoi(valueAfterInvoke)
	if valueInt != valueAfterInvokeInt {
		t.Fatalf("SendTransaction didn't change the QueryValue")
	}

	// Register callback for chaincode event
	done1, rce1 := fabricTxn.RegisterCCEvent(testSetup.ChainCodeID, eventID, testSetup.EventHub)

	err = insertFunds(&testSetup)
	fmt.Println("################FUNCTION:INVOKE insert c 100")
	if err != nil {
		t.Fatalf("Insert funds return error: %v", err)
	}

	select {
	case <-done1:
	case <-time.After(time.Second * 20):
		t.Fatalf("Did NOT receive CC for eventId(%s)\n", eventID)
	}

	testSetup.EventHub.UnregisterChaincodeEvent(rce1)

}

// moveFunds ...
func moveFunds(setup *BaseSetupImpl) error {
	fcn := "invoke"

	var args []string
	args = append(args, "move")
	args = append(args, "a")
	args = append(args, "b")
	args = append(args, "1")

	transientDataMap := make(map[string][]byte)
	transientDataMap["result"] = []byte("Transient data in move funds...")

	_, err := fabricTxn.InvokeChaincode(setup.Client, setup.Channel, []apitxn.ProposalProcessor{setup.Channel.PrimaryPeer()}, setup.EventHub, setup.ChainCodeID, fcn, args, transientDataMap)
	return err
}

// moveFunds ...
func insertFunds(setup *BaseSetupImpl) error {
	fcn := "invoke"

	var args []string
	args = append(args, "insert")
	args = append(args, "c")
//	args = append(args, "b")
	args = append(args, "1")

	transientDataMap := make(map[string][]byte)
	transientDataMap["result"] = []byte("Transient data in move funds...")

	_, err := fabricTxn.InvokeChaincode(setup.Client, setup.Channel, []apitxn.ProposalProcessor{setup.Channel.PrimaryPeer()}, setup.EventHub, setup.ChainCodeID, fcn, args, transientDataMap)
	return err
}
