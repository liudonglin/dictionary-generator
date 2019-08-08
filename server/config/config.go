package config

import (
	"strings"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

func init() {
	// envFile := "./.env"

	// driver := flag.String("driver", "sqlite3", "db driver, mysql or sqlite3")
	// dbsource := flag.String("data source", "codebuild.db", "data source")
	// admin := flag.String("admin", "admin:123456", "init admin user")
	// flag.Parse()

	// if _, err := os.Stat(envFile); os.IsNotExist(err) {
	// 	f, err1 := os.Create(envFile)
	// 	if err1 != nil {
	// 		log.Fatalf("env文件创建失败: %s", err1.Error())
	// 	}

	// 	envString := fmt.Sprintf("CODEBUILDER_DATABASE_DRIVER=%s\n", *driver)
	// 	envString += fmt.Sprintf("CODEBUILDER_DATABASE_DATASOURCE=%s\n", *dbsource)
	// 	envString += fmt.Sprintf("CODEBUILDER_USER_CREATE=%s\n", *admin)
	// 	_, err1 = io.WriteString(f, envString)
	// 	if err1 != nil {
	// 		log.Fatalf("env文件写入失败: %s", err1.Error())
	// 	}
	// }

	// 从.env文件中加载配置信息
	godotenv.Load()
}

// Config provides the system configuration.
type Config struct {

	// Database provides the database configuration.
	Database struct {
		Driver     string `envconfig:"CODEBUILDER_DATABASE_DRIVER"     default:"sqlite3"`
		Datasource string `envconfig:"CODEBUILDER_DATABASE_DATASOURCE" default:"codebuild.db"`
	}

	Admin AdminCreate `envconfig:"CODEBUILDER_USER_CREATE" default:"admin:123456"`
}

// Environ returns the settings from the environment.
func Environ() (Config, error) {
	cfg := Config{}
	err := envconfig.Process("", &cfg)
	return cfg, err
}

// AdminCreate ...
type AdminCreate struct {
	Username string
	Password string
}

// Decode ...
func (u *AdminCreate) Decode(value string) error {
	parts := strings.Split(value, ":")
	if len(parts) != 2 {
		return nil
	}
	u.Username = parts[0]
	u.Password = parts[1]
	return nil
}
