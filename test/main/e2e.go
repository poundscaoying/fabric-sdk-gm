/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"fmt"

	"strconv"
	//"testing"
	"time"
	"sync"
	"strings"
	base "github.com/hyperledger/fabric-sdk-go/test/integration"
	//fabricTxn "github.com/hyperledger/fabric-sdk-go/fabric-txn"
	"os"
)

//var ch_p chan int

var waitgroup sync.WaitGroup

var total_t int
var current_t int

var count_p int
var count_t int

//var p_t int
//var plock sync.Mutex

var key_lab string
var q_all string
var org1Name = "peerorg1"




func init_channel()(ts base.BaseSetupImpl) {

	//testSetup := base.BaseSetupImpl{
	//	ConfigFile:      "../fixtures/config/config_test.yaml",
	//	ChannelID:       "mychannel",
	//	OrgID:           "peerorg1",
	//	ChannelConfig:   "../fixtures/channel-artifacts-nokafka/channel.kafka.tx",
	//	ConnectEventHub: true,
	//}

	//testSetup := base.BaseSetupImpl{
	//	ConfigFile:      "../fixtures/config/config_test.yaml",
	//	ChannelID:       "mychannel",
	//	OrgID:           "peerorg1",
	//	ChannelConfig:   "../fixtures/channel-artifacts/channel.tx",
	//	ConnectEventHub: true,
	//}

	testSetup := base.BaseSetupImpl{
		//ConfigFile:      "../fixtures/config/config_test.yaml",
		ConfigFile:      "../fixtures/config/config_test.yaml",
		ChannelID:       "mychannel",
		OrgID:           org1Name,
		ChannelConfig:   "../fixtures/channel-artifacts/channel.tx",
		ConnectEventHub: true,
	}

	if err := testSetup.Initialize_noCA1(); err != nil {
		fmt.Printf(err.Error())
	}
	testSetup.ChainCodeID ="mycc"

	if err := testSetup.InstallAndInstantiateExampleCC(); err != nil {
		fmt.Printf("InstallAndInstantiateExampleCC return error: %v", err)
	}
	ts=testSetup
	return ts
}

func query_value(testSetup base.BaseSetupImpl) {

	// Get Query value before invoke
	//value, err := testSetup.QueryAsset()
	value, err := testSetup.QueryAssetA1()
	if err != nil {
		fmt.Printf("getQueryValue return error: %v", err)
	}
	fmt.Printf("*** QueryValue ----------- %s\n", value)

}

func query_all_value(testSetup base.BaseSetupImpl) {

	// Get Query value before invoke
	for i := 1; i < count_t+1; i++ {
		for j:=1;j<count_p+1;j++{
			query_single_value(testSetup,j)
		}
	}

}


func query_single_value(testSetup base.BaseSetupImpl,i int) {

	// Get Query value before invoke
	key:=key_lab+strconv.Itoa(i)
	value, err := testSetup.QueryAssetSingle(key)
	fmt.Println("#######Get Query value before invoke key:",key,"value:",value)
	if err != nil {
		fmt.Printf("getQueryValue return error: %v", err)
	}
	
	if(!strings.EqualFold(key,value)){
		fmt.Printf("*** %s-----QueryValue -------- %s\n", key,value)
	}

}

func invoke_single(testSetup base.BaseSetupImpl,i int) {	
	defer waitgroup.Done()
	
	key:=key_lab+strconv.Itoa(i)
	//val:=key_lab+strconv.Itoa(i)
	val:=strconv.Itoa(i)
	fmt.Println("insert key is ",key)
	err := testSetup.InsertFunds(key,val)
	//err := testSetup.MoveFunds()
	if err != nil {
		fmt.Printf("Insert funds return error: %v", err)
	}
	
}


func invoke_pp(testSetup base.BaseSetupImpl) {	
	//fmt.Println("-------------pp---------------")
	for i := 0; i < count_p; i++ {
		waitgroup.Add(1)
		//current_t++
		//fmt.Println(current_t)

		go invoke_single(testSetup, i+1)
	}
	
	waitgroup.Wait()


}


func chaincode_invoke() {

	testSetup :=init_channel()
	time.Sleep(1e9*10)

	fmt.Println("--------------beigin---------------\n")

	time_b:=time.Now().UnixNano()
	//time_b := time.Now()

	for i := 0; i < count_t; i++ {

		invoke_pp(testSetup)
	}

	time_e:=time.Now().UnixNano()
	time_d := time_e - time_b
	tps := 1e9*int64(total_t)/time_d

	//time_e := time.Now()
	//duration := time_e.Sub(time_b)
	//nsecs := duration.Nanoseconds()
	//tps := 1e9*int64(total_t)/nsecs


	fmt.Println("\n--------------------staic--------------------------\n")
	fmt.Printf("\n------------time consuming: %d ms--------------------\n",time_d/1e6)
	//fmt.Printf("\n-------------time consuming: %d ms--------------------\n", nsecs/1e6)
	fmt.Printf("\n-------------total_t:%d--------------------\n",total_t)
	fmt.Printf("\n-------------tps:%d--------------------\n",tps)


	//query_all_value(testSetup)

}

func main(){

	args := os.Args

	if args != nil && len(args) >= 4{

		count_p, _ = strconv.Atoi(args[1])
		count_t, _ = strconv.Atoi(args[2])
		key_lab=args[3]


	}else{

		count_p=10
		count_t=1
		total_t=10
		key_lab="a"
		//q_all="1"
	}

	current_t=0

	total_t=count_p*count_t

	chaincode_invoke()
	//init_channel()
}
