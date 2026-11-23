package lea

// 对比密码是否一致
// 因为之前密码是用md5加密的, 所以通过密码长度来判断
// rawPwd 原始, 用户输入的密码
func ComparePwd(rawPwd, dbPwd string) bool {
	if len(dbPwd) == 32 {
		return Md5(rawPwd) == dbPwd
	}

	return CompareHash([]byte(dbPwd), rawPwd)
}

// NeedsRehash reports a legacy MD5 password that must be upgraded to bcrypt
// after the next successful login.
func NeedsRehash(dbPwd string) bool {
	return len(dbPwd) == 32
}

// PasswordMinLength is the minimum accepted password length for new passwords.
const PasswordMinLength = 8

// 加密
func GenPwd(rawPwd string) string {
	digest, err := GenerateHash(rawPwd)
	if err != nil {
		return ""
	}
	return string(digest)
}
