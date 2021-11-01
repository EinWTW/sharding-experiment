package chaincode

import (
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

var namespace = hexdigest("sharding")[:6]
var lockspace = "L_"

// SmartContract provides functions for managing an Account
type SmartContract struct {
	contractapi.Contract
}

// Account describes basic details of what makes up a simple account
type Account struct {
	ID    string `json:"ID"`
	Value int    `json:"Value"`
}

type WLock struct {
	ID       string `json:"ID"`
	Checksum string `json:"Checksum"`
}

// InitLedger adds a base set of accounts to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface, prefix string, count int) error {

	for i := 1; i < count+1; i++ {
		account := Account{ID: prefix + strconv.Itoa(i), Value: 10000}
		err := s.saveAccount(ctx, account.ID, &account)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	return nil
}

func (s *SmartContract) AcquireLock(ctx contractapi.TransactionContextInterface, id string, identifier string) (bool, error) {
	lockId := lockKey(id)
	checksum := hexdigest(identifier)
	locksum, err := s.WLockExists(ctx, id)
	if err != nil {
		return false, err
	}
	if locksum == checksum { // wlock exists
		return true, nil
	}
	if locksum != "" {
		return false, fmt.Errorf("the write lock %s - %s is occupied", lockId, locksum)
	}

	wlock := WLock{
		ID:       lockId,
		Checksum: checksum,
	}

	wlockJSON, err := json.Marshal(wlock)
	if err != nil {
		return false, err
	}
	key := accountKey(lockId)
	err = ctx.GetStub().PutState(key, wlockJSON)
	if err != nil {
		return false, err
	}

	return true, nil
}

// DeleteLock deletes the write lock for an given account.
func (s *SmartContract) DeleteLock(ctx contractapi.TransactionContextInterface, id string, identifier string) error {
	lockId := lockKey(id)
	checksum := hexdigest(identifier)
	locksum, err := s.WLockExists(ctx, id)
	if err != nil {
		return err
	}
	if locksum == checksum { // wlock own
		err = ctx.GetStub().DelState(accountKey(lockId))
		if err != nil {
			return err
		} else {
			return nil
		}
	}
	if locksum != "" {
		return fmt.Errorf("the write lock %s - %s is occupied", lockId, locksum)
	}

	return nil
}

func (s *SmartContract) SendAmount(ctx contractapi.TransactionContextInterface, senderId string, receiverId string, amount int) error {

	senderAccount, err1 := s.ReadAccount(ctx, senderId)
	if err1 != nil {
		return err1
	}
	receiverAccount, err2 := s.ReadAccount(ctx, receiverId)
	if err2 != nil {
		return err2
	}

	senderAccount.Value -= amount
	receiverAccount.Value += amount

	err1 = s.saveAccount(ctx, senderId, senderAccount)
	if err1 != nil {
		return err1
	}
	err2 = s.saveAccount(ctx, receiverId, receiverAccount)
	if err2 != nil {
		return err2
	}

	return nil
}

// SendAmount update amount from account1 to account2 in world state.
func (s *SmartContract) SendAmountWithLock(ctx contractapi.TransactionContextInterface, senderId string, receiverId string, amount int) error {

	locksum1, err1 := s.WLockExists(ctx, senderId)
	if err1 != nil {
		return err1
	}
	locksum2, err2 := s.WLockExists(ctx, receiverId)
	if err2 != nil {
		return err2
	}
	if locksum1 == locksum2 {
		senderAccount, err1 := s.ReadAccount(ctx, senderId)
		if err1 != nil {
			return err1
		}
		receiverAccount, err2 := s.ReadAccount(ctx, receiverId)
		if err2 != nil {
			return err2
		}
		senderAccount.Value -= amount
		receiverAccount.Value += amount
		err1 = s.saveAccount(ctx, senderId, senderAccount)
		if err1 != nil {
			return err1
		}
		err2 = s.saveAccount(ctx, receiverId, receiverAccount)
		if err2 != nil {
			return err2
		}
		identifier := senderId + receiverId
		err1 = s.DeleteLock(ctx, senderId, identifier)
		if err1 != nil {
			return err1
		}
		err2 = s.DeleteLock(ctx, receiverId, identifier)
		if err2 != nil {
			return err2
		}
	}

	return nil
}

// SendAmount update amount from account1 to account2 in world state.
func (s *SmartContract) SendAmountCrossShards(ctx contractapi.TransactionContextInterface, senderId string, receiverId string, amount int) error {
	exists1, err1 := s.AccountExists(ctx, senderId)
	if err1 != nil {
		return err1
	}
	locksum1 := ""
	if exists1 {
		locksum1, err1 = s.WLockExists(ctx, senderId)
		if err1 != nil {
			return err1
		}
	} else {
		locksum1, err1 = s.WLockExistsCrossShards(ctx, senderId)
		if err1 != nil {
			return err1
		}
	}

	exists2, err2 := s.AccountExists(ctx, receiverId)
	if err2 != nil {
		return err2
	}
	locksum2 := ""
	if exists2 {
		locksum2, err2 = s.WLockExists(ctx, receiverId)
		if err1 != nil {
			return err1
		}
	} else {
		locksum2, err2 = s.WLockExistsCrossShards(ctx, receiverId)
		if err2 != nil {
			return err2
		}
	}

	identifier := senderId + receiverId
	checksum := hexdigest(identifier)
	if locksum1 == checksum && locksum2 == checksum {
		if exists1 {
			senderAccount, err1 := s.ReadAccount(ctx, senderId)
			if err1 != nil {
				return err1
			}
			senderAccount.Value -= amount
			err1 = s.saveAccount(ctx, senderId, senderAccount)
			if err1 != nil {
				return err1
			}
			//defer s.DeleteLock(ctx, senderId, identifier)
		}
		if exists2 {
			receiverAccount, err2 := s.ReadAccount(ctx, receiverId)
			if err2 != nil {
				return err2
			}
			receiverAccount.Value += amount
			err2 = s.saveAccount(ctx, receiverId, receiverAccount)
			if err2 != nil {
				return err2
			}
			//defer s.DeleteLock(ctx, receiverId, identifier)
		}

		// err1 = s.DeleteLock(ctx, senderId, identifier)
		// if err1 != nil {
		// 	return err1
		// }
		// err2 = s.DeleteLock(ctx, receiverId, identifier)
		// if err2 != nil {
		// 	return err2
		// }
	} else {
		return fmt.Errorf("the write locks for sender(%s - %s) and receiver(%s - %s) are occupied", senderId, locksum1, receiverId, locksum2)
	}

	return nil
}

func (s *SmartContract) WLockExistsCrossShards(ctx contractapi.TransactionContextInterface, id string) (string, error) {
	channel, smartcontract := ShardFormation(id)
	params := []string{"WLockExists", id}
	queryArgs := make([][]byte, len(params))
	for i, arg := range params {
		queryArgs[i] = []byte(arg)
	}

	response := ctx.GetStub().InvokeChaincode(smartcontract, queryArgs, channel)
	if response.Status != 200 {
		return "", fmt.Errorf("Failed to query chaincode. Got error: %s %s %s", response.Payload, channel, smartcontract)
	}

	return string(response.Payload), nil
}

func (s *SmartContract) WLockExists(ctx contractapi.TransactionContextInterface, id string) (string, error) {
	lockId := lockKey(id)
	key := accountKey(lockId)
	lockJSON, err := ctx.GetStub().GetState(key)
	if err != nil {
		return "", fmt.Errorf("failed to read from world state: %v", err)
	}
	if lockJSON == nil {
		return "", nil
	}

	var wlock WLock
	err = json.Unmarshal(lockJSON, &wlock)
	if err != nil {
		return "", err
	}

	return wlock.Checksum, nil
}

// GetBalance returns the account value stored in the world state with given id.
func (s *SmartContract) GetBalance(ctx contractapi.TransactionContextInterface, id string) (int, error) {
	key := accountKey(id)
	accountJSON, err := ctx.GetStub().GetState(key)
	if err != nil {
		return 0, fmt.Errorf("failed to read from world state: %v", err)
	}
	if accountJSON == nil {
		return 0, fmt.Errorf("the account %s does not exist", id)
	}

	var account Account
	err = json.Unmarshal(accountJSON, &account)
	if err != nil {
		return 0, err
	}

	return account.Value, nil
}

// CreateAccount issues a new account to the world state with given details.
func (s *SmartContract) CreateAccount(ctx contractapi.TransactionContextInterface, id string, value int) error {
	exists, err := s.AccountExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the account %s already exists", id)
	}

	account := Account{
		ID:    id,
		Value: value,
	}

	return s.saveAccount(ctx, id, &account)
}

// ReadAccount returns the account stored in the world state with given id.
func (s *SmartContract) ReadAccount(ctx contractapi.TransactionContextInterface, id string) (*Account, error) {
	key := accountKey(id)
	accountJSON, err := ctx.GetStub().GetState(key)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if accountJSON == nil {
		return nil, fmt.Errorf("the account %s does not exist", id)
	}

	var account Account
	err = json.Unmarshal(accountJSON, &account)
	if err != nil {
		return nil, err
	}

	return &account, nil
}

// UpdateAccount updates an existing account in the world state with provided parameters.
func (s *SmartContract) UpdateAccount(ctx contractapi.TransactionContextInterface, id string, value int) error {
	exists, err := s.AccountExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the account %s does not exist", id)
	}

	// overwriting original account with new account
	account := Account{
		ID:    id,
		Value: value,
	}
	return s.saveAccount(ctx, id, &account)
}

// DeleteAccount deletes an given account from the world state.
func (s *SmartContract) DeleteAccount(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := s.AccountExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the account %s does not exist", id)
	}

	return ctx.GetStub().DelState(accountKey(id))
}

// AccountExists returns true when account with given ID exists in world state
func (s *SmartContract) AccountExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	key := accountKey(id)
	accountJSON, err := ctx.GetStub().GetState(key)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return accountJSON != nil, nil
}

// GetAllAccounts returns all accounts found in world state
func (s *SmartContract) GetAllAccounts(ctx contractapi.TransactionContextInterface) ([]*Account, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all accounts in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var accounts []*Account
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var account Account
		err = json.Unmarshal(queryResponse.Value, &account)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, &account)
	}

	return accounts, nil
}

func (s *SmartContract) saveAccount(ctx contractapi.TransactionContextInterface, id string, account *Account) error {
	accountJSON, err := json.Marshal(account)
	if err != nil {
		return err
	}
	key := accountKey(account.ID)
	return ctx.GetStub().PutState(key, accountJSON)
}
func hexdigest(str string) string {
	hash := sha512.New()
	hash.Write([]byte(str))
	hashBytes := hash.Sum(nil)
	return strings.ToLower(hex.EncodeToString(hashBytes))
}

func lockKey(id string) string {
	return lockspace + id
}
func accountKey(id string) string {
	return namespace + id
}

// Requirements from sharding assessment
func ShardFormation(id string) (string, string) {
	if strings.HasPrefix(id, "A") {
		return "channel1", "sharding1"
	} else if strings.HasPrefix(id, "B") {
		return "channel2", "sharding2"
	} else {
		fmt.Println("Warning: only support for <account> in sharding1 and sharding2")
		return "mychannel", "sharding"
	}
}
