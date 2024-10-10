package helper

import (
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"golang.org/x/crypto/bcrypt"
)

func SetDefaultTimezone() {
	loc, err := time.LoadLocation("UTC")
	if err != nil {
		loc = time.Now().Location()
	}

	time.Local = loc
}

func EnvConfigVariable(filePath string) (cfg *viper.Viper, err error) {
	cfg = viper.New()
	cfg.SetConfigFile(filePath)

	if err = cfg.ReadInConfig(); err != nil {
		err = errors.Wrap(err, "Error while reading config file")

		return
	}

	return
}

func GenerateHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func ValidatePassword(password string, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return false
	}
	return true
}
