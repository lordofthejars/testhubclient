package hub

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/resty.v1"
)

func SendReport(reportFile *os.File, server string, project string, build string) error {

	// we can remove file after sent it to test hub
	defer os.Remove(reportFile.Name())

	dat, err := ioutil.ReadFile(reportFile.Name())

	if err != nil {
		return err
	}

	Debug("Sending report to %s/%s/%s", server, project, build)

	resp, err := resty.R().
		SetHeader("Content-Type", "application/gzip").
		SetPathParams(map[string]string{
			"project": project,
			"build":   build,
		}).
		SetBody(dat).
		Post(server + "/api/{project}/{build}")

	fmt.Printf("\nError: %v", err)
	fmt.Printf("\nResponse Status Code: %v", resp.StatusCode())
	fmt.Printf("\nResponse Status: %v", resp.Status())
	fmt.Printf("\nResponse Body: %v", resp)
	fmt.Printf("\nResponse Time: %v", resp.Time())
	fmt.Printf("\nResponse Recevied At: %v", resp.ReceivedAt())

	return err
}
