package config

import (
	"log"
	"strings"
	"sync"
	"time"

	"github.com/spf13/viper"
)

var (
	configuration *Configurations
	mutex         sync.Once
)

type Configurations struct {
	App         AppConfig         `mapstructure:"app"`
	GoogleCloud GoogleCloudConfig `mapstructure:"google"`
	PubSub      PubSubConfig      `mapstructure:"pubsub"`
	Consumer    ConsumerConfig    `mapstructure:"consumer"`
	Producer    ProducerConfig    `mapstructure:"producer"`
	API         APIConfig         `mapstructure:"api"`
	Database    DatabaseConfig    `mapstructure:"db"`
	CKG         CKGConfig         `mapstructure:"ckg"`
}

type AppConfig struct {
	Environment string `mapstructure:"env"`
	LogLevel    string `mapstructure:"loglevel"`
}

type GoogleCloudConfig struct {
	ProjectID       string `mapstructure:"project"`
	CredentialsPath string `mapstructure:"credentials"`
	Debug           bool   `mapstructure:"debug"`
}

type PubSubConfig struct {
	Topic           string `mapstructure:"topic"`
	Subscription    string `mapstructure:"subscription"`
	MessageOrdering bool   `mapstructure:"messageordering"`
}

type ConsumerConfig struct {
	MaxMessagesPerPull    int               `mapstructure:"maxmessages"`
	SleepTimeBetweenPulls time.Duration     `mapstructure:"sleeptime"`
	AcknowledgeTimeout    time.Duration     `mapstructure:"acktimeout"`
	RetryCount            int               `mapstructure:"retrycount"`
	RetryDelay            time.Duration     `mapstructure:"retrydelay"`
	FlowControl           FlowControlConfig `mapstructure:"flowcontrol"`
	// DeadLetterPolicy      DeadLetterPolicyConfig `mapstructure:"deadletterpolicy"`
}

/*type DeadLetterPolicyConfig struct {
	Enabled                bool  `mapstructure:"enabled"`
	MaxDeliveryAttempts   int  `mapstructure:"maxdelivery"`
	DeadLetterTopicSuffix string `mapstructure:"topicsuffix"`
}*/

type FlowControlConfig struct {
	Enabled                bool  `mapstructure:"enabled"`
	MaxOutstandingMessages int   `mapstructure:"maxmessages"`
	MaxOutstandingBytes    int64 `mapstructure:"maxbytes"`
}

type ProducerConfig struct {
	EnableMessageOrdering bool              `mapstructure:"enableordering"`
	BatchSize             int               `mapstructure:"batchsize"`
	MessageAttributes     map[string]string `mapstructure:"attributes"`
	Compression           CompressionConfig `mapstructure:"compression"`
}

type CompressionConfig struct {
	Enabled   bool   `mapstructure:"enabled"`
	Algorithm string `mapstructure:"algorithm"`
}

type APIConfig struct {
	BaseURL   string        `mapstructure:"baseurl"`
	Timeout   time.Duration `mapstructure:"timeout"`
	APIKey    string        `mapstructure:"apikey"`
	APIHeader string        `mapstructure:"apiheader"`
	BatchSize int           `mapstructure:"batchsize"`
}

type DatabaseConfig struct {
	Driver     string `mapstructure:"driver"`
	Host       string `mapstructure:"host"`
	Port       int    `mapstructure:"port"`
	Username   string `mapstructure:"username"`
	Password   string `mapstructure:"password"`
	Database   string `mapstructure:"database"`
	Attributes string `mapstructure:"attributes"`
}

type CKGConfig struct {
	UseCache           bool   `mapstructure:"usecache"`
	TableMasterWilayah string `mapstructure:"tablemasterwilayah"`
	TableMasterFaskes  string `mapstructure:"tablemasterfaskes"`
	TableSkrining      string `mapstructure:"tableskrining"`
	TableStatus        string `mapstructure:"tablestatus"`
	TableIncoming      string `mapstructure:"tableincoming"`
	TableOutgoing      string `mapstructure:"tableoutgoing"`
	MarkerField        string `mapstructure:"markerfield"`
	MarkerConsume      string `mapstructure:"markerconsume"`
	MarkerProduce      string `mapstructure:"markerproduce"`
}

func GetConfig() *Configurations {
	mutex.Do(func() {
		configuration = newConfig()
	})

	return configuration
}

func newConfig() *Configurations {
	var cfg Configurations

	// Inisialisasi Viper
	v := viper.New()

	// Load environment variables dari file .env jika ada
	v.SetConfigName(".env")
	v.SetConfigType("env")
	v.AddConfigPath(".") // direktori saat ini
	if err := v.ReadInConfig(); err == nil {
		log.Printf("Muat konfigurasi dari file .env")
	}

	log.Printf("Muat konfigurasi dari environment variables")

	// Enable automatic environment variable support
	v.AutomaticEnv()

	// Replace dots with underscores for environment variable keys
	replacer := strings.NewReplacer(".", "_")
	v.SetEnvKeyReplacer(replacer)

	// Set defaults with priority to environment variables
	setPriorityDefaults(v, replacer)

	// Unmarshal configuration
	if err := v.Unmarshal(&cfg); err != nil {
		log.Panicf("Gagal memuat konfigurasi: %v", err)
	}

	return &cfg
}

func setPriorityDefaults(v *viper.Viper, replacer *strings.Replacer) {
	// Force binding of specific environment variables
	bindings := setEnvBindings()
	for runtimeKey, envKey := range bindings {
		v.BindEnv(runtimeKey, envKey)
	}

	defaults := setDefaults()

	// log.Printf("Scan Values:")
	for _, runtimeKey := range v.AllKeys() {
		runtimeValue := v.Get(runtimeKey)
		envFilekey := replacer.Replace(runtimeKey)
		if runtimeKey != envFilekey {
			if runtimeValue == nil {
				envFileValue := v.Get(envFilekey)
				if envFileValue != nil {
					// log.Printf(" %s = %v -> [%s]", runtimeKey, envFileValue, envFilekey)
					v.SetDefault(runtimeKey, envFileValue)
				} else if defValue, ok := defaults[runtimeKey]; ok {
					// log.Printf(" %s = %v -> [DEFAULTS]", runtimeKey, defValue)
					v.SetDefault(runtimeKey, defValue)
				}
			} else {
				// log.Printf(" %s = %v -> [RUNTIME]", runtimeKey, runtimeValue)
			}
		}
	}
}
