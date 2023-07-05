package main

import (
	"encoding/json"
	"fmt"
	"time"
	"log"
    "crypto/sha1"
    "encoding/hex"
	
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

type Hash struct {
	Signature   string `json:"signature"`
	Extradata	string `json:"extradata"`
	Time  		string `json:"time"`
}

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	hashs := []Hash{
		{Signature: "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9", Extradata: "hello world", Time: "1991"},
		{Signature: "4d13a6ebd0979cc031cf3650e8060b0457f17614bc62560ae795acc9faabfb98", Extradata: "Patente de Coca Cola", Time: "1991"},
	}

	for _, hash := range hashs {
		hashJSON, err := json.Marshal(hash)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(hash.Signature, hashJSON)
		if err != nil {
			return fmt.Errorf("Error al guardar en Fabric. %v", err)
		}
	}

	return nil
}

func (s *SmartContract) CreateHash(ctx contractapi.TransactionContextInterface, signature string, extradata string) error {
	exists, err := s.HashExists(ctx, signature)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("El hash %s ya existe", signature)
	}

    h := sha1.New()
    h.Write([]byte(signature))
    sha1_hash := hex.EncodeToString(h.Sum(nil))

    fmt.Println(s, sha1_hash)

	hash := Hash{
		Signature:  sha1_hash,
		Extradata:	extradata,
		Time: time.Now().Format(time.RFC3339),
	}
	hashJSON, err := json.Marshal(hash)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(sha1_hash, hashJSON)
}

func (s *SmartContract) ReadHash(ctx contractapi.TransactionContextInterface, signature string) (*Hash, error) {
	hashJSON, err := ctx.GetStub().GetState(signature)
	if err != nil {
		return nil, fmt.Errorf("Error al leer de Fabric: %v", err)
	}
	if hashJSON == nil {
		return nil, fmt.Errorf("El hash %s no existe", signature)
	}

	var hash Hash
	err = json.Unmarshal(hashJSON, &hash)
	if err != nil {
		return nil, err
	}

	return &hash, nil
}

func (s *SmartContract) HashExists(ctx contractapi.TransactionContextInterface, signature string) (bool, error) {
	hashJSON, err := ctx.GetStub().GetState(signature)
	if err != nil {
		return false, fmt.Errorf("Error al leer de Fabric: %v", err)
	}

	return hashJSON != nil, nil
}

func main() {
    hashChaincode, err := contractapi.NewChaincode(&SmartContract{})
    if err != nil {
		log.Panicf("Error al crear el chaincode Hash: %v", err)
    }

    if err := hashChaincode.Start(); err != nil {
		log.Panicf("Error al arrancar el chaincode Hash: %v", err)
    }
  }
