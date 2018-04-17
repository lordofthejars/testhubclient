package hub

type Options struct {
	URL        string
	Project    string
	Build      string
	BuildURL   string
	RepoURL    string
	Commit     string
	Branch     string
	RepoType   string
	RootCA     string
	CertFile   string
	KeyFile    string
	SkipVerify bool
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

func (o Options) IsRootCaSet() bool {
	return len(o.RootCA) > 0
}

func (o Options) IsCertFile() bool {
	return len(o.CertFile) > 0
}

func (o Options) IsKeyFileSet() bool {
	return len(o.KeyFile) > 0
}

func PublishTestReport(options Options, reportDirectory string) {

	zippedResults, error := CreateTestReportsPackage(".", "surefire.tar.gz", reportDirectory)

	if error != nil {
		Error("Couldn't create gzip file with test reports because of %s", error.Error())
	}

	applyDefaults(&options)
	error = SendReport(zippedResults, options, "surefire")

	if error != nil {
		Error("Couldn't send gzip file with test report because of %s", error.Error())
	}

}

func applyDefaults(options *Options) {
	if !options.IsBranchSet() {
		branch, error := getCurrentBranch()

		if error != nil {
			Info("Branch has not been possible to get from repo because of %s", error.Error())
		}

		options.Branch = branch
	}

	if !options.IsCommitSet() {
		commit, error := getCurrentRevision()

		if error != nil {
			Info("Commit Id has not been possible to get from repo because of %s", error.Error())
		}

		options.Commit = commit

	}

}

func RemoveBuild(options Options) {
	error := DeleteBuild(options)

	if error != nil {
		Error("Couldn't delete project %s build %s because of %s", options.Project, options.Build, error.Error())
	}
}
