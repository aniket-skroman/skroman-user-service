package helper

import (
	"errors"
	"fmt"
	"strings"

	"github.com/lib/pq"
)

func Handle_DBError(err error) (err_ error) {
	switch e := err.(type) {
	case *pq.Error:
		fmt.Println("PQ Code : ", e.Code)
		switch e.Code {
		case "23502":
			// not-null constraint violation
			fmt.Println("Some required data was left out:\n\n", e.Message)
			err_ = errors.New(e.Detail)
			return
		case "23505":
			// unique constraint violation
			if strings.Contains(e.Message, "email") {
				err_ = Err_Account_Already_Exists
				return
			}
			if strings.Contains(e.Message, "emp_code") {
				err_ = Err_EMP_Code_Exists
				return
			}
			err_ = errors.New(e.Detail)
			return

		case "23514":
			fmt.Println("Handle_DBError called from constraint check")

			// check constraint violation
			if strings.Contains(e.Message, "contact") {
				err_ = errors.New("contact should not be empty")
				return
			} else if strings.Contains(e.Message, "email") {
				err_ = errors.New("email should not be empty")
				return
			}
			// err_ = validate_err_msg(&e.Message)
			// return
		case "23503":
			err_ = errors.New("invalid id has been provided,please try with valid id's")
			return
		case "2201X":
			// when offset, limit not working, for pagination data
			err_ = errors.New("invalid pagination found")
			return
		default:
			msg := e.Message
			if d := e.Detail; d != "" {
				msg += "\n\n" + d
			}
			if h := e.Hint; h != "" {
				msg += "\n\n" + h
			}
			fmt.Println("Message from default : ", e.Code)
			err_ = errors.New(msg)
			return
		}
	default:
		err_ = nil
		return
	}

	return
}

func validate_err_msg(err_msg *string) error {
	var err error

	// if strings.Contains(*err_msg, "contact") {
	// 	err = errors.New("contact should not be empty")
	// } else if strings.Contains(*err_msg, "email") {
	// 	err = errors.New("email should not be empty")
	// }else if {}

	switch {
	case strings.Contains(*err_msg, "users_full_name_key"):
		err = errors.New("full name is already exists")
	default:
		err = errors.New("can't process the request")
	}

	return err
}
