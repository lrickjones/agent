package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an agent
type SmartContract struct {
	contractapi.Contract
}

// Agent describes basic details of what makes up a simple agent
type Agent struct {
	ID         string `json:"ID"`
	DID        string `json:"DID"`
	Name       string `json:"NAME"`
	Address    string `json:"ADDRESS"`
	Represents string `json:"REPRESENTS"`
	AgentType  string `json:"TYPE"`
	Roles      string `json:"ROLES"`
	IAM        string `json:"IAM"`
	Status     string `json:"STATUS"`
}

// InitLedger adds a base set of agents to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	agent := []Agent{
		{ID: "agent1", DID: "agent1@myssi.org", Name: "Number One", Address: "",
			Represents: "", AgentType: "", Roles: "", IAM: "", Status: ""},
		{ID: "agent2", DID: "agent2@myssi.org", Name: "Number Two", Address: "",
			Represents: "", AgentType: "", Roles: "", IAM: "", Status: ""},
		{ID: "agent3", DID: "agent3@myssi.org", Name: "Number Three", Address: "",
			Represents: "", AgentType: "", Roles: "", IAM: "", Status: ""},
	}

	for _, agent := range agent {
		agentJSON, err := json.Marshal(agent)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(agent.ID, agentJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	return nil
}

// Createagent issues a new agent to the world state with given details.
func (s *SmartContract) Createagent(ctx contractapi.TransactionContextInterface,
	id string, did string, name string, address string, represents string, agentType string,
	roles string, iam string) error {
	exists, err := s.agentExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the agent %s already exists", id)
	}

	agent := Agent{
		ID:         id,
		DID:        did,
		Name:       name,
		Address:    address,
		Represents: represents,
		AgentType:  agentType,
		Roles:      roles,
		IAM:        iam,
		Status:     "active",
	}
	agentJSON, err := json.Marshal(agent)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, agentJSON)
}

// Readagent returns the agent stored in the world state with given id.
func (s *SmartContract) Readagent(ctx contractapi.TransactionContextInterface, id string) (*Agent, error) {
	agentJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if agentJSON == nil {
		return nil, fmt.Errorf("the agent %s does not exist", id)
	}

	var agent Agent
	err = json.Unmarshal(agentJSON, &agent)
	if err != nil {
		return nil, err
	}

	return &agent, nil
}

// Updateagent updates an existing agent in the world state with provided parameters.
func (s *SmartContract) Updateagent(ctx contractapi.TransactionContextInterface,
	id string, did string, name string, address string, represents string, agentType string,
	roles string, iam string, status string) error {
	exists, err := s.agentExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the agent %s does not exist", id)
	}

	// overwriting original agent with new agent
	agent := Agent{
		ID:         id,
		DID:        did,
		Name:       name,
		Address:    address,
		Represents: represents,
		AgentType:  agentType,
		Roles:      roles,
		IAM:        iam,
		Status:     status,
	}
	agentJSON, err := json.Marshal(agent)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, agentJSON)
}

// Deleteagent deletes an given agent from the world state.
func (s *SmartContract) Deleteagent(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := s.agentExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the agent %s does not exist", id)
	}

	return ctx.GetStub().DelState(id)
}

// agentExists returns true when agent with given ID exists in world state
func (s *SmartContract) agentExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	agentJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return agentJSON != nil, nil
}

// GetAllagents returns all agents found in world state
func (s *SmartContract) GetAllagents(ctx contractapi.TransactionContextInterface) ([]*Agent, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all agents in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var agents []*Agent
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var agent Agent
		err = json.Unmarshal(queryResponse.Value, &agent)
		if err != nil {
			return nil, err
		}
		agents = append(agents, &agent)
	}

	return agents, nil
}

func main() {
	agentChaincode, err := contractapi.NewChaincode(&SmartContract{})
	if err != nil {
		log.Panicf("Error creating agent-transfer-basic chaincode: %v", err)
	}

	if err := agentChaincode.Start(); err != nil {
		log.Panicf("Error starting agent-transfer-basic chaincode: %v", err)
	}
}
