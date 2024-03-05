package xutils

import "github.com/pquerna/otp/totp"

// 验证谷歌验证码
func VerifyGoogleCode(secret string, code string) bool {
	return totp.Validate(code, secret)
}

// 创建新的google secret
func NewGoogleSecret(Issuer string, AccountName string) (string, string) {
	key, _ := totp.Generate(totp.GenerateOpts{
		Issuer:      Issuer,
		AccountName: AccountName,
	})
	return key.Secret(), key.URL()
}

// 获取谷歌二维码
func GetGoogleQrCodeUrl(secret string, issuer string, accountname string) string {
	key, _ := totp.Generate(totp.GenerateOpts{
		Issuer:      issuer,
		AccountName: accountname,
		Secret:      []byte(secret),
	})
	return key.URL()
}
