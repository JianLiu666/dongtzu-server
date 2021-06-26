package server

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"dongtzu/config"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gitlab.geax.io/demeter/gologger/logger"
)

func hookNewebpay() fiber.Handler {
	return func(c *fiber.Ctx) error {
		logger.Debugf("Receive notify. %v\n%v\n", c.Method(), string(c.Body()))

		return c.JSON(fiber.Map{"code": 200})
	}
}

func getNewebPayOrderInfo() fiber.Handler {
	return func(c *fiber.Ctx) error {
		timeNow := time.Now().Unix()
		merchantOrderNo := fmt.Sprintf("%v", timeNow)
		amount := 500
		courseName := "TestCourse"
		email := "jianliu0616@gmail.com"

		tradeInfo := fmt.Sprintf("MerchantID=%v&RespondType=%v&TimeStamp=%v&Version=%v&MerchantOrderNo=%v&Amt=%v&ItemDesc=%v&TradeLimit=%v&Email=%v&LoginType=%v",
			config.GetGlobalConfig().NewebPay.MerchantID,
			"JSON",
			timeNow,
			config.GetGlobalConfig().NewebPay.APIVersion,
			merchantOrderNo,
			amount,
			courseName,
			600,
			email,
			0,
		)

		tradeInfo, tradeSHA := genereateTradeInfoAndTradeSHA(tradeInfo)

		resp, err := http.PostForm(config.GetGlobalConfig().NewebPay.APIUrl, url.Values{
			"MerchantID": []string{config.GetGlobalConfig().NewebPay.MerchantID},
			"TradeInfo":  []string{tradeInfo},
			"TradeSha":   []string{tradeSHA},
			"Version":    []string{config.GetGlobalConfig().NewebPay.APIVersion},
		})
		if err != nil {
			logger.Debugf("[Server][NewebPay] failed to post form: %v", err)
		}
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)

		return c.SendString(string(body))
	}
}

func genereateTradeInfoAndTradeSHA(rawData string) (string, string) {
	key := config.GetGlobalConfig().NewebPay.MerchantHashKey
	iv := config.GetGlobalConfig().NewebPay.MerchantIV

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		logger.Debugf("new cipher error: %v", err)
	}
	origData := usePKCS7Padding([]byte(rawData), 32)
	blockMode := cipher.NewCBCEncrypter(block, []byte(iv))
	encrypted := make([]byte, len(origData))
	blockMode.CryptBlocks(encrypted, origData)

	shaString := fmt.Sprintf("HashKey=%v&%v&HashIV=%v",
		key,
		hex.EncodeToString(encrypted),
		iv,
	)
	sum := sha256.Sum256([]byte(shaString))

	return hex.EncodeToString(encrypted), strings.ToUpper(fmt.Sprintf("%x", sum[:]))
}

func usePKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)

	return append(ciphertext, padtext...)
}
