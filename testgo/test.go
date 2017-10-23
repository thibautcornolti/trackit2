package main

import (
	"./costs"
	"./es"
	"./parsator"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	//"gopkg.in/olivere/elastic.v5"
)

func main() {
	client, err := signed_elasticsearch_client.NewSignedElasticClient("search-job-msol-prod-tracstg-tnqmu4c6vijsgul7rx5moba2ce.us-west-2.es.amazonaws.com", credentials.NewSharedCredentials("", ""))
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	ress, err := costs.GetElasticSearchParams([]string{"394125495069"}, time.Date(2017, 10, 1, 0, 0, 0, 0, time.Local), time.Date(2017, 11, 1, 0, 0, 0, 0, time.Local), []string{"tag:Environment", "account", "year", "region", "month", "product", "day"}, client, "awsdetailedlineitem").Do(context.Background())
	//ress, err := client.Search().Index("awsdetailedlineitem").Size(1).Aggregation("&buckets.accounts", elastic.NewTermsAggregation().Field("linked_account_id").SubAggregation("&buckets.products", elastic.NewTermsAggregation().Field("product_name").Size(2).SubAggregation("&buckets.resources", elastic.NewTermsAggregation().Field("resource_id").Size(3).SubAggregation("*value.cost", elastic.NewSumAggregation().Field("cost")))).Size(5)).Do(context.Background())
	if err != nil {
		fmt.Printf("ERROR %v\n", err)
	} else {
		pars := parsator.GetParsedElasticSearchResult(ress)
		jsonRes, _ := json.MarshalIndent(pars, "", "  ")
		fmt.Printf("JSONRES\n%v\n", string(jsonRes))
	}
}
