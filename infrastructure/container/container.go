package container

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type (
	Container struct {
		Apps  *Apps
		Redis *Redis
	}

	Apps struct {
		Name    string `json:"name"`
		Version string `json:"version"`
		Host    string `json:"host"`
		Port    int    `json:"port"`
	}

	Redis struct {
		RedisAddr      string
		RedisPassword  string
		RedisDB        int
		RedisDefaultdb int
		MinIdleConns   int
		PoolSize       int
		PoolTimeout    int
	}
)

func (c *Container) Validate() *Container {
	if c.Apps == nil {
		panic("Apps config is nill")
	}
	// if c.Redis == nil {
	// 	panic("Redis config is nill")
	// }

	return c
}

func New(envpath string) *Container {
	v := viper.New()
	v.SetConfigFile(envpath)
	pathDir, err := os.Executable()
	if err != nil {
		panic(err)
	}
	dir := filepath.Dir(pathDir)
	v.AddConfigPath(dir)

	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	appHost := v.GetString("app.host")
	appPort := v.GetInt("app.port")
	appName := v.GetString("app.name")
	appVersion := v.GetString("app.version")

	redisAddr := v.GetString("redis.Addr")
	redisPassword := v.GetString("redis.Password")
	redisDB := v.GetInt("redis.Db")
	redisDefaultDB := v.GetInt("redis.Defaultdb")
	redisMinIdleConns := v.GetInt("redis.MinIdleConns")
	redisPoolSize := v.GetInt("redis.PoolSize")
	redisPoolTimeout := v.GetInt("redis.PoolTimeout")

	appConf := &Apps{
		Name:    appName,
		Version: appVersion,
		Host:    appHost,
		Port:    appPort,
	}

	redisConf := &Redis{
		RedisAddr:      redisAddr,
		RedisPassword:  redisPassword,
		RedisDB:        redisDB,
		RedisDefaultdb: redisDefaultDB,
		MinIdleConns:   redisMinIdleConns,
		PoolSize:       redisPoolSize,
		PoolTimeout:    redisPoolTimeout,
	}

	cont := &Container{
		Apps:  appConf,
		Redis: redisConf,
	}

	cont.Validate()

	return cont

}
