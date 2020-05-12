package encrypt

import "testing"

func TestAesEncryptString(t *testing.T) {
	str := "123456789"
	hex, err := AesEncrypt(str)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(hex)

	dec, err := AesDecrypt(hex)
	if err != nil {
		t.Fatal(err)
	}
	if dec != str {
		t.Fatal()
	}
}
