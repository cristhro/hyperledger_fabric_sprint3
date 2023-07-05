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

type SmartContract struct {
      contractapi.Contract
    }

type Farmaco struct {
      ID             string `json:"ID"`
      Ph             string `json:"ph"`
      Temperature    int    `json:"temperature"`
      User           string `json:"user"`
	    Geo            string `json:"geo"`
      Estado         string `json:"estado"`
      Extradata      string `json:"extradata"`
	}

	type HistoryQueryResult struct {
		Record    *Farmaco    `json:"record"`
		TxId     string    `json:"txId"`
		Timestamp time.Time `json:"timestamp"`
		IsDelete  bool      `json:"isDelete"`
	}

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {

	myuuids := uuid.NewV4().String()
	fmt.Println("Your UUID is:", myuuids)

	myuuid2 := uuid.NewV4().String()
	fmt.Println("Your UUID2 is:", myuuid2)
	
    farmacos := []Farmaco{
      {ID: myuuids, Ph: "5.5", Temperature: 12, User: "USR122", Geo: "62 , -11", Estado: "Fabricado", Extradata: "-"},
      {ID: myuuid2, Ph: "5.4", Temperature: 12, User: "USR099", Geo: "45 , -66", Estado: "Fabricando", Extradata: "-"},
    }
    for _, farmaco := range farmacos {
      farmacoJSON, err := json.Marshal(farmaco)
      if err != nil {
        return err
      }
    
      err = ctx.GetStub().PutState(farmaco.ID, farmacoJSON)
      if err != nil {
        return fmt.Errorf("failed to put to world state. %v", err)
      }
    }

    return nil
  }

func (s *SmartContract) CreateFarmaco(ctx contractapi.TransactionContextInterface, ph string, temperature int, user string, geo string, estado string, extradata string) error {

	myuuid := uuid.NewV4().String()
	fmt.Println("Nuevo fármaco UUID:", myuuid)

    farmaco := Farmaco{
      ID:          myuuid,
      Ph:          ph,
      Temperature: temperature,
      User:        user,
	    Geo:         geo,
      Estado:		   estado,
      Extradata:   extradata,
    }
    farmacoJSON, err := json.Marshal(farmaco)
    if err != nil {
      return err
    }

    return ctx.GetStub().PutState(myuuid, farmacoJSON)
  }

// ReadFarmaco devuelve el Farmaco almacenado en el World State con su última modificación
   func (s *SmartContract) ReadFarmaco(ctx contractapi.TransactionContextInterface, id string) (*Farmaco, error) {
    farmacoJSON, err := ctx.GetStub().GetState(id)
    if err != nil {
      return nil, fmt.Errorf("Error al leer de Fabric: %v", err)
    }
    if farmacoJSON == nil {
      return nil, fmt.Errorf("El fármaco %s no existe", id)
    }

    var farmaco Farmaco
    err = json.Unmarshal(farmacoJSON, &farmaco)
    if err != nil {
      return nil, err
    }

    return &farmaco, nil
  }

func (s *SmartContract) UpdateFarmaco(ctx contractapi.TransactionContextInterface, id string, ph string, temperature int, user string, geo string, estado string, extradata string) error {
  exists, err := s.FarmacoExists(ctx, id)
  if err != nil {
      return err
    }
  if !exists {
      return fmt.Errorf("El fármaco %s no existe", id)
    }

  farmaco := Farmaco{
      ID:          id,
      Ph:          ph,
      Temperature: temperature,
      User:        user,
	    Geo:         geo,
      Estado:		   estado,
      Extradata:   extradata,
    }
  
  farmacoJSON, err := json.Marshal(farmaco)
  
  if err != nil {
      return err
    }

    return ctx.GetStub().PutState(id, farmacoJSON)
  }

  // DeleteFarmaco elimina a un Farmaco del world state
  func (s *SmartContract) DeleteFarmaco(ctx contractapi.TransactionContextInterface, id string) error {
    exists, err := s.FarmacoExists(ctx, id)
    if err != nil {
      return err
    }
    if !exists {
      return fmt.Errorf("the farmaco %s does not exist", id)
    }

    return ctx.GetStub().DelState(id)
  }

// FarmacoExists devuelve un true si existe el ID del Farmaco consultado
   func (s *SmartContract) FarmacoExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
    farmacoJSON, err := ctx.GetStub().GetState(id)
    if err != nil {
      return false, fmt.Errorf("failed to read from world state: %v", err)
    }

    return farmacoJSON != nil, nil
  }

// TransferFarmaco actualiza el dueño del Farmaco en el world state
   func (s *SmartContract) TransferFarmaco(ctx contractapi.TransactionContextInterface, id string, newUser string) error {
    farmaco, err := s.ReadFarmaco(ctx, id)
    if err != nil {
      return err
    }

    farmaco.User = newUser
    farmacoJSON, err := json.Marshal(farmaco)
    if err != nil {
      return err
    }

    return ctx.GetStub().PutState(id, farmacoJSON)
  }

// StatusFarmaco actualiza el estado del Farmaco en el world state
func (s *SmartContract) StatusFarmaco(ctx contractapi.TransactionContextInterface, id string, newStatus string) error {
    farmaco, err := s.ReadFarmaco(ctx, id)
    if err != nil {
      return err
    }

    farmaco.Estado = newStatus
    farmacoJSON, err := json.Marshal(farmaco)
    if err != nil {
      return err
    }

    return ctx.GetStub().PutState(id, farmacoJSON)
  }

// GetAllFarmacos devuelve todos los Farmacos del world state
   func (s *SmartContract) GetAllFarmacos(ctx contractapi.TransactionContextInterface) ([]*Farmaco, error) {
    resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
    if err != nil {
      return nil, err
    }
    defer resultsIterator.Close()

    var farmacos []*Farmaco
    for resultsIterator.HasNext() {
      queryResponse, err := resultsIterator.Next()
      if err != nil {
        return nil, err
      }

      var farmaco Farmaco
      err = json.Unmarshal(queryResponse.Value, &farmaco)
      if err != nil {
        return nil, err
      }
      farmacos = append(farmacos, &farmaco)
    }

    return farmacos, nil
  }

// GetFarmacoHistory devuelve la información histórica custodiada por el ledger
func (t *SmartContract) GetFarmacoHistory(ctx contractapi.TransactionContextInterface, farmacoID string) ([]HistoryQueryResult, error) {
	log.Printf("GetFarmacoHistory: ID %v", farmacoID)

	resultsIterator, err := ctx.GetStub().GetHistoryForKey(farmacoID)
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

		var farmaco Farmaco
		if len(response.Value) > 0 {
			err = json.Unmarshal(response.Value, &farmaco)
			if err != nil {
				return nil, err
			}
		} else {
			farmaco = Farmaco{
				ID: farmacoID,
			}
		}

		timestamp, err := ptypes.Timestamp(response.Timestamp)
		if err != nil {
			return nil, err
		}

		record := HistoryQueryResult{
			TxId:      response.TxId,
			Timestamp: timestamp,
			Record:    &farmaco,
			IsDelete:  response.IsDelete,
		}
		records = append(records, record)
	}

	return records, nil
}

  func main() {
    farmacoChaincode, err := contractapi.NewChaincode(&SmartContract{})
    if err != nil {
      log.Panicf("Error creating farmaco-transfer-basic chaincode: %v", err)
    }

    if err := farmacoChaincode.Start(); err != nil {
      log.Panicf("Error starting farmaco-transfer-basic chaincode: %v", err)
    }
  }

