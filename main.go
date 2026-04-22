package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

func main() {
	// Check if git is available
	if _, err := exec.LookPath("git"); err != nil {
		log.Fatal("git not found in PATH")
	}

	// Run git diff --staged --name-only
	cmd := exec.Command("git", "diff", "--staged", "--name-only")
	output, err := cmd.Output()
	if err != nil {
		log.Fatal("error running git diff --name-only: ", err)
	}

	filesStr := strings.TrimSpace(string(output))
	if filesStr == "" {
		log.Fatal("no staged files to commit")
	}
	files := strings.Split(filesStr, "\n")

	// Run git diff --staged
	cmd = exec.Command("git", "diff", "--staged")
	diffOutput, err := cmd.Output()
	if err != nil {
		log.Fatal("error running git diff: ", err)
	}
	diff := string(diffOutput)

	// Check if this is the first commit
	cmd = exec.Command("git", "log", "--oneline", "-1")
	_, err = cmd.Output()
	isFirstCommit := err != nil

	var lastCommitType string
	if !isFirstCommit {
		// Get last commit message to infer type
		cmd = exec.Command("git", "log", "--oneline", "-1")
		lastOutput, _ := cmd.Output()
		lastMsg := strings.TrimSpace(string(lastOutput))
		if strings.Contains(lastMsg, ": ") {
			lastCommitType = strings.Split(lastMsg, ": ")[0]
		}
	}

	// Determine commit type based on rules
	commitType := determineCommitType(files, diff, isFirstCommit)

	// Generate 3 unique suggestions
	suggestions := generateSuggestions(commitType, files, diff, lastCommitType)

	// Print colored suggestions
	cyan := color.New(color.FgCyan)
	for i, sug := range suggestions {
		cyan.Printf("%d. %s\n", i+1, sug)
	}

	// Read user input
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Choose 1-3 or 'q' to quit: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "q" {
		fmt.Println("Exiting without committing.")
		return
	}

	num, err := strconv.Atoi(input)
	if err != nil || num < 1 || num > 3 {
		log.Fatal("invalid input: choose 1, 2, 3, or q")
	}

	msg := suggestions[num-1]

	// Execute git commit
	cmd = exec.Command("git", "commit", "-m", msg)
	err = cmd.Run()
	if err != nil {
		log.Fatal("commit failed: ", err)
	}

	fmt.Println("Committed successfully with message:", msg)
}

func determineCommitType(files []string, diff string, isFirstCommit bool) string {
	if isFirstCommit {
		return "initial"
	}

	hasNew := false
	hasDelete := false
	hasTest := false
	hasDocs := false
	hasDepAdd := false
	hasDepRemove := false
	hasEnv := false
	hasVersion := false
	hasStyle := false
	hasRefactor := false

	// Analyze files
	for _, file := range files {
		lowerFile := strings.ToLower(file)
		if strings.Contains(lowerFile, "test") || strings.HasSuffix(lowerFile, "_test.go") || strings.HasSuffix(lowerFile, ".spec.js") {
			hasTest = true
		}
		if strings.Contains(lowerFile, "readme") || strings.HasSuffix(lowerFile, ".md") || strings.HasSuffix(lowerFile, ".txt") {
			hasDocs = true
		}
		if strings.Contains(lowerFile, "package.json") || strings.Contains(lowerFile, "version") || strings.Contains(lowerFile, "cargo.toml") {
			hasVersion = true
		}
		if strings.Contains(lowerFile, ".env") || strings.Contains(lowerFile, "config") || strings.Contains(lowerFile, "settings") {
			hasEnv = true
		}
		if strings.Contains(lowerFile, "go.mod") || strings.Contains(lowerFile, "requirements.txt") || strings.Contains(lowerFile, "package-lock.json") {
			// Check for additions/removals in diff
			if strings.Contains(diff, "+") && strings.Contains(diff, "-") {
				hasRefactor = true
			} else if strings.Count(diff, "+") > strings.Count(diff, "-") {
				hasDepAdd = true
			} else if strings.Count(diff, "-") > strings.Count(diff, "+") {
				hasDepRemove = true
			}
		}
	}

	// Analyze diff for new/delete files
	lines := strings.Split(diff, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "+++ ") && strings.Contains(line, "/dev/null") {
			hasDelete = true
		} else if strings.HasPrefix(line, "+++ ") {
			hasNew = true
		}
	}

	// Check for style changes (only whitespace)
	addLines := 0
	delLines := 0
	for _, line := range lines {
		if strings.HasPrefix(line, "+") && strings.TrimSpace(line[1:]) == "" {
			addLines++
		}
		if strings.HasPrefix(line, "-") && strings.TrimSpace(line[1:]) == "" {
			delLines++
		}
	}
	if addLines > 0 && delLines > 0 && addLines == delLines {
		hasStyle = true
	}

	// Folder based type
	folderBased := true
	for _, file := range files {
		if !strings.HasPrefix(file, "src/") && !strings.HasPrefix(file, "lib/") && !strings.HasPrefix(file, "app/") {
			folderBased = false
			break
		}
	}

	// Apply rules in order
	if hasVersion {
		return "version"
	}
	if hasNew && !hasDelete {
		return "feat"
	}
	if hasDelete && !hasNew {
		return "feat"
	}
	if hasTest && len(files) == 1 {
		return "test"
	}
	if hasDocs && len(files) == 1 {
		return "docs"
	}
	if hasRefactor {
		return "refactor"
	}
	if hasStyle {
		return "style"
	}
	if hasDepAdd {
		return "deps"
	}
	if hasDepRemove {
		return "deps"
	}
	if hasEnv {
		return "env"
	}
	if folderBased {
		return "feat"
	}
	return "feat" // fallback
}

func generateSuggestions(commitType string, files []string, diff string, lastCommitType string) []string {
	// If last commit type is known, prefer similar or varied
	preferred := []string{"feat", "fix", "refactor", "docs", "style", "test", "chore"}
	if lastCommitType != "" && contains(preferred, lastCommitType) {
		// Rotate to start with last type
		idx := indexOf(preferred, lastCommitType)
		preferred = append(preferred[idx:], preferred[:idx]...)
	}

	switch commitType {
	case "initial":
		return []string{
			"chore: initial commit",
			"feat: initial implementation",
			"docs: project setup",
		}
	case "version":
		return []string{
			"bump: version to latest",
			"release: bump version number",
			"chore: update version",
		}
	case "test":
		return []string{
			"test: add unit tests",
			"test: update test cases",
			"test: fix failing tests",
		}
	case "docs":
		return []string{
			"docs: update documentation",
			"docs: add README section",
			"docs: fix documentation errors",
		}
	case "refactor":
		return []string{
			"refactor: restructure code",
			"refactor: rename variables",
			"refactor: optimize performance",
		}
	case "style":
		return []string{
			"style: format code",
			"style: fix indentation",
			"style: update code style",
		}
	case "deps":
		return []string{
			"deps: add new dependency",
			"deps: update dependencies",
			"deps: remove unused dependency",
		}
	case "env":
		return []string{
			"env: update environment variables",
			"env: add config settings",
			"env: fix environment setup",
		}
	default: // feat or fallback
		// Analyze diff for additions and deletions
		addCount := strings.Count(diff, "\n+")
		delCount := strings.Count(diff, "\n-")

		// Check for specific keywords in diff to make smarter suggestions
		hasFirstCommit := strings.Contains(diff, "isFirstCommit") || strings.Contains(diff, "initial")
		hasHistory := strings.Contains(diff, "lastCommitType") || strings.Contains(diff, "history")
		hasSuggestions := strings.Contains(diff, "generateSuggestions") || strings.Contains(diff, "suggestion")

		if hasFirstCommit && hasHistory && hasSuggestions {
			return []string{
				"feat: enhance commit message generation with first commit detection and history-aware suggestions",
				"feat: implement intelligent commit type analysis based on git history",
				"feat: add context-aware suggestion system with initial commit handling",
			}
		}
		if hasFirstCommit {
			return []string{
				"feat: add first commit detection logic",
				"feat: implement initial commit message suggestions",
				"feat: enhance tool with repository initialization checks",
			}
		}
		if hasHistory {
			return []string{
				"feat: integrate git history analysis for consistent commit types",
				"feat: add last commit type detection",
				"feat: improve suggestions based on previous commits",
			}
		}
		if hasSuggestions {
			return []string{
				"feat: enhance suggestion generation logic",
				"feat: implement smarter commit message proposals",
				"feat: add dynamic suggestion system",
			}
		}

		// Based on diff size
		if addCount > delCount * 2 {
			return []string{
				"feat: add significant new functionality",
				"feat: implement major feature enhancements",
				"feat: expand application capabilities",
			}
		} else if delCount > addCount * 2 {
			return []string{
				"feat: remove obsolete code and simplify logic",
				"refactor: clean up and optimize code",
				"feat: streamline application structure",
			}
		} else if addCount > 0 && delCount > 0 {
			return []string{
				"refactor: modify and improve existing code",
				"feat: update functionality with changes",
				"feat: enhance features through code modifications",
			}
		}

		// Make more descriptive based on files
		desc := "new functionality"
		if len(files) == 1 {
			if strings.Contains(files[0], "main") {
				desc = "main application logic"
			} else if strings.Contains(files[0], "test") {
				desc = "test coverage"
			} else if strings.Contains(files[0], "readme") {
				desc = "documentation"
			}
		} else if len(files) > 1 {
			desc = "multiple files"
		}
		return []string{
			fmt.Sprintf("%s: implement %s", preferred[0], desc),
			fmt.Sprintf("%s: resolve issue", preferred[1]),
			fmt.Sprintf("%s: add %s", preferred[0], desc),
		}
	}
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func indexOf(slice []string, item string) int {
	for i, s := range slice {
		if s == item {
			return i
		}
	}
	return -1
}