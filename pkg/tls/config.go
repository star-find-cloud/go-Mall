package tls

import (
	"crypto/tls"
	"errors"
	"github.com/star-find-cloud/star-mall/conf"
	applog "github.com/star-find-cloud/star-mall/pkg/logger"
)

// 获取 TLS 证书配置
func GetTLS() (*tls.Config, error) {
	var c = conf.GetConfig()
	cert, err := tls.LoadX509KeyPair(c.TLS.Server.CertFile, c.TLS.Server.KeyFile)
	if err != nil {
		applog.AppLogger.Errorln("TLS .crt .key 文件加载失败")
		return nil, errors.New("TLS .crt .key 文件加载失败")
	}
	if c.TLS.Server.CertFile == "" {
		applog.AppLogger.Errorln("配置文件中TLS路径加载失败")
		return nil, errors.New("配置文件中TLS路径加载失败")
	}

	return &tls.Config{
		Certificates:           []tls.Certificate{cert},
		ClientAuth:             tls.ClientAuthType(c.TLS.Server.ClientAuth),
		MinVersion:             parseTLSVersion(c.TLS.Server.MinVersion),
		CipherSuites:           parseTLSCipherSuites(c.TLS.Cipher_suites.Suites),
		SessionTicketsDisabled: c.TLS.Advanced.SessionTickets,
		CurvePreferences:       parseTLSCurvePreferences(c.TLS.Advanced.CurvePreferences),
	}, nil
}
