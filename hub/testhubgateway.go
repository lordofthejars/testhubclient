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

func SendReport(reportFile *os.File, options Options, reportType string) error {

	// we can remove file after sent it to test hub
	defer os.Remove(reportFile.Name())

	dat, err := ioutil.ReadFile(reportFile.Name())

	if err != nil {
		return err
	}

	server := options.URL
	project := options.Project
	build := options.Build

	Debug("Sending report to %s/%s/%s of type %s", server, project, build, reportType)

	resp, err := resty.R().
		SetHeader("Content-Type", "application/gzip").
		SetHeader("x-testhub-type", reportType).
		SetPathParams(map[string]string{
			"project": project,
			"build":   build,
		}).
		SetQueryParams(buildSendReportQueryParams(options)).
		SetBody(dat).
		Post(server + "/api/project/{project}/{build}")

	if err != nil {
		return err
	}

	return checkResponse(resp)
}

func buildSendReportQueryParams(options Options) map[string]string {

	var queryParams map[string]string
	queryParams = make(map[string]string)

	if options.IsBranchSet() {
		queryParams["branch"] = options.Branch
	}

	if options.IsBuildUrlSet() {
		queryParams["buildUrl"] = options.BuildURL
	}

	if options.IsCommitSet() {
		queryParams["commit"] = options.Commit
	}

	if options.IsRepoTypeSet() {
		queryParams["repoType"] = options.RepoType
	}

	if options.IsRepoUrlSet() {
		queryParams["repoUrl"] = options.RepoURL
	}

	return queryParams
}

func DeleteBuild(options Options) error {

	server := options.URL
	project := options.Project
	build := options.Build

	Debug("Deleting Build %s/%s/%s", server, project, build)

	resp, err := resty.R().
		SetPathParams(map[string]string{
			"project": project,
			"build":   build,
		}).
		Delete(server + "/api/project/{project}/{build}")

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
