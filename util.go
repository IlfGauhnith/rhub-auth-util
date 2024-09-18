package util

import (
	"log"
	"reflect"
	"runtime"

	"github.com/joho/godotenv"
)

func getFunctionName(i interface{}) string {
	// Obtain the program counter (PC) from the passed function
	pc := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Entry()

	// Get the function name from the PC
	functionName := runtime.FuncForPC(pc).Name()

	return functionName
}

func LogFunctionExecutionStart(i interface{}) {
	log.Printf("Executing %v", getFunctionName(i))
}

func LoadEnvironmentVariables() {
	LogFunctionExecutionStart(LoadEnvironmentVariables)
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	log.Print("Local environment variables sucessfully loaded")
}
