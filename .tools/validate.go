package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

// DefaultRules are the standard validation to perform on git commits
var DefaultRules = []ValidateRule{
	func(c CommitEntry) (vr ValidateResult) {
		vr.CommitEntry = c
		if len(strings.Split(c["parent"], " ")) > 1 {
			vr.Pass = true
			vr.Msg = "merge commits do not require DCO"
			return vr
		}

		hasValid := false
		for _, line := range strings.Split(c["body"], "\n") {
			if validDCO.MatchString(line) {
				hasValid = true
			}
		}
		if !hasValid {
			vr.Pass = false
			vr.Msg = "does not have a valid DCO"
		} else {
			vr.Pass = true
			vr.Msg = "has a valid DCO"
		}

		return vr
	},
	// TODO add something for the cleanliness of the c.Subject
	func(c CommitEntry) (vr ValidateResult) {
		vr.CommitEntry = c
		buf := bytes.NewBuffer([]byte{})
		args := []string{"git", "show", "--check", c["commit"]}
		vr.Msg = strings.Join(args, " ")
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stdout = buf
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			vr.Pass = false
			vr.Detail = string(buf.Bytes())
			return vr
		}
		vr.Pass = true
		return vr
	},
	func(c CommitEntry) (vr ValidateResult) {
		return ExecTree(c, "go", "vet", "./...")
	},
	func(c CommitEntry) (vr ValidateResult) {
		return ExecTree(c, "go", "fmt", "./...")
	},
	func(c CommitEntry) (vr ValidateResult) {
		vr = ExecTree(c, os.ExpandEnv("$HOME/gopath/bin/golint"), "./...")
		if len(vr.Detail) > 0 {
			vr.Pass = false
		}
		return vr
	},
}

var (
	flVerbose     = flag.Bool("v", false, "verbose")
	flCommitRange = flag.String("range", "", "use this commit range instead")

	validDCO = regexp.MustCompile(`^Signed-off-by: ([^<]+) <([^<>@]+@[^<>]+)>$`)
)

func main() {
	flag.Parse()

	var commitrange string
	if *flCommitRange != "" {
		commitrange = *flCommitRange
	} else {
		var err error
		commitrange, err = GitFetchHeadCommit()
		if err != nil {
			log.Fatal(err)
		}
	}

	c, err := GitCommits(commitrange)
	if err != nil {
		log.Fatal(err)
	}

	results := ValidateResults{}

	if *flVerbose {
		fmt.Println("TAP version 13")
		fmt.Printf("1..%d\n", len(c)*len(DefaultRules))
	}
	for i, commit := range c {
		fmt.Printf("# %s %s ... ", commit["abbreviated_commit"], commit["subject"])
		vr := ValidateCommit(commit, DefaultRules)
		results = append(results, vr...)
		if _, fail := vr.PassFail(); fail == 0 {
			fmt.Println("PASS")
		} else {
			fmt.Println("FAIL")
		}
		for j, r := range vr {
			if *flVerbose {
				if r.Pass {
					fmt.Printf("ok")
				} else {
					fmt.Printf("not ok")
				}
				fmt.Printf(" %d - %s\n", i*len(DefaultRules)+j+1, r.Msg)
			} else if !r.Pass {
				fmt.Printf("not ok - %s\n", r.Msg)
			}
			if (*flVerbose || !r.Pass) && len(r.Detail) > 0 {
				m := map[string]string{"message": r.Detail}
				buf, err := yaml.Marshal(m)
				if err != nil {
					log.Fatal(err)
				}
				lines := strings.Split(strings.TrimSpace(string(buf)), "\n")
				fmt.Println(" ---")
				for _, line := range lines {
					fmt.Printf(" %s\n", line)
				}
				fmt.Println(" ...")
			}
		}
	}
	_, fail := results.PassFail()
	if fail > 0 {
		fmt.Printf("%d issues to fix\n", fail)
		os.Exit(1)
	}
}

// ValidateRule will operate over a provided CommitEntry, and return a result.
type ValidateRule func(CommitEntry) ValidateResult

// ValidateCommit processes the given rules on the provided commit, and returns the result set.
func ValidateCommit(c CommitEntry, rules []ValidateRule) ValidateResults {
	results := ValidateResults{}
	for _, r := range rules {
		results = append(results, r(c))
	}
	return results
}

// ValidateResult is the result for a single validation of a commit.
type ValidateResult struct {
	CommitEntry CommitEntry
	Pass        bool
	Msg         string
	Detail      string
}

// ValidateResults is a set of results. This is type makes it easy for the following function.
type ValidateResults []ValidateResult

// PassFail gives a quick over/under of passes and failures of the results in this set
func (vr ValidateResults) PassFail() (pass int, fail int) {
	for _, res := range vr {
		if res.Pass {
			pass++
		} else {
			fail++
		}
	}
	return pass, fail
}

// CommitEntry represents a single commit's information from `git`
type CommitEntry map[string]string

var (
	prettyFormat         = `--pretty=format:`
	formatSubject        = `%s`
	formatBody           = `%b`
	formatCommit         = `%H`
	formatAuthorName     = `%aN`
	formatAuthorEmail    = `%aE`
	formatCommitterName  = `%cN`
	formatCommitterEmail = `%cE`
	formatSigner         = `%GS`
	formatCommitNotes    = `%N`
	formatMap            = `{"commit": "%H", "abbreviated_commit": "%h", "tree": "%T", "abbreviated_tree": "%t", "parent": "%P", "abbreviated_parent": "%p", "refs": "%D", "encoding": "%e", "sanitized_subject_line": "%f", "verification_flag": "%G?", "signer_key": "%GK", "author_date": "%aD" , "committer_date": "%cD" }`
)

// GitLogCommit assembles the full information on a commit from its commit hash
func GitLogCommit(commit string) (*CommitEntry, error) {
	buf := bytes.NewBuffer([]byte{})
	cmd := exec.Command("git", "log", "-1", prettyFormat+formatMap, commit)
	cmd.Stdout = buf
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Println(strings.Join(cmd.Args, " "))
		return nil, err
	}
	c := CommitEntry{}
	output := buf.Bytes()
	if err := json.Unmarshal(output, &c); err != nil {
		fmt.Println(string(output))
		return nil, err
	}

	// any user provided fields can't be sanitized for the mock-json marshal above
	for k, v := range map[string]string{
		"subject":         formatSubject,
		"body":            formatBody,
		"author_name":     formatAuthorName,
		"author_email":    formatAuthorEmail,
		"committer_name":  formatCommitterName,
		"committer_email": formatCommitterEmail,
		"commit_notes":    formatCommitNotes,
		"signer":          formatSigner,
	} {
		output, err := exec.Command("git", "log", "-1", prettyFormat+v, commit).Output()
		if err != nil {
			return nil, err
		}
		c[k] = strings.TrimSpace(string(output))
	}

	return &c, nil
}

// GitCommits returns a set of commits.
// If commitrange is a git still range 12345...54321, then it will be isolated set of commits.
// If commitrange is a single commit, all ancestor commits up through the hash provided.
func GitCommits(commitrange string) ([]CommitEntry, error) {
	output, err := exec.Command("git", "log", prettyFormat+formatCommit, commitrange).Output()
	if err != nil {
		return nil, err
	}
	commitHashes := strings.Split(strings.TrimSpace(string(output)), "\n")
	commits := make([]CommitEntry, len(commitHashes))
	for i, commitHash := range commitHashes {
		c, err := GitLogCommit(commitHash)
		if err != nil {
			return commits, err
		}
		commits[i] = *c
	}
	return commits, nil
}

// GitFetchHeadCommit returns the hash of FETCH_HEAD
func GitFetchHeadCommit() (string, error) {
	output, err := exec.Command("git", "rev-parse", "--verify", "FETCH_HEAD").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// GitHeadCommit returns the hash of HEAD
func GitHeadCommit() (string, error) {
	output, err := exec.Command("git", "rev-parse", "--verify", "HEAD").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// GitCheckoutTree extracts the tree associated with the given commit
// to the given directory.  Unlike 'git checkout ...', it does not
// alter the HEAD.
func GitCheckoutTree(commit string, directory string) error {
	pipeReader, pipeWriter := io.Pipe()
	gitCmd := exec.Command("git", "archive", commit)
	gitCmd.Stdout = pipeWriter
	gitCmd.Stderr = os.Stderr
	tarCmd := exec.Command("tar", "-xC", directory)
	tarCmd.Stdin = pipeReader
	tarCmd.Stderr = os.Stderr
	err := gitCmd.Start()
	if err != nil {
		return err
	}
	defer gitCmd.Process.Kill()
	err = tarCmd.Start()
	if err != nil {
		return err
	}
	defer tarCmd.Process.Kill()
	err = gitCmd.Wait()
	if err != nil {
		return err
	}
	err = pipeWriter.Close()
	if err != nil {
		return err
	}
	err = tarCmd.Wait()
	if err != nil {
		return err
	}
	return nil
}

// ExecTree executes a command in a checkout of the commit's tree,
// wrapping any errors in a ValidateResult object.
func ExecTree(c CommitEntry, args ...string) (vr ValidateResult) {
	vr.CommitEntry = c
	stdout, err := execTree(c, args...)
	vr.Detail = strings.TrimSpace(stdout)
	if err == nil {
		vr.Pass = true
		vr.Msg = strings.Join(args, " ")
	} else {
		vr.Pass = false
		vr.Msg = fmt.Sprintf("%s : %s", strings.Join(args, " "), err.Error())
	}
	return vr
}

// execTree executes a command in a checkout of the commit's tree
func execTree(c CommitEntry, args ...string) (string, error) {
	dir, err := ioutil.TempDir("", "go-validate-")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(dir)
	err = GitCheckoutTree(c["commit"], dir)
	if err != nil {
		return "", err
	}
	buf := bytes.NewBuffer([]byte{})
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = dir
	cmd.Stdout = buf
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	stdout := string(buf.Bytes())
	if err != nil {
		return stdout, err
	}
	return stdout, nil
}
