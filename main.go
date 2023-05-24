package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

func main() {
	prefixFlag := flag.String("prefix", "", "Add a common prefix to all created env variables")
	regionFlag := flag.String("region", "eu-central-1", "Which AWS region we should use")
	idFlag := flag.String("secret-id", "", "Name or ARN of the secret to retrieve")

	flag.Parse()

	data := make(map[string]string)

	if *idFlag == "" {
		fi, err := os.Stdin.Stat()
		if err != nil {
			panic(err)
		}

		if fi.Mode()&os.ModeNamedPipe == 0 {
			// No piped input into the tool
			os.Exit(0)
		}

		input, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Println(err)
		}
		json.Unmarshal(input, &data)

	} else {
		session := session.Must(session.NewSession())
		svc := secretsmanager.New(session, aws.NewConfig().WithRegion(*regionFlag))

		x, err := svc.GetSecretValue(&secretsmanager.GetSecretValueInput{
			SecretId: idFlag,
		})

		if err != nil {
			fmt.Println(err)
		}

		json.Unmarshal([]byte(*x.SecretString), &data)
	}

	// print all variables as export statements to the terminal
	for key, value := range data {
		fmt.Printf("export %s%s=%s\n", *prefixFlag, key, value)
	}
}
