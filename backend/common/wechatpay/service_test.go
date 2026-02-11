package wechatpay

import (
	"strings"
	"testing"
)

func newTestService() *Service {
	return NewService(&Config{
		AppID:  "wx-app",
		MchID:  "mch-1",
		APIKey: "key123",
	})
}

func TestGenerateAndVerifySign(t *testing.T) {
	s := newTestService()
	params := map[string]string{
		"appid":        "wx-app",
		"mch_id":       "mch-1",
		"nonce_str":    "abc",
		"out_trade_no": "o-1",
		"total_fee":    "100",
	}
	sign := s.GenerateMD5Sign(params)
	if sign == "" {
		t.Fatalf("empty sign")
	}
	if !s.VerifyMD5Sign(params, sign) {
		t.Fatalf("verify sign failed")
	}
	if s.VerifyMD5Sign(params, "BAD") {
		t.Fatalf("verify should fail")
	}
}

func TestCreateNativePay(t *testing.T) {
	s := newTestService()
	resp, err := s.CreateNativePay("order-1", 520, "body")
	if err != nil {
		t.Fatalf("create native pay error: %v", err)
	}
	if resp == nil || resp.ReturnCode != "SUCCESS" {
		t.Fatalf("unexpected response")
	}
	if !strings.Contains(resp.CodeURL, "weixin://wxpay") {
		t.Fatalf("unexpected code_url: %s", resp.CodeURL)
	}
}

func TestParseNotifyRequest(t *testing.T) {
	s := newTestService()
	params := map[string]string{
		"return_code":    "SUCCESS",
		"appid":          "wx-app",
		"mch_id":         "mch-1",
		"device_info":    "d",
		"nonce_str":      "n",
		"out_trade_no":   "o-1",
		"transaction_id": "tx-1",
		"total_fee":      "100",
		"fee_type":       "CNY",
		"time_end":       "20260101010101",
		"attach":         "a",
		"is_subscribe":   "Y",
	}
	sign := s.GenerateMD5Sign(params)
	xml := "<xml>" +
		"<return_code>SUCCESS</return_code>" +
		"<appid>wx-app</appid>" +
		"<mch_id>mch-1</mch_id>" +
		"<device_info>d</device_info>" +
		"<nonce_str>n</nonce_str>" +
		"<sign>" + sign + "</sign>" +
		"<out_trade_no>o-1</out_trade_no>" +
		"<transaction_id>tx-1</transaction_id>" +
		"<total_fee>100</total_fee>" +
		"<fee_type>CNY</fee_type>" +
		"<time_end>20260101010101</time_end>" +
		"<attach>a</attach>" +
		"<is_subscribe>Y</is_subscribe>" +
		"</xml>"

	notify, err := s.ParseNotifyRequest(xml)
	if err != nil {
		t.Fatalf("parse notify error: %v", err)
	}
	if notify == nil || notify.OutTradeNo != "o-1" {
		t.Fatalf("unexpected notify")
	}
}

func TestBuildNotifyResponseAndNonce(t *testing.T) {
	x := BuildNotifyResponse()
	if x != "" {
		t.Fatalf("unexpected notify response: %s", x)
	}
	n1 := generateNonce()
	n2 := generateNonce()
	if n1 == "" || n2 == "" || n1 == n2 {
		t.Fatalf("nonce invalid")
	}
}
