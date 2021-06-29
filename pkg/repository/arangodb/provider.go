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
	"github.com/google/uuid"
	"gitlab.geax.io/demeter/gologger/logger"
)

func GetProviders(ctx context.Context) ([]*model.Provider, int) {
	results := []*model.Provider{}

	query := fmt.Sprintf(`
		FOR d IN %s
			FILTER d.lineAtChannelSecret != "" AND
				d.lineAtAccessToken != ""
			RETURN d`, collectionProviders)
	bindVars := map[string]interface{}{}
	cursor, err := db.Query(ctx, query, bindVars)
	defer closeCursor(cursor)
	if err != nil {
		logger.Errorf("[ArangoDB][GetProviders] failed to query: %v", err)
		return []*model.Provider{}, constant.ArangoDB_Driver_Failed
	}

	for {
		var doc *model.Provider
		_, err := cursor.ReadDocument(ctx, &doc)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			logger.Errorf("[ArangoDB][GetProviders] failed to read doc: %v", err)
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
		RETURN p`, collectionProviders)
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

	u4 := uuid.New()
	genSharableCode := u4.String()

	aql := `
function (Params) {
	const db = require('@arangodb').db;
	const providerCol = db._collection("Providers");
	let savedObj = Params[0];

	const providerDoc = providerCol.firstExample({lineUserId: savedObj.lineUserId});
	const inviterId = "";
	if (savedObj.inviteCode && savedObj.inviteCode.length > 0) {
		inviterDoc = providerCol.firstExample({inviteCode: savedObj.inviteCode});
		inviterId = (inviterDoc && inviterDoc._key) ? inviterDoc._key : "";
	}
	if (providerDoc.gmailAddr != savedObj.gmailAddr && savedObj.gToken === "") {
		savedObj.gCalSync = false;
		savedObj.guuid = "";
		savedObj.gToken = "";
		savedObj.gRawData = "";
	}
	if (providerDoc && providerDoc._key && providerDoc._key.length > 0 &&
		(providerDoc.status === 0 || providerDoc.status === 1)) {
		providerDoc = {
			...providerDoc,
			...savedObj
			inviterId,
			sharableCode: Params[1],
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
		WriteCollections:   []string{collectionProviders},
		ReadCollections:    []string{collectionProviders},
		Params:             []interface{}{savedMap, genSharableCode},
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
		WriteCollections:   []string{collectionProviders},
		ReadCollections:    []string{collectionProviders},
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

func GetPaymentsByLineUserID(ctx context.Context, lineUserID string, status int) ([]model.Payment, error) {
	results := []model.Payment{}

	cursor, err := db.Query(ctx, `
		LET providerId = (
			FOR p IN @@CollectionP
				FILTER p.lineUserId == @valueId
			RETURN p._key
		)
		FOR p IN @@collectionPayments
			FILTER p.providerId == providerId AND
				p.status == @valueStatus
		RETURN p`,
		map[string]interface{}{
			"@CollectionP":        collectionProviders,
			"@collectionPayments": collectionPayments,
			"valueId":             lineUserID,
			"valueStatus":         status,
		})

	defer closeCursor(cursor)
	if err != nil {
		logger.Errorf("[ArangoDB][GetPaymentsByLineUserID] failed to query: %v", err)
		return results, err
	}

	for {
		var doc model.Payment
		_, err := cursor.ReadDocument(ctx, &doc)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			logger.Errorf("[ArangoDB][GetPaymentsByLineUserID] failed to read doc: %v", err)
			return results, err
		}

		results = append(results, doc)
	}

	return results, nil
}

func GetMonthReceiptByLineUserID(ctx context.Context, lineUserID string) ([]model.MonthReceipt, error) {
	results := []model.MonthReceipt{}
	now := time.Now()
	lastDayOfTwoMonthAgo := now.AddDate(0, -2, -now.Day())

	cursor, err := db.Query(ctx, `
		LET providerId = (
			FOR p IN @@CollectionP
				FILTER p.lineUserId == @valueId
			RETURN p._key
		)
		FOR m IN @@CollectionM
			FILTER m.providerId == providerId AND
				m.clearingStartedAt >= @valueStart
		RETURN m`,
		map[string]interface{}{
			"@CollectionP": collectionProviders,
			"@CollectionM": collectionMonthReceipts,
			"valueId":      lineUserID,
			"valueStart":   lastDayOfTwoMonthAgo,
		})

	defer closeCursor(cursor)
	if err != nil {
		logger.Errorf("[ArangoDB][GetMonthReceiptByLineUserID] failed to query: %v", err)
		return results, err
	}

	for {
		var doc model.MonthReceipt
		_, err := cursor.ReadDocument(ctx, &doc)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			logger.Errorf("[ArangoDB][GetMonthReceiptByLineUserID] failed to read doc: %v", err)
			return results, err
		}

		results = append(results, doc)
	}

	return results, nil
}

func GetScheduleByLineUserID(ctx context.Context, lineUserID string) ([]model.ScheduleCourse, error) {
	results := []model.ScheduleCourse{}
	weekAgoFromNow := time.Now().AddDate(0, 0, -7)
	// weekAfterFromNow := time.Now().AddDate(0, 0, 7)
	weekAgoTimestamp := int64(time.Nanosecond) * weekAgoFromNow.UnixNano() / int64(time.Millisecond)
	// weekAfterTimestamp := int64(time.Nanosecond) * weekAfterFromNow.UnixNano() / int64(time.Millisecond)

	cursor, err := db.Query(ctx, `
		LET providerId = (
			FOR p IN @@CollectionP
				FILTER p.lineUserId == @valueId
			RETURN p._key
		)
		FOR s IN @@CollectionS
			FILTER s.providerId == providerId AND
				s.CourseStartAt >= @valueStart
			LET c = (
				FOR c IN @@CollectionC
					FILTER c._key == s.courseId
				RETURN c
			)
		RETURN MERGE(s, {title: c[0].title, content: c[0].content})`,
		map[string]interface{}{
			"@CollectionP": collectionProviders,
			"@CollectionS": collectionSchedules,
			"@CollectionC": collectionCourses,
			"valueId":      lineUserID,
			"valueStart":   weekAgoTimestamp,
		})

	defer closeCursor(cursor)
	if err != nil {
		logger.Errorf("[ArangoDB][GetScheduleByLineUserID] failed to query: %v", err)
		return results, err
	}

	for {
		var doc model.ScheduleCourse
		_, err := cursor.ReadDocument(ctx, &doc)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			logger.Errorf("[ArangoDB][GetScheduleByLineUserID] failed to read doc: %v", err)
			return results, err
		}

		results = append(results, doc)
	}

	return results, nil
}

func GetAppointmentsByLineUserID(ctx context.Context, lineUserID string) ([]model.Appointment, error) {
	results := []model.Appointment{}
	weekAgoFromNow := time.Now().AddDate(0, 0, -7)
	// weekAfterFromNow := time.Now().AddDate(0, 0, 7)
	weekAgoTimestamp := int64(time.Nanosecond) * weekAgoFromNow.UnixNano() / int64(time.Millisecond)
	// weekAfterTimestamp := int64(time.Nanosecond) * weekAfterFromNow.UnixNano() / int64(time.Millisecond)

	cursor, err := db.Query(ctx, `
		LET providerId = (
			FOR p IN @@CollectionP
				FILTER p.lineUserId == @valueId
			RETURN p._key
		)
		FOR a IN @@CollectionA
			FILTER a.providerId == providerId AND
				a.CourseStartAt >= @valueStart AND
				a.Status != @valueStatus
		RETURN a`,
		map[string]interface{}{
			"@CollectionP": collectionProviders,
			"@CollectionA": collectionAppointments,
			"valueId":      lineUserID,
			"valueStart":   weekAgoTimestamp,
			"valueStatus":  -1,
		})

	defer closeCursor(cursor)
	if err != nil {
		logger.Errorf("[ArangoDB][GetAppointmentsByLineUserID] failed to query: %v", err)
		return results, err
	}

	for {
		var doc model.Appointment
		_, err := cursor.ReadDocument(ctx, &doc)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			logger.Errorf("[ArangoDB][GetAppointmentsByLineUserID] failed to read doc: %v", err)
			return results, err
		}

		results = append(results, doc)
	}

	return results, nil
}

func GetMonthReceiptList(ctx context.Context, lineUserID string) ([]model.MonthReceipt, error) {
	results := []model.MonthReceipt{}

	cursor, err := db.Query(ctx, `
		LET providerId = (
			FOR p IN @@CollectionP
				FILTER p.lineUserId == @valueId
			RETURN p._key
		)
		FOR m IN @@collectionMonthReceipts
			FILTER m.providerId == providerId
		RETURN m`,
		map[string]interface{}{
			"@CollectionP":             collectionProviders,
			"@collectionMonthReceipts": collectionMonthReceipts,
			"valueId":                  lineUserID,
		})

	defer closeCursor(cursor)
	if err != nil {
		logger.Errorf("[ArangoDB][GetMonthReceiptList] failed to query: %v", err)
		return results, err
	}

	for {
		var doc model.MonthReceipt
		_, err := cursor.ReadDocument(ctx, &doc)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			logger.Errorf("[ArangoDB][GetMonthReceiptList] failed to read doc: %v", err)
			return results, err
		}

		results = append(results, doc)
	}

	return results, nil
}

func GetServiceProductsByLineUserID(ctx context.Context, lineUserID string) ([]model.ServiceProduct, error) {
	results := []model.ServiceProduct{}

	cursor, err := db.Query(ctx, `
		LET providerId = (
			FOR p IN @@CollectionP
				FILTER p.lineUserId == @valueId
			RETURN p._key
		)
		FOR s IN @@collectionServiceProducts
			FILTER s.providerId == providerId AND
			(IS_NULL(s.deletedAt) OR s.deletedAt == 0)
		RETURN s`,
		map[string]interface{}{
			"@CollectionP":               collectionProviders,
			"@collectionServiceProducts": collectionServiceProducts,
			"valueId":                    lineUserID,
		})

	defer closeCursor(cursor)
	if err != nil {
		logger.Errorf("[ArangoDB][GetServiceProductsByLineUserID] failed to query: %v", err)
		return results, err
	}

	for {
		var doc model.ServiceProduct
		_, err := cursor.ReadDocument(ctx, &doc)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			logger.Errorf("[ArangoDB][GetServiceProductsByLineUserID] failed to read doc: %v", err)
			return results, err
		}

		results = append(results, doc)
	}

	return results, nil
}

func CreateOrUpdateServiceProduct(ctx context.Context, lineUserID string,
	params model.CreateOrUpdateServiceProductsReq) error {

	aql := `
function (Params) {
	const db = require('@arangodb').db;
	const providerCol = db._collection("Providers");
	const svcProductsCol = db._collection("ServiceProducts");
	const lineUserId = Params[0];
	const savedObjList = Params[1];

	providerDoc = providerCol.firstExample({lineUserId: lineUserId});
	if (!(providerDoc && providerDoc._key && providerDoc._key.length > 0 &&
		(providerDoc.status >= 3))) {
		return 2;
	}

	if (savedObjList.length == 0) {
		return 1;
	}

	savedObjList.forEach(function(obj) {
		const keyExist = obj._key && obj._key.length > 0;
		if (keyExist) {
			svcProductsCol.update({ _key: obj._key }, obj);
		} else {
			svcProductsCol.save(obj);
		}
	})

	return 1;
}
`
	options := &driver.TransactionOptions{
		MaxTransactionSize: 100000,
		WriteCollections:   []string{collectionProviders, collectionServiceProducts},
		ReadCollections:    []string{collectionProviders, collectionServiceProducts},
		Params:             []interface{}{lineUserID, params.ReqList},
		WaitForSync:        false,
	}

	resCode, err := db.Transaction(ctx, aql, options)
	if err != nil {
		logger.Errorf("CreateOrUpdateServiceProduct TX execution failure")
		return err
	}
	if resCode == 2 {
		return errors.New("Not found document")
	}

	logger.Debugf("CreateOrUpdateServiceProduct TX execution resCode is : %v", resCode)

	return nil
}

func GetScheduleList(ctx context.Context, lineUserID string, start int64) ([]model.Schedule, error) {
	results := []model.Schedule{}

	cursor, err := db.Query(ctx, `
		LET providerId = (
			FOR p IN @@CollectionP
				FILTER p.lineUserId == @valueId
			RETURN p._key
		)
		FOR s IN @@collectionSchedules
			FILTER s.providerId == providerId AND
			s.courseStartAt >= @valueStart
		RETURN s`,
		map[string]interface{}{
			"@CollectionP":         collectionProviders,
			"@collectionSchedules": collectionSchedules,
			"valueId":              lineUserID,
			"valueStart":           start,
		})

	defer closeCursor(cursor)
	if err != nil {
		logger.Errorf("[ArangoDB][GetScheduleList] failed to query: %v", err)
		return results, err
	}

	for {
		var doc model.Schedule
		_, err := cursor.ReadDocument(ctx, &doc)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			logger.Errorf("[ArangoDB][GetScheduleList] failed to read doc: %v", err)
			return results, err
		}

		results = append(results, doc)
	}

	return results, nil
}

func CreateProviderSchedule(ctx context.Context, lineUserID string,
	params model.CreateServiceScheduleReq) error {

	aql := `
function (Params) {
	const db = require('@arangodb').db;
	const providerCol = db._collection("Providers");
	const courseCol = db._collection("Courses");
	const scheduleCol = db._collection("Schedules");
	const lineId = Params[0];
	const savedObj = Params[1];

	const courseRes = courseCol.save({
		providerId: lineId,
		title: savedObj.title,
		content: savedObj.content
	});

	const scheduleRes = scheduleCol.save({
		courseId: courseRes._key,
		providerId: lineId,
		cousreStartAt: savedObj.cousreStartAt,
		courseEndAt: savedObj.courseEndAt,
		minConsumerLimit: savedObj.minConsumerLimit,
		maxConsumerLimit: savedObj.maxConsumerLimit,
		count: savedObj.count
	});

	courseCol.update({_key: courseRes._key}, {scheduleId: scheduleRes._key})

	return 1;
}
`
	options := &driver.TransactionOptions{
		MaxTransactionSize: 100000,
		WriteCollections:   []string{collectionProviders, collectionSchedules, collectionCourses},
		ReadCollections:    []string{collectionProviders, collectionSchedules, collectionCourses},
		Params:             []interface{}{lineUserID, params},
		WaitForSync:        false,
	}

	resCode, err := db.Transaction(ctx, aql, options)
	if err != nil {
		logger.Errorf("CreateProviderSchedule TX execution failure")
		return err
	}

	logger.Debugf("CreateProviderSchedule TX execution resCode is : %v", resCode)

	return nil
}

/**
 * Private methods
 */
func formatCreateProviderMap(registerInfo model.RegisterProviderReq) (map[string]interface{}, error) {
	dataMap := map[string]interface{}{}
	var now int64
	if registerInfo.Status == constant.Provider_Status_Saved {
		now = time.Now().Unix()
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
		InviteCode:          registerInfo.InviteCode,
		SharableCode:        "",
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
		GUUID:              providerInfo.GUUID,
		GToken:             providerInfo.GToken,
		GRawData:           providerInfo.GRawData,
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
	delete(dataMap, "inviteCode")
	delete(dataMap, "inviterId")
	delete(dataMap, "sharableCode")
	if providerInfo.Status != constant.Provider_Status_Auditing {
		delete(dataMap, "memeberTerm")
		delete(dataMap, "privacyTerm")
	}
	delete(dataMap, "createdAt")
	if providerInfo.GToken == "" {
		delete(dataMap, "guuid")
		delete(dataMap, "gToken")
		delete(dataMap, "gRawData")
	}

	return dataMap, nil
}
