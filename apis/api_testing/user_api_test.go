package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"testing"

	db "github.com/aniket-skroman/skroman-user-service/sqlc_lib"
	"github.com/aniket-skroman/skroman-user-service/utils"
	"github.com/stretchr/testify/require"
)

func TestCreaUserAPI(t *testing.T) {
	args := []struct {
		TestName      string
		RequestBody   db.CreateNewUserParams
		ExpectedCode  int
		ExpectedError bool
	}{
		{TestName: "First",
			RequestBody: db.CreateNewUserParams{
				FullName: utils.RandomString(7),
				Email:    fmt.Sprintf("%s@gmail.com", utils.RandomString(5)),
				Password: "user123",
				Contact:  "1234567890",
				UserType: "ADMIN",
			},
			ExpectedCode:  http.StatusCreated,
			ExpectedError: false,
		},
		{TestName: "Second",
			RequestBody: db.CreateNewUserParams{
				FullName: utils.RandomString(12),
				Email:    fmt.Sprintf("%s@gmail.com[]", utils.RandomString(8)),
				Password: "user123",
				Contact:  "1234567890",
				UserType: "ADMIN",
			},
			ExpectedCode:  http.StatusBadRequest,
			ExpectedError: true,
		},
		{TestName: "Third",
			RequestBody: db.CreateNewUserParams{
				FullName: utils.RandomString(12),
				Email:    fmt.Sprintf("%s@gmail.com", utils.RandomString(5)),
				Password: "",
				Contact:  "1234567890",
				UserType: "ADMIN",
			},
			ExpectedCode:  http.StatusBadRequest,
			ExpectedError: true,
		},
		{TestName: "Fourth",
			RequestBody: db.CreateNewUserParams{
				FullName: utils.RandomString(12),
				Email:    fmt.Sprintf("%s@gmail.com", utils.RandomString(4)),
				Password: "",
				Contact:  "123456789022121",
				UserType: "ADMIN",
			},
			ExpectedCode:  http.StatusBadRequest,
			ExpectedError: true,
		},
		{
			TestName: "Five",
			RequestBody: db.CreateNewUserParams{
				FullName: utils.RandomString(12),
				Email:    "rohit.dhavale-1@skromanglobal.com",
				Password: "user123",
				Contact:  "7720830172",
				UserType: "EMP",
			},
			ExpectedCode:  http.StatusConflict,
			ExpectedError: true,
		},
	}

	log_file, err := os.Create("app.log")

	if err != nil {
		log.Fatal(err)
	}

	debug_logger := log.New(log_file, "DEBUG : ", log.Flags())

	// url := "http://13.233.196.149:8080/api/create-user"
	url := "http://localhost:8080/api/create-user"

	for _, arg := range args {
		t.Run(arg.TestName, func(t *testing.T) {
			req_body, err := json.Marshal(arg.RequestBody)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(req_body))
			require.NoError(t, err)

			response, err := http.DefaultClient.Do(request)
			response_body, err := io.ReadAll(response.Body)

			debug_logger.Println(arg.TestName)
			debug_logger.Println("REQUEST : ", arg.RequestBody)
			debug_logger.Println("RESPONSE : ", string(response_body))
			debug_logger.Println("RESPONSE ERROR : ", err)
			debug_logger.Println("RESPONSE STATUS CODE : ", response.StatusCode)
			debug_logger.Println("EXPECTED STATUS CODE : ", arg.ExpectedCode)
			debug_logger.Println()

			// if arg.ExpectedError {
			// 	require.Error(t, err)
			// } else {
			// 	require.NoError(t, err)
			// }

			fmt.Println("Response err : ", err, string(response_body))
			require.Equal(t, arg.ExpectedCode, response.StatusCode)

		})
	}
}
