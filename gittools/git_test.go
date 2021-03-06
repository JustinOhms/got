package git

import (
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/ghthor/gospec"
	. "github.com/ghthor/gospec"
	. "github.com/ghthor/journal/git/gittest"
)

func TestIntegrationSpecs(t *testing.T) {
	r := gospec.NewRunner()

	r.AddSpec(DescribeGitIntegration)
	r.AddSpec(DescribeCommit)

	gospec.MainGoTest(r, t)
}

func DescribeGitIntegration(c gospec.Context) {
	c.Specify("a git repository will be created", func() {
		d, err := ioutil.TempDir("", "git_integration_test")
		c.Assume(err, IsNil)

		defer func(dir string) {
			c.Expect(os.RemoveAll(dir), IsNil)
		}(d)

		d = path.Join(d, "a_git_repo")
		c.Expect(Init(d), IsNil)
		c.Expect(d, IsAGitRepository)

		c.Specify("and will be clean", func() {
			c.Expect(IsClean(d), IsNil)
		})

		testFile := path.Join(d, "test_file")
		c.Assume(ioutil.WriteFile(testFile, []byte("some data\n"), 0666), IsNil)

		c.Specify("and will be dirty", func() {
			c.Expect(IsClean(d).Error(), Equals, "directory is dirty")
		})

		sd := d + "/test_sub_dir"
		err = os.Mkdir(sd, os.FileMode(0777))
		c.Assume(err, IsNil)

		defer func(dir string) {
			c.Expect(os.RemoveAll(dir), IsNil)
		}(sd)

		c.Specify("and the sub directory's repo root is the root ", func() {
			rr, _ := RepoRoot(sd)
			//we must strip off "/private" if it is there it's just a quirk of where the temp fs is
			rr = strings.TrimSpace(strings.Replace(rr, "/private", "", 1))
			c.Expect(rr, Equals, d)
		})

		c.Specify("and will add a file", func() {
			c.Expect(AddFilepath(d, testFile), IsNil)
			o, err := Command(d, "status", "-s").Output()
			c.Assume(err, IsNil)
			c.Expect(string(o), Equals, "A  test_file\n")
		})

		c.Specify("and will commit all staged changes", func() {
			c.Assume(AddFilepath(d, testFile), IsNil)
			c.Expect(CommitWithMessage(d, "a commit msg"), IsNil)

			o, err := Command(d, "show", "--no-color", "--pretty=format:\"%s%b\"").Output()
			c.Assume(err, IsNil)

			c.Expect(string(o), Equals, `"a commit msg"
diff --git a/test_file b/test_file
new file mode 100644
index 0000000..4268632
--- /dev/null
+++ b/test_file
@@ -0,0 +1 @@
+some data
`)
		})

		c.Specify("and will create an empty commit with message", func() {
			c.Expect(CommitEmpty(d, "an empty commit"), IsNil)

			o, err := Command(d, "show", "--no-color", "--pretty=format:\"%s%b\"").Output()
			c.Assume(err, IsNil)

			c.Expect(string(o), Equals, `"an empty commit"`)
		})
	})
}
