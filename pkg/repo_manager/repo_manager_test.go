package repo_manager_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"fmt"
	"os"
	"path"
	"strings"

	"github.com/0xlucius/multi-git/pkg/helpers"
	. "github.com/0xlucius/multi-git/pkg/repo_manager"
)

var _ = Describe("RepoManager", func() {
	const baseDir = "/tmp/test-multi-git"
	var repoList []string

	removeTestDirs := func() {
		err := os.RemoveAll(baseDir)
		Expect(err).To(BeNil())
	}

	BeforeEach(func() {
		err := helpers.CreateDir(baseDir, "test-dir1", true)
		Expect(err).To(BeNil())
		repoList = []string{"test-dir-1"}
	})

	AfterEach(removeTestDirs)

	Describe("Initializing a new repoManager", func () {
		It("Should fail with invalid base dir", func() {
			_, err := NewRepoManager("/no-such-dir", repoList, true)
			Expect(err).ToNot(BeNil())
		})

		It("Should fail with empty repo list", func() {
			_, err := NewRepoManager(baseDir, []string{}, true)
			Expect(err).ToNot(BeNil())
		})
	})

	Describe("Get repos", func ()  {
		It("Should get repo list successfully", func() {
			rm, _ := NewRepoManager(baseDir, repoList, true)

			repos := rm.GetRepos()
			Expect(repos).To(HaveLen(1))
			Expect(repos[0] == path.Join(baseDir, repoList[0])).To(BeTrue())
		})
	})

	Describe("Executing Git commands", func ()  {
		It("Should create branches successfully", func() {
			repoList = append(repoList, "dir-2")
			helpers.CreateDir(baseDir, repoList[1], true)
			rm, _ := NewRepoManager(baseDir, repoList, true)

			output, err := rm.Exec("checkout -b test-branch")
			Expect(err).To(BeNil())

			for _, out := range output {
				Expect(out).To(Equal("Switched to a new branch 'test-branch'\n"))
			}
		})

		It("Should commit files successfully", func() {
			rm, _ := NewRepoManager(baseDir, repoList, true)

			output, err := rm.Exec("checkout -b test-branch")
			Expect(err).To(BeNil())

			for _, out := range output {
				Expect(out).To(Equal("Switched to a new branch 'test-branch'\n"))
			}

			err = helpers.AddFiles(baseDir, repoList[0], true, "file_1.txt", "file_2.txt")
			Expect(err).To(BeNil())

			// Restore working directory after executing the command
			wd, _ := os.Getwd()
			defer os.Chdir(wd)

			dir := path.Join(baseDir, repoList[0])
			err = os.Chdir(dir)
			Expect(err).To(BeNil())

			output, err = rm.Exec("log --oneline")
			fmt.Println(output)
			Expect(err).To(BeNil())

			ok := strings.HasSuffix(output[dir], "added some files...\n")
			Expect(ok).To(BeTrue())
		})
	})



})
