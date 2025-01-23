package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	Port      int             `yaml:"port" mapstructure:"port"`
	ZooKeeper ZooKeeperConfig `yaml:"zookeeper" mapstructure:"zookeeper"`
	// ... 其他配置项
}

// LoadEnv 加载 .env 文件，支持指定环境
func LoadEnv(env string) error {
	var envPath string

	// 优先查找指定环境的文件
	if env != "" {
		envPath = fmt.Sprintf(".env.%s", env)
		if _, err := os.Stat(envPath); os.IsNotExist(err) {
			fmt.Printf("env file %s not found, try .env\n", envPath)
		} else if err != nil {
			return fmt.Errorf("stat env file %s error: %w", envPath, err)
		} else {
			fmt.Printf("load env file %s\n", envPath)
			return godotenv.Load(envPath)
		}
	}

	// 如果指定环境的文件不存在，则查找 .env 文件
	envPath = ".env"
	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		//如果既没有.env.uat也没有.env文件，则不报错，使用环境变量
		fmt.Println("env file .env not found")
		return nil
	} else if err != nil {
		return fmt.Errorf("stat env file .env error: %w", err)
	}

	fmt.Printf("load env file .env\n")
	return godotenv.Load(envPath)
}

func LoadEnvWithDir(dir, env string) error {
	var envPath string
	if dir == "" {
		dir, _ = os.Getwd()
	}
	if !strings.HasSuffix(dir, "/") {
		dir = dir + "/"
	}
	// 优先查找指定环境的文件
	if env != "" {
		envPath = filepath.Join(dir, fmt.Sprintf(".env.%s", env))
		if _, err := os.Stat(envPath); os.IsNotExist(err) {
			fmt.Printf("env file %s not found, try .env\n", envPath)
		} else if err != nil {
			return fmt.Errorf("stat env file %s error: %w", envPath, err)
		} else {
			fmt.Printf("load env file %s\n", envPath)
			return godotenv.Load(envPath)
		}
	}

	// 如果指定环境的文件不存在，则查找 .env 文件
	envPath = filepath.Join(dir, ".env")
	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		fmt.Println("env file .env not found")
		return nil
	} else if err != nil {
		return fmt.Errorf("stat env file .env error: %w", envPath, err)
	}

	fmt.Printf("load env file .env\n")
	return godotenv.Load(envPath)
}
