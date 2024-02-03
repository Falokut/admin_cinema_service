package config

import (
	"crypto/tls"
	"crypto/x509"
	"sync"

	"github.com/Falokut/admin_cinema_service/internal/repository"
	"github.com/Falokut/admin_cinema_service/pkg/jaeger"
	"github.com/Falokut/admin_cinema_service/pkg/metrics"
	logging "github.com/Falokut/online_cinema_ticket_office.loggerwrapper"
	"github.com/ilyakaznacheev/cleanenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type DialMethod = string

const (
	Insecure                 DialMethod = "INSECURE"
	NilTlsConfig             DialMethod = "NIL_TLS_CONFIG"
	ClientWithSystemCertPool DialMethod = "CLIENT_WITH_SYSTEM_CERT_POOL"
	Server                   DialMethod = "SERVER"
)

type ConnectionSecureConfig struct {
	Method DialMethod `yaml:"dial_method"`
	// Only for client connection with system pool
	ServerName string `yaml:"server_name"`
	CertName   string `yaml:"cert_name"`
	KeyName    string `yaml:"key_name"`
}

type Config struct {
	LogLevel        string `yaml:"log_level" env:"LOG_LEVEL"`
	HealthcheckPort string `yaml:"healthcheck_port" env:"HEALTHCHECK_PORT"`
	Listen          struct {
		Host string `yaml:"host" env:"HOST"`
		Port string `yaml:"port" env:"PORT"`
		Mode string `yaml:"server_mode" env:"SERVER_MODE"` // support GRPC, REST, BOTH
	} `yaml:"listen"`

	PrometheusConfig struct {
		Name         string                      `yaml:"service_name" env:"PROMETHEUS_SERVICE_NAME"`
		ServerConfig metrics.MetricsServerConfig `yaml:"server_config"`
	} `yaml:"prometheus"`
	MoviesService struct {
		Addr             string                 `yaml:"addr" env:"MOVIES_SERVICE_ADDRESS"`
		ConnectionConfig ConnectionSecureConfig `yaml:"connection_config"`
	} `yaml:"movies_service"`

	DBConfig     repository.DBConfig `yaml:"db_config"`
	JaegerConfig jaeger.Config       `yaml:"jaeger"`
}

var instance *Config
var once sync.Once

const configsPath = "configs/"

func GetConfig() *Config {
	once.Do(func() {
		logger := logging.GetLogger()
		instance = &Config{}

		if err := cleanenv.ReadConfig(configsPath+"config.yml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Fatal(help, " ", err)
		}
	})

	return instance
}

func (c ConnectionSecureConfig) GetGrpcTransportCredentials() (grpc.DialOption, error) {
	if c.Method == Insecure {
		return grpc.WithTransportCredentials(insecure.NewCredentials()), nil
	}

	if c.Method == NilTlsConfig {
		return grpc.WithTransportCredentials(credentials.NewTLS(nil)), nil
	}

	if c.Method == ClientWithSystemCertPool {
		certPool, err := x509.SystemCertPool()
		if err != nil {
			return grpc.EmptyDialOption{}, err
		}
		return grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(certPool, c.ServerName)), nil
	}

	cert, err := tls.LoadX509KeyPair(c.CertName, c.KeyName)
	if err != nil {
		return grpc.EmptyDialOption{}, err
	}
	return grpc.WithTransportCredentials(credentials.NewServerTLSFromCert(&cert)), nil
}
