package arangodb

import (
	"context"
	"dongtzu/constant"
	"dongtzu/pkg/model"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/arangodb/go-driver"
	"gitlab.geax.io/demeter/gologger/logger"
)

const (
	CollectionProviders = "Providers"
)

func GetProviders(ctx context.Context) ([]*model.Provider, int) {
	results := []*model.Provider{}

	query := fmt.Sprintf(`
		FOR d IN %s
			FILTER d.lineAtChannelSecret != "" AND
				d.lineAtAccessToken != ""
			RETURN d`, CollectionProviders)
	bindVars := map[string]interface{}{}
	cursor, err := db.Query(ctx, query, bindVars)
	defer closeCursor(cursor)
	if err != nil {
		logger.Errorf("[ArangoDB] GetProviders query failed: %v", err)
		return []*model.Provider{}, constant.ArangoDB_Driver_Failed
	}

	for {
		var doc *model.Provider
		_, err := cursor.ReadDocument(ctx, &doc)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			logger.Errorf("[ArangoDB] GetProviders cursor failed: %v", err)
			return []*model.Provider{}, constant.ArangoDB_Driver_Failed
		}

		results = append(results, doc)
	}

	return results, constant.ArangoDB_Success
}

func GetProviderProfileByLineUserID(ctx context.Context, lineUserID string) (*model.Provider, error) {
	var result model.Provider

	query := fmt.Sprintf(`
		FOR p IN %s 
			FILTER p.lineUserId == @lineUserId
		RETURN p`, CollectionProviders)
	bindVars := map[string]interface{}{
		"lineUserId": lineUserID,
	}
	cursor, err := db.Query(ctx, query, bindVars)
	defer func() {
		if err := cursor.Close(); err != nil {
			logger.Errorf("[ArangoDB] cursor close falied: %v", err)
		}
	}()

	if err != nil {
		logger.Errorf("[ArangoDB] GetProviderProfileByLineUserID query falied: %v", err)
		return nil, err
	}

	for {
		_, err := cursor.ReadDocument(ctx, &result)
		if driver.IsNotFound(err) {
			return nil, nil
		} else if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			logger.Errorf("[ArangoDB] GetProviderProfileByLineUserID cursor falied: %v", err)
			return nil, err
		}
	}

	return &result, nil
}

func CreateProviderProfile(ctx context.Context, registerInfo model.RegisterProviderReq) error {
	savedMap, err := formatCreateProviderMap(registerInfo)
	if err != nil {
		return err
	}

	aql := `
function (Params) {
	const db = require('@arangodb').db;
	const providerCol = db._collection("Providers");
	const savedObj = Params[0];

	providerDoc = providerCol.firstExample({lineUserId: savedObj.lineUserId});
	if (providerDoc && providerDoc._key && providerDoc._key.length > 0 &&
		(providerDoc.status === 0 || providerDoc.status === 1)) {
		providerDoc = {
			...providerDoc,
			...savedObj
		}
		providerCol.update({ _key: providerDoc._key }, providerDoc);
	} else {
		providerCol.insert(savedObj);
	}

	return 1;
}
`
	options := &driver.TransactionOptions{
		MaxTransactionSize: 100000,
		WriteCollections:   []string{CollectionProviders},
		ReadCollections:    []string{CollectionProviders},
		Params:             []interface{}{savedMap},
		WaitForSync:        false,
	}

	resCode, err := db.Transaction(ctx, aql, options)
	if err != nil {
		logger.Errorf("CreateProviderProfile TX execution failure")
		return err
	}

	logger.Debugf("CreateProviderProfile TX execution resCode is : %v\n", resCode)

	return nil
}

func UpdateProviderByLineUserID(ctx context.Context, lineUserID string, providerInfo model.UpdateProviderInfoReq) error {
	savedMap, err := formatUpdateProviderMap(providerInfo)
	if err != nil {
		return err
	}

	aql := `
function (Params) {
	const db = require('@arangodb').db;
	const providerCol = db._collection("Providers");
	const lineUserId = Params[0];
	const savedObj = Params[1];

	providerDoc = providerCol.firstExample({lineUserId: lineUserId});
	if (providerDoc && providerDoc._key && providerDoc._key.length > 0 &&
		(providerDoc.status >= 1)) {
		providerDoc = {
			...providerDoc,
			...savedObj
		}
		providerCol.update({ _key: providerDoc._key }, providerDoc);
	} else {
		return 2;
	}

	return 1;
}
`
	options := &driver.TransactionOptions{
		MaxTransactionSize: 100000,
		WriteCollections:   []string{CollectionProviders},
		ReadCollections:    []string{CollectionProviders},
		Params:             []interface{}{lineUserID, savedMap},
		WaitForSync:        false,
	}

	resCode, err := db.Transaction(ctx, aql, options)
	if err != nil {
		logger.Errorf("UpdateProviderByLineUserID TX execution failure")
		return err
	}
	if resCode == 2 {
		return errors.New("Not found document")
	}

	logger.Debugf("UpdateProviderByLineUserID TX execution resCode is : %v", resCode)

	return nil
}

/**
 * Private methods
 */
func formatCreateProviderMap(registerInfo model.RegisterProviderReq) (map[string]interface{}, error) {
	dataMap := map[string]interface{}{}
	var now int
	if registerInfo.Status == constant.Provider_Status_Saved {
		now = int(time.Now().Unix())
	}
	provider := model.Provider{
		LineUserID:          registerInfo.LineUserID,
		RealName:            registerInfo.RealName,
		LineAtName:          registerInfo.LineAtName,
		LineAtID:            registerInfo.LineAtID,
		LineAtChannelID:     "",
		LineAtChannelSecret: "",
		LineAtAccessToken:   "",
		CountryCode:         "886",
		LineID:              registerInfo.LineID,
		PhoneNum:            registerInfo.PhoneNum,
		ConfirmedPhoneNum:   registerInfo.PhoneNum,
		GmailAddr:           registerInfo.GmailAddr,
		ConfirmedGmailAddr:  registerInfo.GmailAddr,
		GCalSync:            false,
		InviteCode:          registerInfo.InivteCode,
		MemeberTerm:         false,
		PrivacyTerm:         false,
		Status:              registerInfo.Status,
		CreatedAt:           now,
	}
	j, _ := json.Marshal(provider)
	_ = json.Unmarshal(j, &dataMap)
	delete(dataMap, "_key")
	return dataMap, nil
}

func formatUpdateProviderMap(providerInfo model.UpdateProviderInfoReq) (map[string]interface{}, error) {
	dataMap := map[string]interface{}{}
	provider := model.Provider{
		RealName:           providerInfo.RealName,
		LineAtName:         providerInfo.LineAtName,
		LineAtID:           providerInfo.LineAtID,
		CountryCode:        "886",
		LineID:             providerInfo.LineID,
		PhoneNum:           providerInfo.PhoneNum,
		ConfirmedPhoneNum:  providerInfo.PhoneNum,
		GmailAddr:          providerInfo.GmailAddr,
		ConfirmedGmailAddr: providerInfo.GmailAddr,
		Status:             providerInfo.Status,
	}
	j, _ := json.Marshal(provider)
	_ = json.Unmarshal(j, &dataMap)
	delete(dataMap, "_key")
	delete(dataMap, "lineUserId")
	delete(dataMap, "lineAtChannelId")
	delete(dataMap, "lineAtChannelSecret")
	delete(dataMap, "lineAtAccessToken")
	delete(dataMap, "gCalSync")
	delete(dataMap, "inviteCode")
	if providerInfo.Status != constant.Provider_Status_Auditing {
		delete(dataMap, "memeberTerm")
		delete(dataMap, "privacyTerm")
	}
	delete(dataMap, "createdAt")

	return dataMap, nil
}
