package signed_elasticsearch_client

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/sha1sum/aws_signing_client"
	"gopkg.in/olivere/elastic.v5"
	"strings"
)

func checkError(splittedEndpoint []string, creds *credentials.Credentials) error {
	if creds == nil {
		return errors.New("Nil credential")
	} else if _, err := creds.Get(); err != nil {
		return err
	} else if len(splittedEndpoint) <= 4 {
		return errors.New("Wrong endpoint parameter")
	}
	return nil
}

func NewSignedElasticClient(endpoint string, creds *credentials.Credentials) (*elastic.Client, error) {
	signer := v4.NewSigner(creds)
	splittedEndpoint := strings.Split(endpoint, ".")
	if err := checkError(splittedEndpoint, creds); err != nil {
		return nil, err
	}
	region := splittedEndpoint[len(splittedEndpoint)-4]
	awsClient, err := aws_signing_client.New(signer, nil, "es", region)
	if err != nil {
		return nil, err
	}
	prefix := ""
	if !strings.HasPrefix(endpoint, "http") {
		prefix = "https://"
	}
	return elastic.NewClient(
		elastic.SetURL(prefix+endpoint),
		elastic.SetScheme("https"),
		elastic.SetHttpClient(awsClient),
		elastic.SetSniff(false),
	)
}
