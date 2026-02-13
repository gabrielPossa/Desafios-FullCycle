package configs

import "github.com/spf13/viper"

type conf struct {
	DBDriver          string         `mapstructure:"DB_DRIVER"`
	DBHost            string         `mapstructure:"DB_HOST"`
	DBPort            string         `mapstructure:"DB_PORT"`
	DBUser            string         `mapstructure:"DB_USER"`
	DBPassword        string         `mapstructure:"DB_PASSWORD"`
	DBName            string         `mapstructure:"DB_NAME"`
	WebServerPort     string         `mapstructure:"WEB_SERVER_PORT"`
	GRPCServerPort    string         `mapstructure:"GRPC_SERVER_PORT"`
	GraphQLServerPort string         `mapstructure:"GRAPHQL_SERVER_PORT"`
	RabbitMQ          RabbitMQConfig `mapstructure:",squash"`
}

type RabbitMQConfig struct {
	Host string `mapstructure:"RABBITMQ_HOST"`
	Port string `mapstructure:"RABBITMQ_PORT"`
	User string `mapstructure:"RABBITMQ_USER"`
	Pass string `mapstructure:"RABBITMQ_PASS"`
}

func LoadConfig(path string) (*conf, error) {
	var cfg *conf
	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	_ = viper.BindEnv("DB_DRIVER")
	_ = viper.BindEnv("DB_HOST")
	_ = viper.BindEnv("DB_PORT")
	_ = viper.BindEnv("DB_USER")
	_ = viper.BindEnv("DB_PASSWORD")
	_ = viper.BindEnv("DB_NAME")
	_ = viper.BindEnv("WEB_SERVER_PORT")
	_ = viper.BindEnv("GRPC_SERVER_PORT")
	_ = viper.BindEnv("GRAPHQL_SERVER_PORT")
	_ = viper.BindEnv("RABBITMQ_HOST")
	_ = viper.BindEnv("RABBITMQ_PORT")
	_ = viper.BindEnv("RABBITMQ_USER")
	_ = viper.BindEnv("RABBITMQ_PASS")

	_ = viper.ReadInConfig()
	err := viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}
	return cfg, err
}
