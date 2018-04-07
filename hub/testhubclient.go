package hub

type Options struct {
	URL      string
	Project  string
	Build    string
	BuildURL string
	RepoURL  string
	Commit   string
	Branch   string
	RepoType string
}

func (o Options) IsBuildUrlSet() bool {
	return len(o.BuildURL) > 0
}

func (o Options) IsCommitSet() bool {
	return len(o.Commit) > 0
}

func (o Options) IsBranchSet() bool {
	return len(o.Branch) > 0
}

func (o Options) IsRepoUrlSet() bool {
	return len(o.RepoURL) > 0
}

func (o Options) IsRepoTypeSet() bool {
	return len(o.RepoType) > 0
}

func PublishTestReport(options Options, reportDirectory string) {

	zippedResults, error := CreateTestReportsPackage(".", "surefire.tar.gz", reportDirectory)

	if error != nil {
		Error("Couldn't create gzip file with test reports because of %s", error.Error())
	}

	error = SendReport(zippedResults, options, "surefire")

	if error != nil {
		Error("Couldn't send gzip file with test report because of %s", error.Error())
	}

}

func RemoveBuild(options Options) {
	error := DeleteBuild(options)

	if error != nil {
		Error("Couldn't delete project %s build %s because of %s", options.Project, options.Build, error.Error())
	}
}
