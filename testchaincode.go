package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type Asset struct {
	ID		string `json:id`
	Owner	string `json:owner`
	Value 	int `json:value`
}

type SimpleAsset struct {
	contractapi.Contract	
}

func (s *SimpleAsset) InitLedger(ctx contractapi.TransactionContextInterface) error {
	assets := []Asset{
		{ID: "ASSET1", Owner: "Jipor", Value: 300 },
		{ID: "ASSET2", Owner: "Chipor", Value: 600 },
		{ID: "ASSET3", Owner: "Kimpor", Value: 900 },
	}
	for _, asset := range assets {
		assetJson, err := json.Marshal(asset)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(asset.ID, assetJson)
		if err != nil {
			return fmt.Errorf("failed to init asset %s: %v", asset.ID, err)
		}
	}
	return nil
}



func (s *SimpleAsset) QueryAsset(ctx contractapi.TransactionContextInterface, assetId string) (*Asset, error) {
	assetJson, err := ctx.GetStub().GetState(assetId)
    if err != nil {
        return nil, fmt.Errorf("failed to read asset %s: %v", assetId, err)
    }

	if assetJson == nil {
		return nil, fmt.Errorf("asset does not exist", assetId)
	}

	var asset Asset 
	err = json.Unmarshal(assetJson, &asset)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

func (s *SimpleAsset) CreateAsset(ctx contractapi.TransactionContextInterface, assetId string, owner string, value int) error {
	asset := Asset{
		ID: assetId,
		Owner: owner,
		Value: value,
	}

	assetJson, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(assetId, assetJson)

}

func (s *SimpleAsset) UpdateAsset(ctx contractapi.TransactionContextInterface, assetId string, newOwner string, newValue int) error {
	asset, err := s.QueryAsset(ctx, assetId)	
	if err != nil {
		return err
	}
	asset.Owner = newOwner
	asset.Value = newValue

	assetJson, err := json.Marshal(asset)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(assetId, assetJson)
}

func (s *SimpleAsset) DeleteAsset(ctx contractapi.TransactionContextInterface, assetId string) error {

	asset, err := s.QueryAsset(ctx, assetId)
	if err != nil {
		return err
	}

	return ctx.GetStub().DelState(asset.ID)
}

func (s *SimpleAsset) GetAssetsByRange(ctx contractapi.TransactionContextInterface, startKey string, endKey string) ([]*Asset, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey) 
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

// GetAssetsByRangeWithPagination retrieves assets by range with pagination support
// func (s *SimpleAsset) GetAssetsByRangeWithPagination(ctx contractapi.TransactionContextInterface, startKey string, endKey string, pageSize int32, bookmark string) ([]*Asset, *peer.QueryResponseMetadata, error){
//     resultsIterator, metadata, err := ctx.GetStub().GetStateByRangeWithPagination(startKey, endKey, pageSize, bookmark)
//     if err != nil {
//         return nil, nil, err
//     }
//     defer resultsIterator.Close()

//     var assets []*Asset
//     for resultsIterator.HasNext() {
//         queryResponse, err := resultsIterator.Next()
//         if err != nil {
//             return nil, nil, err
//         }

//         var asset Asset
//         err = json.Unmarshal(queryResponse.Value, &asset)
//         if err != nil {
//             return nil, nil, err
//         }
//         assets = append(assets, &asset)
//     }

//     return assets, metadata, nil
// }

func main()  {
	chaincode, err := contractapi.NewChaincode(new(SimpleAsset))
	if err != nil {
		fmt.Printf("Error creating simple test chaincode: %s", err)
		return
	}

	if err = chaincode.Start(); err != nil {
		fmt.Printf("Error starting simple test chaincode %s", err)
	}
}