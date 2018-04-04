package hub

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/resty.v1"
)

func SendReport(reportFile *os.File, server string, project string, build string, reportType string) error {

	// we can remove file after sent it to test hub
	defer os.Remove(reportFile.Name())

	dat, err := ioutil.ReadFile(reportFile.Name())

	if err != nil {
		return err
	}

	Debug("Sending report to %s/%s/%s of type %s", server, project, build, reportType)

	resp, err := resty.R().
		SetHeader("Content-Type", "application/gzip").
		SetHeader("x-testhub-type", reportType).
		SetPathParams(map[string]string{
			"project": project,
			"build":   build,
		}).
		SetBody(dat).
		Post(server + "/api/{project}/{build}")

	// Need to check status code logic
	fmt.Printf("\nResponse Status Code: %v", resp.StatusCode())
	fmt.Printf("\nResponse Status: %v", resp.Status())
	fmt.Printf("\nResponse Body: %v", resp)
	fmt.Printf("\nResponse Time: %v", resp.Time())
	fmt.Printf("\nResponse Recevied At: %v", resp.ReceivedAt())

	return err
}
