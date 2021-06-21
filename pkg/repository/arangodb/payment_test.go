package arangodb

import (
	"context"
	"dongtzu/pkg/model"
	"testing"
	"time"

	"gitlab.geax.io/demeter/gologger/logger"
)

func TestCreatePayment(t *testing.T) {
	initConfig()

	doc := &model.Payment{
		OrderID:         "503996",
		ConsumerID:      "500232",
		PaymentMethodID: "501713",
		PaidPrice:       1000,
		PlatformFee:     0,
		PaymentFee:      0,
		AgentFee:        0,
		AdFee:           0,
		TaxFee:          0,
		NetAmount:       0,
		Status:          0,
		RawParams:       "",
		CreatedAt:       time.Now().Unix(),
		UpdatedAt:       time.Now().Unix(),
	}
	status := CreatePayment(context.TODO(), doc)
	logger.Debugf("%v", status)
}
