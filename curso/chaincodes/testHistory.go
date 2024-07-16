package main

import (
  "encoding/json"
  "fmt"
  "log"
  "time"

  "github.com/golang/protobuf/ptypes"
  "github.com/satori/go.uuid"
  "github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provee las funciones necesarias para manejar los Assets
   type SmartContract struct {
      contractapi.Contract
    }

// La estructura Asset especifica los campos que se usarán para registrar y modificar el Asset
   type Asset struct {
      ID             string `json:"ID"`
      Color          string `json:"color"`
      Size           int    `json:"size"`
      Owner          string `json:"owner"`
	  AppraisedValue int    `json:"appraisedValue"`
	  Estado         string `json:"estado"`
	}

//Estructura de datos para mostrar el histórico de cada Asset
	type HistoryQueryResult struct {
		Record    *Asset    `json:"record"`
		TxId     string    `json:"txId"`
		Timestamp time.Time `json:"timestamp"`
		IsDelete  bool      `json:"isDelete"`
	}

// InitLedger añade una lista de activos Asset para arrancar el ejemplo
   func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {

	myuuids := uuid.NewV4().String()
	fmt.Println("Your UUID is:", myuuids)

	myuuid2 := uuid.NewV4().String()
	fmt.Println("Your UUID2 is:", myuuid2)

	myuuid3 := uuid.NewV4().String()
	fmt.Println("Your UUID3 is:", myuuid3)

	myuuid4 := uuid.NewV4().String()
	fmt.Println("Your UUID4 is:", myuuid4)

	myuuid5 := uuid.NewV4().String()
	fmt.Println("Your UUID5 is:", myuuid5)

	myuuid6 := uuid.NewV4().String()
	fmt.Println("Your UUID6 is:", myuuid6)
	
    assets := []Asset{
      {ID: myuuids, Color: "blue", Size: 5, Owner: "Tomoko", AppraisedValue: 300, Estado: "Nuevo"},
      {ID: myuuid2, Color: "red", Size: 5, Owner: "Brad", AppraisedValue: 400, Estado: "Nuevo"},
      {ID: myuuid3, Color: "green", Size: 10, Owner: "Jin Soo", AppraisedValue: 500, Estado: "Nuevo"},
      {ID: myuuid4, Color: "yellow", Size: 10, Owner: "Max", AppraisedValue: 600, Estado: "Nuevo"},
      {ID: myuuid5, Color: "black", Size: 15, Owner: "Adriana", AppraisedValue: 700, Estado: "Nuevo"},
      {ID: myuuid6, Color: "white", Size: 15, Owner: "Michel", AppraisedValue: 800, Estado: "Nuevo"},
    }

    for _, asset := range assets {
      assetJSON, err := json.Marshal(asset)
      if err != nil {
        return err
      }

      err = ctx.GetStub().PutState(asset.ID, assetJSON)
      if err != nil {
        return fmt.Errorf("failed to put to world state. %v", err)
      }
    }

    return nil
  }

// CreateAsset genera un nuevo Asset con los siguientes detalles
   func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, color string, size int, owner string, appraisedValue int, estado string) error {

	myuuid := uuid.NewV4().String()
	fmt.Println("Your UUID is:", myuuid)

    asset := Asset{
      ID:             myuuid,
      Color:          color,
      Size:           size,
      Owner:          owner,
	  AppraisedValue: appraisedValue,
	  Estado:		  estado,
    }
    assetJSON, err := json.Marshal(asset)
    if err != nil {
      return err
    }

    return ctx.GetStub().PutState(myuuid, assetJSON)
  }

// ReadAsset devuelve el Asset almacenado en el World State con su última modificación
   func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, id string) (*Asset, error) {
    assetJSON, err := ctx.GetStub().GetState(id)
    if err != nil {
      return nil, fmt.Errorf("failed to read from world state: %v", err)
    }
    if assetJSON == nil {
      return nil, fmt.Errorf("the asset %s does not exist", id)
    }

    var asset Asset
    err = json.Unmarshal(assetJSON, &asset)
    if err != nil {
      return nil, err
    }

    return &asset, nil
  }

// UpdateAsset actualiza un Asset en el world state indicando previamente los campos a modificar
   func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface, id string, color string, size int, owner string, appraisedValue int, estado string) error {
    exists, err := s.AssetExists(ctx, id)
    if err != nil {
      return err
    }
    if !exists {
      return fmt.Errorf("the asset %s does not exist", id)
    }

    // overwriting original asset with new asset
    asset := Asset{
      ID:             id,
      Color:          color,
      Size:           size,
      Owner:          owner,
	  AppraisedValue: appraisedValue,
	  Estado:		  estado,
    }
    assetJSON, err := json.Marshal(asset)
    if err != nil {
      return err
    }

    return ctx.GetStub().PutState(id, assetJSON)
  }

  // DeleteAsset elimina a un Asset del world state
  func (s *SmartContract) DeleteAsset(ctx contractapi.TransactionContextInterface, id string) error {
    exists, err := s.AssetExists(ctx, id)
    if err != nil {
      return err
    }
    if !exists {
      return fmt.Errorf("the asset %s does not exist", id)
    }

    return ctx.GetStub().DelState(id)
  }

// AssetExists devuelve un true si existe el ID del Asset consultado
   func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
    assetJSON, err := ctx.GetStub().GetState(id)
    if err != nil {
      return false, fmt.Errorf("failed to read from world state: %v", err)
    }

    return assetJSON != nil, nil
  }

// TransferAsset actualiza el dueño del Asset en el world state
   func (s *SmartContract) TransferAsset(ctx contractapi.TransactionContextInterface, id string, newOwner string) error {
    asset, err := s.ReadAsset(ctx, id)
    if err != nil {
      return err
    }

    asset.Owner = newOwner
    assetJSON, err := json.Marshal(asset)
    if err != nil {
      return err
    }

    return ctx.GetStub().PutState(id, assetJSON)
  }

// StatusAsset actualiza el estado del Asset en el world state
func (s *SmartContract) StatusAsset(ctx contractapi.TransactionContextInterface, id string, newStatus string) error {
    asset, err := s.ReadAsset(ctx, id)
    if err != nil {
      return err
    }

    asset.Estado = newStatus
    assetJSON, err := json.Marshal(asset)
    if err != nil {
      return err
    }

    return ctx.GetStub().PutState(id, assetJSON)
  }

// GetAllAssets devuelve todos los Assets del world state
   func (s *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]*Asset, error) {
    resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
    if err != nil {
      return nil, err
    }
    defer resultsIterator.Close()

    var assets []*Asset
    for resultsIterator.HasNext() {
      queryResponse, err := resultsIterator.Next()
      if err != nil {
        return nil, err
      }

      var asset Asset
      err = json.Unmarshal(queryResponse.Value, &asset)
      if err != nil {
        return nil, err
      }
      assets = append(assets, &asset)
    }

    return assets, nil
  }

// GetAssetHistory devuelve la información histórica custodiada por el ledger
func (t *SmartContract) GetAssetHistory(ctx contractapi.TransactionContextInterface, assetID string) ([]HistoryQueryResult, error) {
	log.Printf("GetAssetHistory: ID %v", assetID)

	resultsIterator, err := ctx.GetStub().GetHistoryForKey(assetID)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var records []HistoryQueryResult
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset Asset
		if len(response.Value) > 0 {
			err = json.Unmarshal(response.Value, &asset)
			if err != nil {
				return nil, err
			}
		} else {
			asset = Asset{
				ID: assetID,
			}
		}

		timestamp, err := ptypes.Timestamp(response.Timestamp)
		if err != nil {
			return nil, err
		}

		record := HistoryQueryResult{
			TxId:      response.TxId,
			Timestamp: timestamp,
			Record:    &asset,
			IsDelete:  response.IsDelete,
		}
		records = append(records, record)
	}

	return records, nil
}

  func main() {
    assetChaincode, err := contractapi.NewChaincode(&SmartContract{})
    if err != nil {
      log.Panicf("Error creating asset-transfer-basic chaincode: %v", err)
    }

    if err := assetChaincode.Start(); err != nil {
      log.Panicf("Error starting asset-transfer-basic chaincode: %v", err)
    }
  }
