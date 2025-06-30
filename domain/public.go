package domain

import "errors"

func ValidateCode(inputCode, storedCode string) (bool, error) {
	if storedCode == "" {
		return false, errors.New("验证码已过期")
	}

	if inputCode != storedCode {
		return false, errors.New("验证码错误")
	}
	return true, nil
}
