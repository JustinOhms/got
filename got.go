// got
package main

import (
	"fmt"
	"strconv"

	"os"

	"github.com/justinohms/got/gittools"
	"github.com/justinohms/got/gittools/gittest"
)

func main() {
	put("Got 4 Hello World!")
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	//gitRootCommand = "git rev-parse --show-toplevel"

	fmt.Println(pwd)

	match, _, _, err := gittest.IsInsideAGitRepository(pwd, pwd)

	haltOnError(err)

	if !match {
		//actually don't need this should just pass through to git
		fmt.Println("fatal: Not a git repository (or any of the parent directories): .git")
	} else {
		fmt.Println("IS Repo:" + strconv.FormatBool(match))
	}

	reporoot, err := git.RepoRoot(pwd)

	fmt.Println("Root:" + reporoot)

	//	reporoot := gitRepoRootDirPath(pwd) + "/"
	//	fmt.Println("RR:" + reporoot)

	//	repo, err := git.OpenRepository(string(reporoot))

	haltOnError(err)
	//	put(repo.Path)

}

//func gitRepoRootDirPath(pwd string) string {
//	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
//	cmd.Dir = pwd

//	output, err := cmd.CombinedOutput()

//	haltOnError(err)

//	//return string(output[0 : len(output)-1])
//	return strings.TrimSpace(string(output))
//}

func haltOnError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func put(msg string) {
	fmt.Println(msg)
}

func putBytes(outs []byte) {
	if len(outs) > 0 {
		fmt.Printf("==> Output: %s\n", string(outs))
	}
}
