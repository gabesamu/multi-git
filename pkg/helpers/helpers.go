package helpers

import (
	"fmt"
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

func RunMultiGit(command string, ignoreErrors bool, mgRoot string, mgRepos string, useConfigFile bool) (output string, err error) {
	out, err := exec.Command("which", "multi-git").CombinedOutput()
	if err != nil {
		return
	}

	if len(out) == 0 {
		err = fmt.Errorf("multi-git is not in the PATH")
		return
	}

	components := []string{command}
	env := os.Environ()
	if useConfigFile {
		configFile := path.Join(mgRoot, "multi-git-test-config.toml")
		data := fmt.Sprintf("root = \"%s\"\nrepos = \"%s\"\nignore-errors = %v\n", mgRoot, mgRepos, ignoreErrors)
		err = os.WriteFile(configFile, []byte(data), 0644)
		if err != nil {
			return
		}
		components = append(components, "--config", configFile)
	} else {
		if ignoreErrors {
			components = append(components, "--ignore-errors")
		}
		env = append(env, "MG_ROOT="+mgRoot, "MG_REPOS="+mgRepos)
	}

	cmd := exec.Command("multi-git", components...)
	cmd.Env = env
	out, err = cmd.CombinedOutput()
	output = string(out)
	return
}
