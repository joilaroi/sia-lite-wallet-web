package main

import (
	"encoding/json"
	"log"
	"syscall/js"

	"github.com/siacentral/sia-lite-wallet-web/wasm/modules"
	siatypes "gitlab.com/NebulousLabs/Sia/types"
)

func main() {
	js.Global().Set("sia", map[string]interface{}{
		"generateSeed":       js.FuncOf(generateSeed),
		"generateAddresses":  js.FuncOf(generateAddresses),
		"recoverAddresses":   js.FuncOf(recoverAddresses),
		"getTransactions":    js.FuncOf(getTransactions),
		"encodeTransaction":  js.FuncOf(encodeTransaction),
		"signTransaction":    js.FuncOf(signTransaction),
		"encodeUnlockHash":   js.FuncOf(encodeUnlockHash),
		"encodeUnlockHashes": js.FuncOf(encodeUnlockHashes),
	})

	c := make(chan bool, 1)

	<-c
}

func checkArgs(args []js.Value, argTypes ...js.Type) bool {
	if len(args) != len(argTypes) {
		return false
	}

	for i, arg := range args {
		if arg.Type() != argTypes[i] {
			return false
		}
	}

	return true
}

func encodeTransaction(this js.Value, args []js.Value) interface{} {
	if !checkArgs(args, js.TypeString, js.TypeFunction) {
		return false
	}

	jsonTxn := args[0].String()
	callback := args[1]

	var txn siatypes.Transaction

	if err := json.Unmarshal([]byte(jsonTxn), &txn); err != nil {
		callback.Invoke(err.Error(), js.Null())
		return false
	}

	go modules.EncodeTransaction(txn, callback)

	return true
}

func signTransaction(this js.Value, args []js.Value) interface{} {
	var txn siatypes.Transaction

	if !checkArgs(args, js.TypeString, js.TypeString, js.TypeObject, js.TypeFunction) {
		return false
	}

	phrase := args[0].String()
	jsonTxn := args[1].String()
	length := args[2].Length()
	callback := args[3]
	requiredSigs := make([]uint64, length)

	if err := json.Unmarshal([]byte(jsonTxn), &txn); err != nil {
		callback.Invoke(err.Error(), js.Null())
		return false
	}

	for i := 0; i < length; i++ {
		requiredSigs[i] = uint64(args[2].Index(i).Int())
	}

	go modules.SignTransaction(txn, phrase, requiredSigs, callback)

	return true
}

func encodeUnlockHash(this js.Value, args []js.Value) interface{} {
	if !checkArgs(args, js.TypeString, js.TypeFunction) {
		return false
	}

	jsonUnlockConds := args[0].String()
	callback := args[1]

	go func() {
		var unlockConds siatypes.UnlockConditions

		if err := json.Unmarshal([]byte(jsonUnlockConds), &unlockConds); err != nil {
			callback.Invoke(err.Error(), js.Null())
			return
		}

		callback.Invoke(js.Null(), unlockConds.UnlockHash().String())
	}()

	return true
}

func encodeUnlockHashes(this js.Value, args []js.Value) interface{} {
	if !checkArgs(args, js.TypeObject, js.TypeFunction) {
		log.Println("wrong arguments")
		return false
	}

	length := args[0].Length()
	addresses := make([]interface{}, 0, length)
	callback := args[1]

	go func() {
		for i := 0; i < length; i++ {
			var uc siatypes.UnlockConditions

			if err := json.Unmarshal([]byte(args[0].Index(i).String()), &uc); err != nil {
				log.Println(err)
				callback.Invoke(err.Error(), js.Null())
				return
			}

			addresses = append(addresses, uc.UnlockHash().String())
		}

		callback.Invoke(js.Null(), addresses)
	}()

	return true
}

func generateSeed(this js.Value, args []js.Value) interface{} {
	if !checkArgs(args, js.TypeString, js.TypeFunction) {
		return false
	}

	seedType := args[0].String()
	callback := args[1]

	go modules.GenerateSeed(seedType, callback)

	return true
}

func generateAddresses(this js.Value, args []js.Value) interface{} {
	if !checkArgs(args, js.TypeString, js.TypeNumber, js.TypeNumber, js.TypeFunction) {
		return false
	}

	phrase := args[0].String()
	i := args[1].Int()
	n := args[2].Int()
	callback := args[3]

	go modules.GetAddresses(phrase, uint64(i), uint64(n), callback)

	return true
}

func recoverAddresses(this js.Value, args []js.Value) interface{} {
	if !checkArgs(args, js.TypeString, js.TypeNumber, js.TypeNumber, js.TypeNumber, js.TypeNumber, js.TypeFunction) {
		return false
	}

	seed := args[0].String()
	i := uint64(args[1].Int())
	maxEmptyRounds := uint64(args[2].Int())
	addressCount := uint64(args[3].Int())
	lastKnownIdx := uint64(args[4].Int())
	callback := args[5]

	go modules.RecoverAddresses(seed, i, maxEmptyRounds, addressCount, lastKnownIdx, callback)

	return true
}

func getTransactions(this js.Value, args []js.Value) interface{} {
	if !checkArgs(args, js.TypeObject, js.TypeFunction) {
		return false
	}

	count := args[0].Length()
	callback := args[1]
	addresses := make([]string, count)

	for i := 0; i < count; i++ {
		addresses[i] = args[0].Index(i).String()
	}

	go modules.GetTransactions(addresses, callback)

	return true
}
