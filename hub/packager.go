package hub

import (
	"archive/tar"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/mholt/archiver"
	"github.com/yargevad/filepathx"
)

type InvalidLocation struct {
	Location string
}

func (e *InvalidLocation) Error() string {
	return fmt.Sprintf("%s path does not exists", e.Location)
}

func addFile(tw *tar.Writer, path string, testdirectory string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	if stat, err := file.Stat(); err == nil {
		fileLocationGzip := strings.Trim(path, testdirectory)
		// now lets create the header as needed for this file within the tarball
		header := new(tar.Header)
		header.Name = fileLocationGzip
		header.Size = stat.Size()
		header.Mode = int64(stat.Mode())
		header.ModTime = stat.ModTime()
		header.Typeflag = tar.TypeReg
		// write the header to the tarball archive
		if err := tw.WriteHeader(header); err != nil {
			return err
		}
		// copy the file data to the tarball
		if _, err := io.Copy(tw, file); err != nil {
			return err
		}
	}
	return nil
}

func getReportLocation(original string, testdirectory string) string {
	return original
}

func CreateTestReportsPackage(rootdir string, outputFileName string, testdirectory string) (*os.File, error) {

	if exists(rootdir) {

		if !strings.HasSuffix(testdirectory, "/") {
			testdirectory = testdirectory + "/"
		}

		expression := "/**/" + testdirectory + "*.xml"
		pathExpression := rootdir + expression
		reports, err := filepathx.Glob(pathExpression)

		if err != nil {
			return nil, err
		}

		tmpDir, err := ioutil.TempDir("", "testhub")

		if err != nil {
			return nil, err
		}

		output := filepath.Join(tmpDir, outputFileName)
		return compress(output, reports)
	}
	return nil, &InvalidLocation{rootdir}

}

func CreateReportPackage(dir string, outputFileName string) (*os.File, error) {
	if exists(dir) {

		tmpDir, err := ioutil.TempDir("", "testhub")

		if err != nil {
			return nil, err
		}

		output := filepath.Join(tmpDir, outputFileName)
		return compress(output, []string{dir})
	}
	return nil, &InvalidLocation{dir}
}

func exists(path string) bool {

	if _, err := os.Stat(path); err == nil {
		return true
	}

	// Maybe does not exists or there are any permission problem but in anyeway for us is that it does not exists
	return false
}

func compress(output string, paths []string) (*os.File, error) {
	// set up the output file
	Debug("Creating temporary gzipped file at %s", output)
	file, err := os.Create(output)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	// set up the gzip writer
	err = archiver.TarGz.Make(output, paths)

	return file, err
}
