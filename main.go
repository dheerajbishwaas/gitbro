package main

import (
	"bufio"
	"fmt"
	"os" // <-- Ye line add kar
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

type Suggestion struct {
	Msg string
	Rule string
}

func main() {
	green := color.New(color.FgGreen)
	yellow := color.New(color.FgYellow)
	red := color.New(color.FgRed)
	cyan := color.New(color.FgCyan)
	hiBlack := color.New(color.FgHiBlack)

	if _, err := exec.LookPath("git"); err!= nil {
		red.Println("Git not found in PATH. Install from https://git-scm.com")
		os.Exit(1)
	}

	nameOnly, _ := runGit("diff", "--staged", "--name-only")
	filesStr := strings.TrimSpace(nameOnly)
	if filesStr == "" {
		yellow.Println("No staged files found. Run 'git add' first.")
		os.Exit(1)
	}
	files := strings.Split(filesStr, "\n")

	diff, _ := runGit("diff", "--staged")
	stat, _ := runGit("diff", "--staged", "--shortstat")
	status, _ := runGit("diff", "--staged", "--name-status")

	cyan.Println("Analyzing changes...")

	suggestions := []Suggestion{}
	seen := make(map[string]bool)

	add := func(msg, rule string) {
		msg = strings.ToLower(msg)
		if!seen[msg] && len(suggestions) < 3 {
			suggestions = append(suggestions, Suggestion{msg, rule})
			seen[msg] = true
		}
	}

	// Rule 1: Version bump
	for _, f := range files {
		if strings.HasSuffix(f, "package.json") {
			re := regexp.MustCompile(`\+\s*"version":\s*"(.*)"`)
			if m := re.FindStringSubmatch(diff); len(m) > 1 {
				add(fmt.Sprintf("chore: bump version to %s", m[1]), "version bump")
			}
		}
		if strings.HasSuffix(f, "Cargo.toml") {
			re := regexp.MustCompile(`\+\s*version\s*=\s*"(.*)"`)
			if m := re.FindStringSubmatch(diff); len(m) > 1 {
				add(fmt.Sprintf("chore: bump version to %s", m[1]), "version bump")
			}
		}
	}

	// Rule 2: New file
	for _, line := range strings.Split(status, "\n") {
		if strings.HasPrefix(line, "A\t") {
			f := strings.TrimPrefix(line, "A\t")
			scope := getScope(f)
			name := cleanName(f)
			add(fmt.Sprintf("feat%s: add %s", scope, name), "new file")
		}
	}

	// Rule 3: Delete file
	for _, line := range strings.Split(status, "\n") {
		if strings.HasPrefix(line, "D\t") {
			f := strings.TrimPrefix(line, "D\t")
			name := cleanName(f)
			add(fmt.Sprintf("refactor: remove %s", name), "delete file")
		}
	}

	// Rule 4: Only test files
	if onlyMatch(files, []string{".test.", "_test.", ".spec."}) {
		add("test: update tests", "only test files")
	}

	// Rule 5: Only docs
	if onlyExt(files, ".md", ".txt", ".rst") {
		add("docs: update documentation", "only docs")
	}

	// Rule 6: Refactor by stats
	ins, del := parseStat(stat)
	if del > ins*2 && del > 0 {
		add("refactor: simplify code", "deletions >> insertions")
	}

	// Rule 7: Style - tiny change
	if ins+del < 5 && ins+del > 0 {
		add("style: format code", "tiny change")
	}

	// Parse diff for dynamic suggestions
	addedLines, removedLines := parseDiffLines(diff)

	// Rule 8: Function detection
	if funcName := findAddedFunction(addedLines); funcName!= "" {
		scope := getScope(files[0])
		add(fmt.Sprintf("feat%s: add %s", scope, funcName), "new function")
	}

	// Rule 9: Fix pattern
	if isFix(addedLines, removedLines) {
		scope := getScope(files[0])
		add(fmt.Sprintf("fix%s: resolve issue", scope), "bug fix pattern")
	}

	// Rule 10: Debug code
	if hasDebugCode(addedLines) {
		add("chore: add debug logs", "debug code")
	}
	if hasRemovedDebug(removedLines) {
		add("chore: remove debug logs", "cleanup")
	}

	// Rule 11: Import change
	if hasImportChange(addedLines, removedLines) {
		add("refactor: update imports", "import change")
	}

	// Rule 12: Folder based
	if folderType := getFolderType(files[0]); folderType!= "" {
		name := cleanName(files[0])
		add(fmt.Sprintf("%s: update %s", folderType, name), "folder based")
	}

	// FORCE 3 SUGGESTIONS
	baseName := cleanName(files[0])
	scope := getScope(files[0])
	changeDesc := analyzeChanges(addedLines, removedLines)

	fallbacks := []Suggestion{
		{fmt.Sprintf("chore%s: update %s", scope, baseName), "fallback"},
		{fmt.Sprintf("refactor%s: %s", scope, changeDesc), "generic"},
		{fmt.Sprintf("feat%s: enhance functionality", scope), "generic"},
		{fmt.Sprintf("chore: apply code changes"), "generic"},
	}

	for _, g := range fallbacks {
		if len(suggestions) >= 3 {
			break
		}
		add(g.Msg, g.Rule)
	}

	// Print suggestions
	fmt.Println()
	green.Println("Select a commit message:")
	for i, s := range suggestions {
		fmt.Printf("%d. %s\n", i+1, s.Msg)
		hiBlack.Printf(" %s\n", s.Rule)
	}
	hiBlack.Println("q. quit")

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\nChoice [1-3/q]: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "q" {
		yellow.Println("Cancelled.")
		return
	}

	num, err := strconv.Atoi(input)
	if err!= nil || num < 1 || num > len(suggestions) {
		red.Println("Invalid input.")
		os.Exit(1)
	}

	msg := suggestions[num-1].Msg
	cmd := exec.Command("git", "commit", "-m", msg)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err!= nil {
		red.Printf("Commit failed: %v\n", err)
		os.Exit(1)
	}
	green.Printf("Committed: %s\n", msg)
}

func runGit(args...string) (string, error) {
	cmd := exec.Command("git", args...)
	out, err := cmd.Output()
	return string(out), err
}

func getScope(file string) string {
	parts := strings.Split(filepath.ToSlash(file), "/")
	if len(parts) > 1 && parts[0]!= "." && parts[0]!= "" {
		return fmt.Sprintf("(%s)", parts[0])
	}
	return ""
}

func cleanName(file string) string {
	name := filepath.Base(file)
	name = strings.TrimSuffix(name, filepath.Ext(name))
	name = strings.ReplaceAll(name, "_", " ")
	name = strings.ReplaceAll(name, "-", " ")
	return strings.ToLower(name)
}

func parseDiffLines(diff string) ([]string, []string) {
	var added, removed []string
	for _, line := range strings.Split(diff, "\n") {
		if strings.HasPrefix(line, "+") &&!strings.HasPrefix(line, "+++") {
			added = append(added, strings.TrimPrefix(line, "+"))
		}
		if strings.HasPrefix(line, "-") &&!strings.HasPrefix(line, "---") {
			removed = append(removed, strings.TrimPrefix(line, "-"))
		}
	}
	return added, removed
}

func findAddedFunction(addedLines []string) string {
	re := regexp.MustCompile(`(?i)func\s+(\w+)|function\s+(\w+)|const\s+(\w+)\s*=|def\s+(\w+)`)
	for _, line := range addedLines {
		if m := re.FindStringSubmatch(line); len(m) > 1 {
			for i := 1; i < len(m); i++ {
				if m[i]!= "" {
					return strings.ToLower(m[i])
				}
			}
		}
	}
	return ""
}

func isFix(added, removed []string) bool {
	fixKeywords := []string{"fix", "bug", "error", "issue", "null", "undefined", "nil", "panic", "exception"}
	for _, line := range append(added, removed...) {
		l := strings.ToLower(line)
		for _, kw := range fixKeywords {
			if strings.Contains(l, kw) {
				return true
			}
		}
	}
	return false
}

func hasDebugCode(added []string) bool {
	for _, line := range added {
		l := strings.ToLower(strings.TrimSpace(line))
		if strings.HasPrefix(l, "console.log") || strings.HasPrefix(l, "print(") ||
		   strings.HasPrefix(l, "fmt.print") || strings.HasPrefix(l, "debugger") {
			return true
		}
	}
	return false
}

func hasRemovedDebug(removed []string) bool {
	for _, line := range removed {
		l := strings.ToLower(strings.TrimSpace(line))
		if strings.HasPrefix(l, "console.log") || strings.HasPrefix(l, "print(") {
			return true
		}
	}
	return false
}

func hasImportChange(added, removed []string) bool {
	for _, line := range append(added, removed...) {
		l := strings.TrimSpace(line)
		if strings.HasPrefix(l, "import ") || strings.HasPrefix(l, "require(") ||
		   strings.HasPrefix(l, "from ") || strings.HasPrefix(l, "use ") {
			return true
		}
	}
	return false
}

func analyzeChanges(added, removed []string) string {
	if len(added) > len(removed)*2 {
		return "add new code"
	}
	if len(removed) > len(added)*2 {
		return "remove code"
	}
	if len(added) > 0 && len(removed) > 0 {
		return "modify implementation"
	}
	return "update code"
}

func onlyMatch(files []string, patterns []string) bool {
	for _, f := range files {
		ok := false
		lf := strings.ToLower(f)
		for _, p := range patterns {
			if strings.Contains(lf, p) {
				ok = true
				break
			}
		}
		if!ok {
			return false
		}
	}
	return len(files) > 0
}

func onlyExt(files []string, exts...string) bool {
	for _, f := range files {
		ok := false
		for _, ext := range exts {
			if strings.HasSuffix(f, ext) {
				ok = true
				break
			}
		}
		if!ok {
			return false
		}
	}
	return len(files) > 0
}

func parseStat(stat string) (int, int) {
	reIns := regexp.MustCompile(`(\d+) insertion`)
	reDel := regexp.MustCompile(`(\d+) deletion`)
	ins, del := 0, 0
	if m := reIns.FindStringSubmatch(stat); len(m) > 1 {
		ins, _ = strconv.Atoi(m[1])
	}
	if m := reDel.FindStringSubmatch(stat); len(m) > 1 {
		del, _ = strconv.Atoi(m[1])
	}
	return ins, del
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func getFolderType(file string) string {
	f := strings.ToLower(file)
	if strings.Contains(f, "api/") || strings.Contains(f, "routes/") ||
	   strings.Contains(f, "controller") || strings.Contains(f, "handler") {
		return "feat"
	}
	if strings.Contains(f, "util") || strings.Contains(f, "helper") ||
	   strings.Contains(f, "lib/") || strings.Contains(f, "pkg/") {
		return "refactor"
	}
	if strings.Contains(f, "test") || strings.Contains(f, "spec") {
		return "test"
	}
	if strings.Contains(f, "config") || strings.Contains(f, "env") {
		return "chore"
	}
	return ""
}