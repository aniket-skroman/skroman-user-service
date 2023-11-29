package helper

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

type InputValidator interface {
	Validate() (interface{}, error)
}

type ContactValidate struct {
	Contact string
}

func (cv *ContactValidate) Validate() (interface{}, error) {
	regex := `^[0-9]{10}$`
	match, err := regexp.MatchString(regex, cv.Contact)
	return match, err
}

func valide_actual_input(input_data InputValidator) (interface{}, error) {
	return input_data.Validate()
}

func ValidateInput(input_data interface{}) (interface{}, error) {
	if data, ok := input_data.(string); ok {
		contact_val := ContactValidate{Contact: data}
		var validate InputValidator = &contact_val
		return valide_actual_input(validate)
	}

	return false, Err_Invalid_Input
}

type ApiError struct {
	Field string
	Msg   string
}

func Error_handler(err error) []ApiError {
	if err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]ApiError, len(ve))
			for i, fe := range ve {
				out[i] = ApiError{fe.Field(), msgForTag(fe.Tag())}
			}
			return out
		}
		return nil
	}
	return nil
}

func Single_Error_handler(err error) ApiError {
	if err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			for _, fe := range ve {
				return ApiError{fe.Field(), msgForTag(fe.Tag())}
			}

		}
		return ApiError{}
	}
	return ApiError{}
}

func Handle_required_param_error(err error) string {
	var ve validator.ValidationErrors
	var err_msg string
	if errors.As(err, &ve) {
		for _, fe := range ve {
			err_msg = fmt.Sprintf("%v - %v", fe.Field(), msgForTag(fe.Tag()))
			break
		}
	} else {
		fmt.Println("Error MSG : ", err)
		if strings.Contains(err.Error(), "cannot unmarshal string into") {
			err_msg = "required a integer but found string, please check params"
		} else if strings.Contains(err.Error(), "cannot unmarshal number into") {
			err_msg = "required a string but found integer, please check params"
		} else {
			err_msg = "something went's wrong, invalid param detecte"

		}
	}

	return err_msg
}

func msgForTag(tag string) string {
	switch tag {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email"
	case "oneof":
		return "Invalid params"
	}
	return ""
}

var key = "skroman-user-servi-12345"

func EncryptData(plaintext string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(plaintext))

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func DecryptData(ciphertext string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	decodedCiphertext, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	iv := decodedCiphertext[:aes.BlockSize]
	decodedCiphertext = decodedCiphertext[aes.BlockSize:]

	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(decodedCiphertext, decodedCiphertext)

	return string(decodedCiphertext), nil
}
