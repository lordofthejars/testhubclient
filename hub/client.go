package hub

type Options struct {
	URL     string
	Project string
	Build   string
}

func PublishTestReport(options Options, reportDirectory string) {

	zippedResults, error := CreateTestReportsPackage(".", "surefire.tar.gz", reportDirectory)

	if error != nil {
		Error("Couldn't create gzip file with test reports because of %s", error.Error())
	}

	error = SendReport(zippedResults, options.URL, options.Project, options.Build, "surefire")

	if error != nil {
		Error("Couldn't send gzip file with test report because of %s", error.Error())
	}

}

func RemoveBuild(options Options) {
	error := DeleteBuild(options.URL, options.Project, options.Build)

	if error != nil {
		Error("Couldn't delete project %s build %s because of %s", options.Project, options.Build, error.Error())
	}
}
