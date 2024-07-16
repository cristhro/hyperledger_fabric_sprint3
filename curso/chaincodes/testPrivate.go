package main

import (
  "encoding/json"
  "encoding/base64"
  "fmt"
  "log"
  "time"

  "github.com/hyperledger/fabric-chaincode-go/shim"
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
      Type  	     string `json:"objectType"`
      ID             string `json:"ID"`
      Color          string `json:"color"`
      Size           int    `json:"size"`
      Owner          string `json:"owner"`
      Estado         string `json:"estado"`
	}

// AssetPrivateDetails describes details that are private to owners
type AssetPrivateDetails struct {
	ID             string `json:"ID"`
	AppraisedValue int    `json:"appraisedValue"`
}

const assetCollection = "assetCollection"

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
      {ID: myuuids, Color: "blue", Size: 5, Owner: "Tomoko", Estado: "Nuevo"},
      {ID: myuuid2, Color: "red", Size: 5, Owner: "Brad", Estado: "Nuevo"},
      {ID: myuuid3, Color: "green", Size: 10, Owner: "Jin Soo", Estado: "Nuevo"},
      {ID: myuuid4, Color: "yellow", Size: 10, Owner: "Max", Estado: "Nuevo"},
      {ID: myuuid5, Color: "black", Size: 15, Owner: "Adriana", Estado: "Nuevo"},
      {ID: myuuid6, Color: "white", Size: 15, Owner: "Michel", Estado: "Nuevo"},
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
   func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface) error {

// Get new asset from transient map
	transientMap, err := ctx.GetStub().GetTransient()
	if err != nil {
		return fmt.Errorf("error getting transient: %v", err)
	}

	// Asset properties are private, therefore they get passed in transient field, instead of func args
	transientAssetJSON, ok := transientMap["asset_properties"]
	if !ok {
		//log error to stdout
		return fmt.Errorf("asset not found in the transient map input")
	}

	type assetTransientInput struct {
		Type           string `json:"objectType"` //Type is used to distinguish the various types
		ID             string `json:"assetID"`
		Color          string `json:"color"`
		Size           int    `json:"size"`
		Estado	       string `json:"estado"`
		AppraisedValue int    `json:"appraisedValue"`
	}

	var assetInput assetTransientInput
	err = json.Unmarshal(transientAssetJSON, &assetInput)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	if len(assetInput.Type) == 0 {
		return fmt.Errorf("objectType field must be a non-empty string")
	}
	if len(assetInput.ID) == 0 {
		return fmt.Errorf("assetID field must be a non-empty string")
	}
	if len(assetInput.Color) == 0 {
		return fmt.Errorf("color field must be a non-empty string")
	}
	if assetInput.Size <= 0 {
		return fmt.Errorf("size field must be a positive integer")
	}
	if len(assetInput.Estado) == 0 {
                return fmt.Errorf("Estado field must be a non-empty string")
        }
	if assetInput.AppraisedValue <= 0 {
		return fmt.Errorf("appraisedValue field must be a positive integer")
	}

	// Check if asset already exists
	assetAsBytes, err := ctx.GetStub().GetPrivateData(assetCollection, assetInput.ID)
	if err != nil {
		return fmt.Errorf("failed to get asset: %v", err)
	} else if assetAsBytes != nil {
		fmt.Println("Asset already exists: " + assetInput.ID)
		return fmt.Errorf("this asset already exists: " + assetInput.ID)
	}

	// Get ID of submitting client identity
	clientID, err := submittingClientIdentity(ctx)
	if err != nil {
		return err
	}

	// Verify that the client is submitting request to peer in their organization
	// This is to ensure that a client from another org doesn't attempt to read or
	// write private data from this peer.
	err = verifyClientOrgMatchesPeerOrg(ctx)
	if err != nil {
		return fmt.Errorf("CreateAsset cannot be performed: Error %v", err)
	}

	myuuid := uuid.NewV4().String()
	fmt.Println("Your UUID is:", myuuid)

    asset := Asset{
      Type:   assetInput.Type,
      ID:     myuuid,
      Color:  assetInput.Color,
      Size:   assetInput.Size,
      Owner:  clientID,
      Estado: assetInput.Estado,
    }

    assetJSONasBytes, err := json.Marshal(asset)
	if err != nil {
		return fmt.Errorf("failed to marshal asset into JSON: %v", err)
	}

    log.Printf("CreateAsset Put: collection %v, ID %v, owner %v", assetCollection, myuuid, clientID)

    err = ctx.GetStub().PutPrivateData(assetCollection, myuuid, assetJSONasBytes)

    if err != nil {
		return fmt.Errorf("failed to put asset into private data collecton: %v", err)
	}

	// Save asset details to collection visible to owning organization
	assetPrivateDetails := AssetPrivateDetails{
		ID:             myuuid,
		AppraisedValue: assetInput.AppraisedValue,
	}
assetPrivateDetailsAsBytes, err := json.Marshal(assetPrivateDetails) // marshal asset details to JSON
	if err != nil {
		return fmt.Errorf("failed to marshal into JSON: %v", err)
	}

	// Get collection name for this organization.
	orgCollection, err := getCollectionName(ctx)
	if err != nil {
		return fmt.Errorf("failed to infer private collection name for the org: %v", err)
	}

	// Put asset appraised value into owners org specific private data collection
	log.Printf("Put: collection %v, ID %v", orgCollection, myuuid)
	err = ctx.GetStub().PutPrivateData(orgCollection, myuuid, assetPrivateDetailsAsBytes)
	if err != nil {
		return fmt.Errorf("failed to put asset private details: %v", err)
	}
	return nil
}

// verifyClientOrgMatchesPeerOrg is an internal function used verify client org id and matches peer org id.
func verifyClientOrgMatchesPeerOrg(ctx contractapi.TransactionContextInterface) error {
	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("failed getting the client's MSPID: %v", err)
	}
	peerMSPID, err := shim.GetMSPID()
	if err != nil {
		return fmt.Errorf("failed getting the peer's MSPID: %v", err)
	}

	if clientMSPID != peerMSPID {
		return fmt.Errorf("client from org %v is not authorized to read or write private data from an org %v peer", clientMSPID, peerMSPID)
	}

	return nil
}


func submittingClientIdentity(ctx contractapi.TransactionContextInterface) (string, error) {
	b64ID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return "", fmt.Errorf("Failed to read clientID: %v", err)
	}
	decodeID, err := base64.StdEncoding.DecodeString(b64ID)
	if err != nil {
		return "", fmt.Errorf("failed to base64 decode clientID: %v", err)
	}
	return string(decodeID), nil
}

// ReadAsset devuelve el Asset almacenado en el World State con su última modificación

func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, assetID string) (*Asset, error) {

	log.Printf("ReadAsset: collection %v, ID %v", assetCollection, assetID)
	assetJSON, err := ctx.GetStub().GetPrivateData(assetCollection, assetID) //get the asset from chaincode state
	if err != nil {
		return nil, fmt.Errorf("failed to read asset: %v", err)
	}

	//No Asset found, return empty response
	if assetJSON == nil {
		log.Printf("%v does not exist in collection %v", assetID, assetCollection)
		return nil, nil
	}

	var asset *Asset
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	return asset, nil

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

// getCollectionName is an internal helper function to get collection of submitting client identity.
func getCollectionName(ctx contractapi.TransactionContextInterface) (string, error) {

	// Get the MSP ID of submitting client identity
	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", fmt.Errorf("failed to get verified MSPID: %v", err)
	}

	// Create the collection name
	orgCollection := clientMSPID + "PrivateCollection"

	return orgCollection, nil
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
