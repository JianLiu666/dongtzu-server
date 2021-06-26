package server

import (
	"dongtzu/config"
	"fmt"

	"github.com/spf13/viper"
	"gitlab.geax.io/demeter/gologger/logger"
)

func initConfig() {
	logger.Init("debug")

	viper.SetConfigFile("./../../../conf.d/env.yaml")
	if err := viper.ReadInConfig(); err != nil {
		logger.Errorf("ReadInConfig file failed: %v", err)
	}

	c, err := config.NewFromViper()
	if err != nil {
		logger.Errorf("Init config failed: %v", err)
	}
	config.SetConfig(c)
}

func Example_genereateTradeInfoAndTradeSHA() {
	initConfig()

	tradeInfo, tradeSHA := genereateTradeInfoAndTradeSHA("MerchantID=MS120977897&RespondType=JSON&TimeStamp=1624633851&Version=1.5&MerchantOrderNo=1624633851&Amt=49&ItemDesc=TEST&Email=ts0542648%40gmail.com")

	fmt.Println(tradeInfo)
	fmt.Println(tradeSHA)

	// Output:
	// 9caec7216c452a3aa99e0d437c0c4dcc0b374544e2979b4d47831167ea430ae045c8eeb7c1694ebd794326f884e08cfd067202fda2068e463407c0cd1c99af4318b5e5defb509cf4f350f60ac2f88c103e0be7abcf7a0c3999572f6e76d97920525b65629f5087a86c2ad876952c96a6c504fb5d5460b86e4b7e1c1407d768f1dc35241a11791b2eca355ad382cb06ac4cbe8b59516e1ae88958ba2a3ab08a55
	// 1EBF9A2B83F75E6A61FE96552F8A633C3FF663CD3661AA8E7D746439B63985D9
}
