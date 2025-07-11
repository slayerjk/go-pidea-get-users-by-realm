package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"strings"
	"time"

	pidea "github.com/slayerjk/go-pideaapi"
	"golang.org/x/term"

	vafswork "github.com/slayerjk/go-vafswork"
	vawebwork "github.com/slayerjk/go-vawebwork"
)

const (
	appName = "pi-get-users-by-realm"
)

type piData struct {
	PideaURL         string `json:"pideaUrl"`
	PideaApiUser     string `json:"pideaApiUser"`
	PideaRealm       string
	PideaApiPassword string
	PideaApiToken    string
}

func main() {
	// defining default values
	var (
		workDir         string    = vafswork.GetExePath()
		logsPathDefault string    = workDir + "/logs" + "_" + appName
		startTime       time.Time = time.Now()
		dataFile        string    = workDir + "/data.json"
		resultsDir      string    = workDir + "/Results"
		piData          piData
	)

	// flags
	logsDir := flag.String("log-dir", logsPathDefault, "set custom log dir")
	logsToKeep := flag.Int("keep-logs", 7, "set number of logs to keep after rotation")
	realm := flag.String("realm", "NONE", "Pidea realm to get users from")

	flag.Usage = func() {
		fmt.Println("Get PrivacyIdea users by Realm")
		fmt.Println("Version = 0.0.1")
		fmt.Println("Usage: <app> -r <YOUR REALM>")
		fmt.Println("Flags:")
		flag.PrintDefaults()
	}

	flag.Parse()

	// logging
	// create log dir
	if err := os.MkdirAll(*logsDir, os.ModePerm); err != nil {
		fmt.Fprintf(os.Stdout, "failed to create log dir %s:\n\t%v", *logsDir, err)
		os.Exit(1)
	}
	// set current date
	dateNow := time.Now().Format("02.01.2006")
	// create log file
	logFilePath := fmt.Sprintf("%s/%s_%s.log", *logsDir, appName, dateNow)
	// open log file in append mode
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Fprintf(os.Stdout, "failed to open created log file %s:\n\t%v", logFilePath, err)
		os.Exit(1)
	}
	defer logFile.Close()
	// set logger
	logger := slog.New(slog.NewTextHandler(logFile, nil))

	// starting programm notification
	logger.Info("Program Started", "app name", appName)

	// rotate logs
	logger.Info("Log rotation first", "logsDir", *logsDir, "logs to keep", *logsToKeep)
	if err := vafswork.RotateFilesByMtime(*logsDir, *logsToKeep); err != nil {
		fmt.Fprintf(os.Stdout, "failed to rotate logs:\n\t%v", err)
	}

	// main code here
	if strings.ToUpper(*realm) == "NONE" {
		fmt.Println("realm is not set, use '-r <YOUR REALM>' flag")
		os.Exit(1)
	}
	piData.PideaRealm = *realm

	// make results dir
	err = os.MkdirAll(resultsDir, os.ModePerm)
	if err != nil {
		logger.Error("failed to create Results dir")
	}

	// check data.json exists
	if _, err := os.Stat(dataFile); errors.Is(err, os.ErrNotExist) {
		logger.Error("data file doesn't exist, exiting")
		os.Exit(1)
	}

	// reading json data
	dataBytes, err := os.ReadFile(dataFile)
	if err != nil {
		logger.Error("failed to read Data file, exiting", "err", err)
		os.Exit(1)
	}

	// check if data.json is valid json
	if !json.Valid(dataBytes) {
		logger.Error("data file is not Valid json file, exiting")
		os.Exit(1)
	}

	// writing data file into a struct
	err = json.Unmarshal(dataBytes, &piData)
	if err != nil {
		logger.Error("failed to Unmarshal data file, exiting")
		os.Exit(1)
	}

	// getting API password from user input
	for {
		fmt.Print("Enter Pidea API user's Password and press Enter: ")
		passwordBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			fmt.Println("error, check logs")
			logger.Error("failed to get user's password, exiting", "err", err)
			os.Exit(1)
		}

		if len(string(passwordBytes)) == 0 {
			fmt.Println("")
			fmt.Println("password may not be empty")
			logger.Error("empty password entered")
			continue
		}

		piData.PideaApiPassword = string(passwordBytes)
		break
	}

	// create HTTP Client
	httpClient := vawebwork.NewInsecureClient()

	// getting API token
	piData.PideaApiToken, err = pidea.GetApiToken(&httpClient, piData.PideaURL, piData.PideaApiUser, piData.PideaApiPassword)
	if err != nil {
		fmt.Println(err)
		logger.Error("failed to get Pidea API Token", "err", err)
		os.Exit(1)
	}
	// fmt.Println("Token:", piData.PideaApiToken)

	// getting Pidea users by realm
	usersResult, err := pidea.GetPideaUsersByRealm(&httpClient, piData.PideaURL, piData.PideaRealm, piData.PideaApiToken)
	if err != nil {
		fmt.Println(err)
		logger.Error("failed to get Pidea Users, exiting", "err", err)
		os.Exit(1)
	}

	// for _, user := range usersResult {
	// 	fmt.Println(user)
	// }

	// making csv result file for users
	resultFilePath := fmt.Sprintf("%s/result_%s_%s.csv", resultsDir, *realm, dateNow)

	resultFile, err := os.Create(resultFilePath)
	if err != nil {
		logger.Error("failed to create result CSV file, exiting", "err", err)
	}
	defer resultFile.Close()

	csvWriter := csv.NewWriter(resultFile)
	defer csvWriter.Flush()

	// getting headers for csv from User struct
	for i, user := range usersResult {
		rowToWrite := make([]string, 0)
		// form CSV headers
		if i == 0 {
			csvHeaders := GetStructKeys(user)
			if len(csvHeaders) == 0 {
				fmt.Println("internal error, check logs")
				logger.Error("empty headers list for result csv, check logic")
			}
			if err = csvWriter.Write(csvHeaders); err != nil {
				fmt.Println("internal error, check logs")
				logger.Error("failed to write headers to csv")
			}
		}

		// get string value of struct fields
		values := reflect.ValueOf(user)
		for j := 0; j < values.NumField(); j++ {
			stringToWrite := fmt.Sprintf("%v", values.Field(j))
			rowToWrite = append(rowToWrite, stringToWrite)
		}

		// write user row in result
		csvWriter.Write(rowToWrite)
	}

	resultFile.Sync()
	fmt.Printf("\nDONE, result is in: %s", resultFilePath)

	// count & print estimated time
	logFile.Close()
	endTime := time.Now()
	logger.Info("Program Done", slog.Any("estimated time(sec)", endTime.Sub(startTime).Seconds()))
}
