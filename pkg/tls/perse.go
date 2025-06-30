package tls

import (
	"crypto/tls"
	"strings"
)

// 将配置文件中的配置转为符合 tls 包要求的配置
// 将版本转换为 tls 包要求
func parseTLSVersion(version string) uint16 {
	switch strings.ToUpper(strings.TrimSpace(version)) {
	case "TLS1.3", "TLSV1.3":
		return tls.VersionTLS13
	case "TLS1.2", "TLSV1.2":
		return tls.VersionTLS12
	case "TLS1.1", "TLSV1.1":
		return tls.VersionTLS11
	case "TLS1.0", "TLSV1.0":
		return tls.VersionTLS10
	default:
		// 默认版本
		return tls.VersionTLS12
	}
}

// 将加密套件转换为 tls 包要求
func parseTLSCipherSuites(names []string) []uint16 {
	var cipherSuites []uint16
	for _, name := range names {
		switch strings.ToUpper(strings.TrimSpace(name)) {
		case "TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384":
			cipherSuites = append(cipherSuites, tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384)
		case "TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384":
			cipherSuites = append(cipherSuites, tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384)
		}
	}
	//fmt.Println(cipherSuites)
	return cipherSuites
}

// 将加密椭圆曲线优先级转换为 tls 包要求
func parseTLSCurvePreferences(CurvePreferences []string) []tls.CurveID {
	var curvePreferencesSlice []tls.CurveID
	for _, name := range CurvePreferences {
		switch strings.ToUpper(strings.TrimSpace(name)) {
		case "X25519":
			curvePreferencesSlice = append(curvePreferencesSlice, tls.X25519)
		case "P256":
			curvePreferencesSlice = append(curvePreferencesSlice, tls.CurveP256)
		case "P384":
			curvePreferencesSlice = append(curvePreferencesSlice, tls.CurveP384)
		case "P521":
			curvePreferencesSlice = append(curvePreferencesSlice, tls.CurveP521)
		}
	}
	//fmt.Println(curvePreferencesSlice)
	return curvePreferencesSlice
}
