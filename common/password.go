package common

import (
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type PasswordConfig struct {
	IncludeChars     bool
	IncludeDigits    bool
	MinLength        int
	MaxLength        int
	IncludeUppercase bool
	IncludeLowercase bool
}

type OtpConfig struct {
	ExpireTime time.Duration
	Digits     int
	Limiter    time.Duration
}

var (
	lowerCharSet   = "abcdedfghijklmnopqrst"
	upperCharSet   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	specialCharSet = "!@#$%&*"
	numberSet      = "0123456789"
	allCharSet     = lowerCharSet + upperCharSet + specialCharSet + numberSet
)

func (cfg PasswordConfig) CheckPassword(password string) bool {
	if len(password) < cfg.MinLength {
		return false
	}

	if cfg.IncludeChars && !HasLetter(password) {
		return false
	}

	if cfg.IncludeDigits && !HasDigits(password) {
		return false
	}

	if cfg.IncludeLowercase && !HasLower(password) {
		return false
	}

	if cfg.IncludeUppercase && !HasUpper(password) {
		return false
	}

	return true
}

func (cfg PasswordConfig) GeneratePassword() string {
	var password strings.Builder

	passwordLength := cfg.MinLength + 2
	minSpecialChar := 2
	minNum := 3
	if !cfg.IncludeDigits {
		minNum = 0
	}

	minUpperCase := 3
	if !cfg.IncludeUppercase {
		minUpperCase = 0
	}

	minLowerCase := 3
	if !cfg.IncludeLowercase {
		minLowerCase = 0
	}

	//Set special character
	for i := 0; i < minSpecialChar; i++ {
		random := rand.Intn(len(specialCharSet))
		password.WriteString(string(specialCharSet[random]))
	}

	//Set numeric
	for i := 0; i < minNum; i++ {
		random := rand.Intn(len(numberSet))
		password.WriteString(string(numberSet[random]))
	}

	//Set uppercase
	for i := 0; i < minUpperCase; i++ {
		random := rand.Intn(len(upperCharSet))
		password.WriteString(string(upperCharSet[random]))
	}

	//Set lowercase
	for i := 0; i < minLowerCase; i++ {
		random := rand.Intn(len(lowerCharSet))
		password.WriteString(string(lowerCharSet[random]))
	}

	remainingLength := passwordLength - minSpecialChar - minNum - minUpperCase
	for i := 0; i < remainingLength; i++ {
		random := rand.Intn(len(allCharSet))
		password.WriteString(string(allCharSet[random]))
	}
	inRune := []rune(password.String())
	rand.Shuffle(len(inRune), func(i, j int) {
		inRune[i], inRune[j] = inRune[j], inRune[i]
	})
	return string(inRune)
}

func (cfg OtpConfig) GenerateOtp() string {
	rand.Seed(time.Now().UnixNano())
	min := int(math.Pow(10, float64(cfg.Digits-1)))   // 10^d-1 100000
	max := int(math.Pow(10, float64(cfg.Digits)) - 1) // 999999 = 1000000 - 1 (10^d) -1

	var num = rand.Intn(max-min) + min
	return strconv.Itoa(num)
}