package conf

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
	"time"
)

type Config struct {
	App      AppConfig
	Database DatabaseConf
	Log      LogConf
	Mail     MailConf
	Cookie   CookieConf
	Etcd     EtcdConf
	TLS      TLSConfig
}

type AppConfig struct {
	Name                 string `mapstructure:"name"`
	Port                 uint64 `mapstructure:"port"`
	Env                  string `mapstructure:"env"`
	Debug                bool   `mapstructure:"debug"`
	Version              string `mapstructure:"version"`
	Sessiongcmaxlifetime int    `mapstructure:"sessiongcmaxlifetime"`
	SessionName          string `mapstructure:"session_name"`
}

type DatabaseConf struct {
	MySQL MySQLConf `mapstructure:"mysql"`
	Redis RedisConf `mapstructure:"redis"`
}

type MySQLConf struct {
	MasterHost   string       `mapstructure:"master_host"`
	MasterPort   uint64       `mapstructure:"master_port"`
	User         string       `mapstructure:"user"`
	Password     string       `mapstructure:"password"`
	DBName       string       `mapstructure:"db_name"`
	MaxOpenConns int          `mapstructure:"max_open_conns"`
	MaxIdleConns int          `mapstructure:"max_idle_conns"`
	Timeout      string       `mapstructure:"timeout"`
	Slaves       []MySQLSlave `mapstructure:"slaves"`
}

type MySQLSlave struct {
	Host string `mapstructure:"host"`
	Port uint64 `mapstructure:"port"`
}

type RedisConf struct {
	MasterHost string       `mapstructure:"master_host"`
	MasterPort uint64       `mapstructure:"master_port"`
	Password   string       `mapstructure:"password"`
	Slaves     []RedisSlave `mapstructure:"slaves"`
}

type RedisSlave struct {
	Host string `mapstructure:"host"`
	Port uint64 `mapstructure:"port"`
}

type LogConf struct {
	Version   string `mapstructure:"version"`
	Level     string `mapstructure:"level"`
	Path      string `mapstructure:"path"`
	MaxSize   int    `mapstructure:"max_size"`
	MaxBackup int    `mapstructure:"max_backup"`
	MaxAge    int    `mapstructure:"max_age"`
	Compress  bool   `mapstructure:"compress"`
}

type MailConf struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	SMTPCode string `mapstructure:"SMTP-code"`
	From     string `mapstructure:"from"`
}

type CookieConf struct {
	Name     string
	Domain   string
	Path     string
	MaxAge   int
	Secure   bool
	HttpOnly bool
	SameSite string
}

type EtcdConf struct {
	Endpoints []string `mapstructure:"endpoints"`
	//Username           string   `mapstructure:"username"`
	//Password           string   `mapstructure:"password"`
	DialTimeout        time.Duration `mapstructure:"dial_timeout"`           // e.g. "5s"
	AutoSyncInterval   time.Duration `mapstructure:"auto_sync_interval"`     // e.g. "1m"
	MaxCallSendMsgSize int           `mapstructure:"max_call_send_msg_size"` // 单位: bytes
	MaxCallRecvMsgSize int           `mapstructure:"max_call_recv_msg_size"` // 单位: bytes
	EnableTLS          bool          `mapstructure:"enable_tls"`
	TLS                TLSConfig
}

type TLSConfig struct {
	Server        ServerConfig
	Client        ClientConfig
	Cipher_suites CipherSuitesConfig
	Advanced      AdvancedConfig
}

type ServerConfig struct {
	CertFile    string `mapstructure:"cert_file"`    // PEM格式服务端证书文件路径（与PfxFile互斥）
	KeyFile     string `mapstructure:"key_file"`     // PEM格式服务端私钥文件路径
	ClientAuth  int    `mapstructure:"client_auth"`  // 是否启用客户端证书验证（双向认证）
	MinVersion  string `mapstructure:"min_version"`  // 最低支持的TLS协议版本（如"TLS1.2"）
	PfxFile     string `mapstructure:"pfx_file"`     // PKCS12格式证书文件路径（替代CertFile+KeyFile）
	PfxPassword string `mapstructure:"pfx_password"` // PKCS12文件的解密密码
}

type ClientConfig struct {
	CaCertFile string `mapstructure:"ca_cert_file"` // CA证书路径，用于验证服务端证书合法性
}

type CipherSuitesConfig struct {
	Suites []string `mapstructure:"suites"` // 启用的加密套件列表（如TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256）
}

type AdvancedConfig struct {
	SessionTickets   bool     `mapstructure:"session_tickets"`   // 是否启用会话票据复用以提升性能
	CurvePreferences []string `mapstructure:"curve_preferences"` // 椭圆曲线优先级（如"X25519"、"P-256"）
}

var c Config

func init() {
	viper.SetConfigName("StarMall")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")

	logFile, err := os.OpenFile("E:/var/log/StarMall/conf.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Println(err)
	}
	defer logFile.Close()

	log.SetOutput(logFile)
	// 读取配置文件
	err = viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
		log.Printf("config load Error: %v \n", err)
	} else {
		log.Println("configuration file was read successfully")
	}

	// 将 viper 读到的数据序列化写入 config
	if err := viper.Unmarshal(&c); err != nil {
		now := time.Now()
		log.Printf("%v: viper Unmarshal err:%s \n", now.Format("2006-01-02 15:04:05"), err)
	}
}

func GetConfig() *Config {
	if c.App.Name == "" {
		fmt.Println("配置文件读取失败")
		logFile, err := os.OpenFile("E:/var/log/StarMall/conf.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Println(err)
		}
		defer logFile.Close()

		log.SetOutput(logFile)
		log.Printf("config load Error: %v \n", err)
	}
	return &c
}
