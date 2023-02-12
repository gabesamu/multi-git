package helpers

import (
	"os"
	"os/exec"
	"path"
)

func CreateDir(baseDir string, name string, initGit bool) (err error) {
	path := path.Join(baseDir, name)
	err = os.MkdirAll(path, os.ModePerm)

	if err != nil {
		return
	}

	if !initGit {
		return
	}

	wd, _ := os.Getwd()
	defer os.Chdir(wd)
	os.Chdir(path)

	err = exec.Command("git", "init").Run()
	return
}

func AddFiles(baseDir string,
							dirName string,
							commit bool,
							filesNames ...string) (err error) {
	dir := path.Join(baseDir, dirName)
	for _, file := range filesNames {
		data := []byte("File data for: " + file)
		err = os.WriteFile(path.Join(dir, file), data, os.ModePerm)
		if err != nil {
			return
		}
	}

	if !commit {
		return
	}

	wd, _ := os.Getwd()
	defer os.Chdir(wd)

	os.Chdir(dir)
	err = exec.Command("git", "add", "-A").Run()
	if err != nil {
		return
	}
	err = exec.Command("git", "commit", "-m", "added files").Run()
	return
}
