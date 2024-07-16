package main

import (
  "encoding/json"
  "fmt"
  "log"

  "github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
      contractapi.Contract
    }

type Transaccion struct {
      ID             string `json:"ID"`
      SapId          string `json:"sapid"`
      Cantidad       int    `json:"cantidad"`
      Precio         int    `json:"precio"`
    }

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
    transaccions := []Transaccion{
      {ID: "T1", SapId: "918723", Cantidad: 583, Precio: 300},
      {ID: "T2", SapId: "202021", Cantidad: 9616,  Precio: 400},
    }

    for _, transaccion := range transaccions {
      transaccionJSON, err := json.Marshal(transaccion)
      if err != nil {
        return err
      }

      err = ctx.GetStub().PutState(transaccion.ID, transaccionJSON)
      if err != nil {
        return fmt.Errorf("failed to put to world state. %v", err)
      }
    }

    return nil
  }

func (s *SmartContract) CreateTransaccion(ctx contractapi.TransactionContextInterface, id string, sapid string, cantidad int, precio int) error {
    exists, err := s.TransaccionExists(ctx, id)
    if err != nil {
      return err
    }
    if exists {
      return fmt.Errorf("the transaccion %s already exists", id)
    }

    if precio < 299 {
      return fmt.Errorf("No se puede vender tan barato")
    }

      transaccion := Transaccion{
        ID:             id,
        SapId:          sapid,
        Cantidad:           cantidad,
        Precio: precio,
      }
      transaccionJSON, err := json.Marshal(transaccion)
      if err != nil {
        return err
      }

      return ctx.GetStub().PutState(id, transaccionJSON)
  }

func (s *SmartContract) ReadTransaccion(ctx contractapi.TransactionContextInterface, id string) (*Transaccion, error) {
    transaccionJSON, err := ctx.GetStub().GetState(id)
    if err != nil {
      return nil, fmt.Errorf("failed to read from world state: %v", err)
    }
    if transaccionJSON == nil {
      return nil, fmt.Errorf("the transaccion %s does not exist", id)
    }

    var transaccion Transaccion
    err = json.Unmarshal(transaccionJSON, &transaccion)
    if err != nil {
      return nil, err
    }

    return &transaccion, nil
  }


func (s *SmartContract) UpdateTransaccion(ctx contractapi.TransactionContextInterface, id string, sapid string, cantidad int, precio int) error {
      exists, err := s.TransaccionExists(ctx, id)
      if err != nil {
        return err
      }
      if !exists {
        return fmt.Errorf("the transaccion %s does not exist", id)
      }

      if precio < 299 {
        return fmt.Errorf("No se puede vender tan barato")
      }

      asset, err := s.ReadTransaccion(ctx, id)
      
      fmt.Println("EEEOOO la cantidad es %d", asset.Cantidad)

      if asset.Cantidad < cantidad {
        return fmt.Errorf("No hay suficiente stock")
      }

      transaccion := Transaccion{
        ID:             id,
        SapId:          sapid,
        Cantidad:           cantidad,
        Precio: precio,
      }
      transaccionJSON, err := json.Marshal(transaccion)
      if err != nil {
        return err
      }

      return ctx.GetStub().PutState(id, transaccionJSON)
  }


func (s *SmartContract) DeleteTransaccion(ctx contractapi.TransactionContextInterface, id string) error {
    exists, err := s.TransaccionExists(ctx, id)
    if err != nil {
      return err
    }
    if !exists {
      return fmt.Errorf("the transaccion %s does not exist", id)
    }

    return ctx.GetStub().DelState(id)
  }

// TransaccionExists returns true when transaccion with given ID exists in world state
   func (s *SmartContract) TransaccionExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
    transaccionJSON, err := ctx.GetStub().GetState(id)
    if err != nil {
      return false, fmt.Errorf("failed to read from world state: %v", err)
    }

    return transaccionJSON != nil, nil
  }

func (s *SmartContract) GetAllTransaccions(ctx contractapi.TransactionContextInterface) ([]*Transaccion, error) {
    resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
    if err != nil {
      return nil, err
    }
    defer resultsIterator.Close()

    var transaccions []*Transaccion
    for resultsIterator.HasNext() {
      queryResponse, err := resultsIterator.Next()
      if err != nil {
        return nil, err
      }

      var transaccion Transaccion
      err = json.Unmarshal(queryResponse.Value, &transaccion)
      if err != nil {
        return nil, err
      }
      transaccions = append(transaccions, &transaccion)
    }

    return transaccions, nil
  }

func main() {
    transaccionChaincode, err := contractapi.NewChaincode(&SmartContract{})
    if err != nil {
      log.Panicf("Error creating transaccion-transfer-basic chaincode: %v", err)
    }

    if err := transaccionChaincode.Start(); err != nil {
      log.Panicf("Error starting transaccion-transfer-basic chaincode: %v", err)
    }
  }


