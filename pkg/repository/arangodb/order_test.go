package arangodb

import (
	"context"
	"testing"

	"gitlab.geax.io/demeter/gologger/logger"
)

func TestCreateOrder(t *testing.T) {
	initConfig()

	order, status := CreateOrder(context.TODO(), "500232", "499903", "501555", 1)
	logger.Debugf("%v", status)
	logger.Debugf("%v", order)
}
