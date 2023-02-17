package e2e_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"os"
	"strings"

	"github.com/0xlucius/multi-git/pkg/helpers"
)

var _ = Describe("MultiGit", func() {
	const baseDir = "/tmp/multi-git"
	var repoList string
	var err error


	BeforeEach(func ()  {
		err := helpers.CreateDir(baseDir, "", false)
		Expect(err).To(BeNil())
	})

	AfterEach(func ()  {
		err := os.RemoveAll(baseDir)
		Expect(err).To(BeNil())
	})

	Context("When ran with empty/undefined evironment", func ()  {
		It("Should fail with empty base directory", func ()  {
			out, err := helpers.RunMultiGit("status", false, "/does-not-exist", repoList , false)
			Expect(err).ToNot(BeNil())
			errMessage := "base dir: '/does-not-exist/' doesn't exist\n"
			Expect(out).To(HaveSuffix(errMessage))
		})

		It("Should fail with empty repoList", func ()  {
			out, err := helpers.RunMultiGit("status", false, baseDir, repoList , false)
			Expect(err).ToNot(BeNil())
			Expect(out).To(ContainSubstring("repo list can't be empty"))
		})
	})

	Describe("Running Git commands", func ()  {
		Context("Given proper input", func() {
			It("Should do git init successfully", func() {
				err = helpers.CreateDir(baseDir, "dir-1", false)
				Expect(err).To(BeNil())
				err = helpers.CreateDir(baseDir, "dir-2", false)
				Expect(err).To(BeNil())
				repoList = "dir-1,dir-2"

				output, err := helpers.RunMultiGit("init", false, baseDir, repoList, false)
				Expect(err).To(BeNil())
				count := strings.Count(output, "Initialized empty Git repository")
				Expect(count).To(Equal(2))
			})

			It("Should do git status successfully for git directories", func() {
				err = helpers.CreateDir(baseDir, "dir-1", true)
				Expect(err).To(BeNil())
				err = helpers.CreateDir(baseDir, "dir-2", true)
				Expect(err).To(BeNil())
				repoList = "dir-1,dir-2"

				output, err := helpers.RunMultiGit("status", false, baseDir, repoList, false)
				Expect(err).To(BeNil())
				count := strings.Count(output, "nothing to commit")
				Expect(count).To(Equal(2))

				output, err = helpers.RunMultiGit("status", false, baseDir, repoList, false)
				Expect(err).To(BeNil())
				count = strings.Count(output, "nothing to commit")
				Expect(count).To(Equal(2))

			})

			It("Should create branches successfully", func() {
				err = helpers.CreateDir(baseDir, "dir-1", true)
				Expect(err).To(BeNil())
				err = helpers.CreateDir(baseDir, "dir-2", true)
				Expect(err).To(BeNil())
				repoList = "dir-1,dir-2"

				output, err := helpers.RunMultiGit("checkout -b test-branch", false, baseDir, repoList, false)
				Expect(err).To(BeNil())
				count := strings.Count(output, "Switched to a new branch 'test-branch'")
				Expect(count).To(Equal(2))

				output, err = helpers.RunMultiGit("checkout -b test-branch", false, baseDir, repoList, false)
				Expect(err).To(BeNil())
				count = strings.Count(output, "Switched to a new branch 'test-branch'")
				Expect(count).To(Equal(2))
			})
		})

		Context("When ran on non-git repo", func() {
			It("Should fail git status", func() {
				err = helpers.CreateDir(baseDir, "dir-1", false)
				Expect(err).To(BeNil())
				err = helpers.CreateDir(baseDir, "dir-2", false)
				Expect(err).To(BeNil())
				repoList = "dir-1,dir-2"

				output, err := helpers.RunMultiGit("status", false, baseDir, repoList, false)
				Expect(err).To(BeNil())
				Expect(output).To(ContainSubstring("fatal: not a git repository"))

				output, err = helpers.RunMultiGit("status", false, baseDir, repoList, false)
				Expect(err).To(BeNil())
				Expect(output).To(ContainSubstring("fatal: not a git repository"))

			})
		})
	})

	Describe("ignoreErrors flag", func ()  {
		Context("When a directory is invalid", func ()  {
			Context("ignoreErrors is true", func ()  {
				It("Should succeed for all valid directories", func ()  {
					err = helpers.CreateDir(baseDir, "dir-1", false)
					Expect(err).To(BeNil())
					err = helpers.CreateDir(baseDir, "dir-2", true)
					Expect(err).To(BeNil())
					repoList = "dir-1,dir-2"

					output, err := helpers.RunMultiGit("status", true, baseDir, repoList, false)
					Expect(err).To(BeNil())
					Expect(output).To(ContainSubstring("[dir-1]: git status\nfatal: not a git repository"))
					Expect(output).To(ContainSubstring("[dir-2]: git status\nOn branch main"))

					output, err = helpers.RunMultiGit("status", true, baseDir, repoList, false)
					Expect(err).To(BeNil())
					Expect(output).To(ContainSubstring("[dir-1]: git status\nfatal: not a git repository"))
					Expect(output).To(ContainSubstring("[dir-2]: git status\nOn branch main"))
				})
			})

			Context("ignoreErrors is false", func ()  {
				It("Should fail on first invalid directory", func ()  {
					err = helpers.CreateDir(baseDir, "dir-1", false)
					Expect(err).To(BeNil())
					err = helpers.CreateDir(baseDir, "dir-2", true)
					Expect(err).To(BeNil())
					repoList = "dir-1,dir-2"

					output, err := helpers.RunMultiGit("status", false, baseDir, repoList, false)
					Expect(err).To(BeNil())
					Expect(output).To(ContainSubstring("[dir-1]: git status\nfatal: not a git repository"))
					Expect(output).ShouldNot(ContainSubstring("[dir-2]"))

					output, err = helpers.RunMultiGit("status", false, baseDir, repoList, false)
					Expect(err).To(BeNil())
					Expect(output).To(ContainSubstring("[dir-1]: git status\nfatal: not a git repository"))
					Expect(output).ToNot(ContainSubstring("[dir-2]"))
				})
			})
		})
	})


})
