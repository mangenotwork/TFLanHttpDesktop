package utils

import (
	"testing"
)

func Test_Sign(t *testing.T) {
	data := "hello world"

	for i := 0; i < 100; i++ {
		// 生成签名
		signature, err := GenerateSignature(data)
		if err != nil {
			t.Logf("生成签名失败: %v\n", err)
			return
		}
		t.Logf("生成的签名: %s\n", signature)

		// 校验签名（正常情况）
		valid, err := VerifySignature(data, signature)
		if err != nil {
			t.Logf("校验失败: %v\n", err)
		} else {
			t.Logf("签名是否有效: %v\n", valid)
		}

		// 校验被篡改的数据（测试用）
		invalidData := "hello-world-tampered"
		valid, err = VerifySignature(invalidData, signature)
		if err != nil {
			t.Logf("篡改数据校验失败: %v\n", err)
		} else {
			t.Logf("篡改数据签名是否有效: %v\n", valid)
		}
	}

}
