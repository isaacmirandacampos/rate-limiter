package configs

import "github.com/spf13/viper"

type conf struct {
	Timeout                     int64  `mapstructure:"TIMEOUT"`
	RedisAddress                string `mapstructure:"REDIS_ADDRESS"`
	RequestsPerSecondByIp       int32  `mapstructure:"REQUESTS_PER_SECOND_BY_IP"`
	RequestsPerSecondByApiToken int32  `mapstructure:"REQUESTS_PER_SECOND_BY_API_TOKEN"`
}

func LoadConfig(path string) (*conf, error) {
	var cfg *conf
	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}
	return cfg, err
}
