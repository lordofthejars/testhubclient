package hub

import (
	"crypto/tls"
	"crypto/x509"
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

func Login(options Options) (string, error) {
	server := options.URL

	auth := make(map[string]string)

	auth["username"] = options.Username
	auth["password"] = options.Password

	err := configureHttps(options)

	if err != nil {
		return "", err
	}

	Debug("Authenticating against %s with username %s", server, options.Username)

	resp, err := resty.R().
		SetHeader("Content-Type", "application/json").
		SetBody(auth).
		Post(server + "/api/login")

	if err != nil {
		return "", err
	}
	err = checkResponse(resp)

	if err != nil {
		return "", err
	}

	return resp.String(), nil
}

func SendReportWithToken(reportFile *os.File, options Options, token string) error {
	partialRequest := resty.R().
		SetAuthToken(token)
	return sendReport(reportFile, options, partialRequest)
}

func SendReport(reportFile *os.File, options Options) error {
	partialRequest := resty.R()
	return sendReport(reportFile, options, partialRequest)
}

func sendReport(reportFile *os.File, options Options, partialRequest *resty.Request) error {
	// we can remove file after sent it to test hub
	defer os.Remove(reportFile.Name())

	dat, err := ioutil.ReadFile(reportFile.Name())

	if err != nil {
		return err
	}

	server := options.URL
	project := options.Project
	build := options.Build
	reportType := options.ReportTestType.ReportType

	err = configureHttps(options)

	if err != nil {
		return err
	}

	Debug("Sending report to %s/%s/%s/%s", server, project, build, reportType)

	partialRequest.
		SetHeader("Content-Type", "application/gzip").
		SetHeader("x-testhub-type", "html")

	resp, err := partialRequest.
		SetPathParams(map[string]string{
			"project": project,
			"build":   build,
			"report":  reportType,
		}).
		SetQueryParams(buildSendReportQueryParams(options)).
		SetBody(dat).
		Post(server + "/api/project/{project}/{build}/report/{report}")

	if err != nil {
		return err
	}

	return checkResponse(resp)
}

func SendTestReportWithToken(reportFile *os.File, options Options, reportType, token string) error {
	partialRequest := resty.R().
		SetAuthToken(token)
	return sendTestReport(reportFile, options, reportType, partialRequest)
}

func SendTestReport(reportFile *os.File, options Options, reportType string) error {
	partialRequest := resty.R()
	return sendTestReport(reportFile, options, reportType, partialRequest)
}

func sendTestReport(reportFile *os.File, options Options, reportType string, partialRequest *resty.Request) error {

	// we can remove file after sent it to test hub
	defer os.Remove(reportFile.Name())

	dat, err := ioutil.ReadFile(reportFile.Name())

	if err != nil {
		return err
	}

	server := options.URL
	project := options.Project
	build := options.Build

	err = configureHttps(options)

	if err != nil {
		return err
	}

	Debug("Sending report to %s/%s/%s of type %s", server, project, build, reportType)

	partialRequest.
		SetHeader("Content-Type", "application/gzip").
		SetHeader("x-testhub-type", reportType)

	resp, err := partialRequest.
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

func configureHttps(options Options) error {
	clientConfig := tls.Config{}

	if options.IsRootCaSet() {
		caCert, err := ioutil.ReadFile(options.RootCA)
		if err != nil {
			return err
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		clientConfig.RootCAs = caCertPool
	}

	if options.IsCertFile() && options.IsKeyFileSet() {
		cert, err := tls.LoadX509KeyPair(options.CertFile, options.KeyFile)
		if err != nil {
			return err
		}

		clientConfig.Certificates = []tls.Certificate{cert}
	}

	if options.SkipVerify {
		clientConfig.InsecureSkipVerify = true
	}

	resty.SetTLSClientConfig(&clientConfig)

	return nil
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

	if options.IsReportTestTypeHomeSet() {
		queryParams["homePage"] = options.ReportTestType.Home
	}

	return queryParams
}

func DeleteBuildWithToken(options Options, token string) error {
	partialRequest := resty.R().
		SetAuthToken(token)
	return deleteBuild(options, partialRequest)
}

func DeleteBuild(options Options) error {
	partialRequest := resty.R()
	return deleteBuild(options, partialRequest)
}

func deleteBuild(options Options, partialRequest *resty.Request) error {

	server := options.URL
	project := options.Project
	build := options.Build

	Debug("Deleting Build %s/%s/%s", server, project, build)

	partialRequest.
		SetPathParams(map[string]string{
			"project": project,
			"build":   build,
		})

	resp, err := partialRequest.
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
