package s3

import (
	"context"
	"testing"

	"github.com/trackit/jsonlog"

	taws "github.com/trackit/trackit2/aws"
)

func init() {
	jsonlog.DefaultLogger = jsonlog.DefaultLogger.WithLogLevel(jsonlog.LogLevelDebug)
}

/*
func TestEpitechio(t *testing.T) {
	var count int
	var start time.Time
	var end time.Time
	time.Sleep(10 * time.Second)
	err := ReadBills(
		context.Background(),
		// taws.AwsAccount{
		// 	RoleArn:  "arn:aws:iam::895365654851:role/trackit",
		// 	External: "RLuxJFYhaZYjWHNYY_pfeAgF@lzymhUKNxiwq_IQ",
		// },
		taws.AwsAccount{
			RoleArn:  "arn:aws:iam::394125495069:role/delegation-trackit",
			External: "5e22M7zqmMNo7t4G7SUSD9IK",
		},
		// BillRepository{
		// 	Bucket: "epitechio-reports",
		// 	Prefix: "constandusage/",
		// },
		BillRepository{
			Bucket: "trackit-billing-report",
			Prefix: "usagecost/AllHourlyToS3/",
		},
		func(li LineItem, ok bool) bool {
			if ok {
				if count == 0 {
					start = time.Now()
				}
				count++
			} else {
				end = time.Now()
			}
			return true
		},
		acceptAllManifests,
	)
	fmt.Printf("Parsed %d records in %s.", count, end.Sub(start).String())
	if err != nil {
		println(err.Error())
	}
}
*/

/*
func TestEpitechio(t *testing.T) {
	var count int
	var start time.Time
	var end time.Time
	client, err := elastic.NewClient(
		elastic.SetBasicAuth("elastic", "changeme"),
	)
	if err != nil {
		println(err.Error())
		return
	}
	err = ReadBills(
		context.Background(),
		// taws.AwsAccount{
		// 	RoleArn:  "arn:aws:iam::895365654851:role/trackit",
		// 	External: "RLuxJFYhaZYjWHNYY_pfeAgF@lzymhUKNxiwq_IQ",
		// },
		taws.AwsAccount{
			RoleArn:  "arn:aws:iam::394125495069:role/delegation-trackit",
			External: "5e22M7zqmMNo7t4G7SUSD9IK",
		},
		// BillRepository{
		// 	Bucket: "epitechio-reports",
		// 	Prefix: "constandusage/",
		// },
		BillRepository{
			Bucket: "trackit-billing-report",
			Prefix: "usagecost/AllHourlyToS3/",
		},
		addToElasticsearch(client),
		acceptAllManifests,
	)
	fmt.Printf("Parsed %d records in %s.", count, end.Sub(start).String())
	if err != nil {
		println(err.Error())
	}
}
*/

func TestUpdate(t *testing.T) {
	latestManifest, err := UpdateReport(
		context.Background(),
		taws.AwsAccount{
			RoleArn:  "arn:aws:iam::394125495069:role/delegation-trackit",
			External: "5e22M7zqmMNo7t4G7SUSD9IK",
		},
		BillRepository{
			Bucket: "trackit-billing-report",
			Prefix: "usagecost/AllHourlyToS3/",
		},
	)
	println(latestManifest.String())
	if err != nil {
		println(err.Error())
	}
}
