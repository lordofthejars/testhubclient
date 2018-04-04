package hub

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/resty.v1"
)

type HttpError struct {
	Resp *resty.Response
}

func (e *HttpError) Error() string {
	return fmt.Sprintf("Http error with status code %d and message %v", e.Resp.StatusCode(), e.Resp)
}

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

	if err != nil {
		return err
	}

	return checkResponse(resp)
}

func DeleteBuild(server string, project string, build string) error {

	Debug("Deleting Build %s/%s/%s", server, project, build)

	resp, err := resty.R().
		SetPathParams(map[string]string{
			"project": project,
			"build":   build,
		}).
		Delete(server + "/api/{project}/{build}")

	if err != nil {
		return err
	}

	return checkResponse(resp)
}

func checkResponse(resp *resty.Response) error {
	statusCode := resp.StatusCode()

	if statusCode < 200 || statusCode > 299 {
		return &HttpError{
			resp,
		}
	}

	return nil
}
