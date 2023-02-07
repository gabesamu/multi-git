package repo_manager

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type RepoManager struct {
	repos []string
	ignoreErrors bool
}

func NewRepoManager(baseDir string,
										repoNames []string,
										ignoreErrors bool) (repoManager *RepoManager, err error) {

	_, err = os.Stat(baseDir)
	if err != nil {
		if os.IsNotExist(err) {
			err = errors.New(fmt.Sprintf("base dir: '%s' doesn't exist", baseDir))
		}
		return
	}

	if len(repoNames) == 0 {
		err = errors.New(fmt.Sprintf("repo list can't be empty"))
		return
	}

	if baseDir[len(baseDir) - 1] != '/' {
		baseDir += "/"
	}

	repoManager = &RepoManager{
		ignoreErrors: ignoreErrors,
	}

	for _, repo := range repoNames {
		path := baseDir + repo
		repoManager.repos = append(repoManager.repos, path)
	}

	return
}

func (m *RepoManager) GetRepos() ([]string) {
	return m.repos
}

func (m *RepoManager) Exec(cmd string) (output map[string]string, err error) {
	var (
		start int
		end int
		components []string
		insentence bool
	)

	for i, component := range strings.Split(cmd, " ") {
		if insentence == true {
			if !strings.HasSuffix(component, "\""){
				continue
			}
			insentence = false
			end = i
			component = cmd[start:end]
		}

		if strings.HasPrefix(component, "\"") {
			insentence = true
			start = i
			continue
		}

		components = append(components, component)
	}

	// change back to pwd after commands are ran
	wd, _ := os.Getwd()
	defer os.Chdir(wd)

	var out []byte

	for _, repo := range m.repos {
		os.Chdir(repo)

		out, err = exec.Command("git", components...).CombinedOutput()

		output[repo] = string(out)

		if err != nil && !m.ignoreErrors {
			return
		}
	}
	return
}
