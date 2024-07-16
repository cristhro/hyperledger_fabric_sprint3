package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Define el chaincode
type SmartContract struct {
	contractapi.Contract
}

// Estructura del alumno
type Alumno struct {
	ID       string `json:"id"`
	Nombre   string `json:"nombre"`
	Apellido string `json:"apellido"`
	Edad     int    `json:"edad"`
	Carrera  string `json:"carrera"`
}

// Inicialización del ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	alumnos := []Alumno{
		{ID: "1", Nombre: "Juan", Apellido: "Pérez", Edad: 21, Carrera: "Ingeniería"},
		{ID: "2", Nombre: "Ana", Apellido: "García", Edad: 22, Carrera: "Medicina"},
	}

	for _, alumno := range alumnos {
		alumnoAsBytes, _ := json.Marshal(alumno)
		err := ctx.GetStub().PutState(alumno.ID, alumnoAsBytes)

		if err != nil {
			return fmt.Errorf("Error al inicializar el ledger: %s", err.Error())
		}
	}

	return nil
}

// Crear un nuevo alumno
func (s *SmartContract) CrearAlumno(ctx contractapi.TransactionContextInterface, id string, nombre string, apellido string, edad int, carrera string) error {
	alumno := Alumno{
		ID:       id,
		Nombre:   nombre,
		Apellido: apellido,
		Edad:     edad,
		Carrera:  carrera,
	}

	alumnoAsBytes, _ := json.Marshal(alumno)
	return ctx.GetStub().PutState(id, alumnoAsBytes)
}

// Leer la información de un alumno
func (s *SmartContract) LeerAlumno(ctx contractapi.TransactionContextInterface, id string) (*Alumno, error) {
	alumnoAsBytes, err := ctx.GetStub().GetState(id)

	if err != nil {
		return nil, fmt.Errorf("No se puede leer el estado del mundo: %s", err.Error())
	}

	if alumnoAsBytes == nil {
		return nil, fmt.Errorf("El alumno %s no existe", id)
	}

	alumno := new(Alumno)
	_ = json.Unmarshal(alumnoAsBytes, alumno)

	return alumno, nil
}

// Actualizar la información de un alumno
func (s *SmartContract) ActualizarAlumno(ctx contractapi.TransactionContextInterface, id string, nombre string, apellido string, edad int, carrera string) error {
	alumno := Alumno{
		ID:       id,
		Nombre:   nombre,
		Apellido: apellido,
		Edad:     edad,
		Carrera:  carrera,
	}

	alumnoAsBytes, _ := json.Marshal(alumno)
	return ctx.GetStub().PutState(id, alumnoAsBytes)
}

// Eliminar un alumno
func (s *SmartContract) EliminarAlumno(ctx contractapi.TransactionContextInterface, id string) error {
	return ctx.GetStub().DelState(id)
}

// Obtener todos los alumnos
func (s *SmartContract) ObtenerTodosLosAlumnos(ctx contractapi.TransactionContextInterface) ([]*Alumno, error) {
	queryString := "{\"selector\":{}}"

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var alumnos []*Alumno
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var alumno Alumno
		_ = json.Unmarshal(queryResponse.Value, &alumno)
		alumnos = append(alumnos, &alumno)
	}

	return alumnos, nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(new(SmartContract))
	if err != nil {
		fmt.Printf("Error al crear el chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error al iniciar el chaincode: %s", err.Error())
	}
}
