package hub

type Options struct {
	URL     string
	Project string
	Build   string
}

func PublishTestReport(options Options, reportDirectory string) {

	zippedResults, error := CreateTestReportsPackage(".", "tests.tar.gz", reportDirectory)

	if error != nil {
		Error("Couldn't create gzip file with test reports because of %s", error.Error())
	}

	error = SendReport(zippedResults, options.URL, options.Project, options.Build)

	if error != nil {
		Error("Couldn't send gzip file with test report because of %s", error.Error())
	}

}
