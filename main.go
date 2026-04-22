package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

type Suggestion struct {
	Subject string
	Body string
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

	add := func(subject, body, rule string) {
		subject = strings.ToLower(subject)
		key := subject + body
		if!seen[key] && len(suggestions) < 3 {
			suggestions = append(suggestions, Suggestion{subject, body, rule})
			seen[key] = true
		}
	}

	// Parse diff for details
	addedLines, removedLines := parseDiffLines(diff)
	changes := analyzeChangesDetailed(addedLines, removedLines, files)

	// Rule 1: Version bump
	for _, f := range files {
		if strings.HasSuffix(f, "package.json") {
			re := regexp.MustCompile(`\+\s*"version":\s*"(.*)"`)
			if m := re.FindStringSubmatch(diff); len(m) > 1 {
				body := "- Bump package version\n- Update dependencies"
				add(fmt.Sprintf("chore: bump version to %s", m[1]), body, "version bump")
			}
		}
	}

	// Rule 2: New files with details
	newFiles := getNewFiles(status)
	if len(newFiles) > 0 {
		scope := getScope(newFiles[0])
		body := buildBodyFromFiles("Added:", newFiles)
		add(fmt.Sprintf("feat%s: add new files", scope), body, "new files")
	}

	// Rule 3: Deleted files
	delFiles := getDeletedFiles(status)
	if len(delFiles) > 0 {
		body := buildBodyFromFiles("Removed:", delFiles)
		add("refactor: remove obsolete files", body, "delete files")
	}

	// Rule 4: Function additions with details
	if len(changes.NewFunctions) > 0 {
		scope := getScope(files[0])
		body := buildBodyFromItems("Added functions:", changes.NewFunctions)
		if len(changes.NewFunctions) == 1 {
			add(fmt.Sprintf("feat%s: add %s", scope, changes.NewFunctions[0]), body, "new function")
		} else {
			add(fmt.Sprintf("feat%s: add multiple functions", scope), body, "new functions")
		}
	}

	// Rule 5: Bug fixes with context
	if changes.IsFix {
		scope := getScope(files[0])
		body := buildBodyFromItems("Fixed:", changes.FixDetails)
		add(fmt.Sprintf("fix%s: resolve issues", scope), body, "bug fix")
	}

	// Rule 6: Refactor detection
	if changes.IsRefactor {
		scope := getScope(files[0])
		body := fmt.Sprintf("- Simplified logic in %s\n- Reduced complexity", cleanName(files[0]))
		add(fmt.Sprintf("refactor%s: improve code structure", scope), body, "refactor")
	}

	// Rule 7: Test files
	if onlyMatch(files, []string{".test.", "_test.", ".spec."}) {
		body := buildBodyFromFiles("Updated tests:", files)
		add("test: update test cases", body, "tests")
	}

	// Rule 8: Docs
	if onlyExt(files, ".md", ".txt", ".rst") {
		body := "- Update documentation\n- Improve clarity"
		add("docs: update documentation", body, "docs")
	}

	// Rule 9: Import/dependency changes
	if changes.ImportChanges {
		body := buildBodyFromItems("Updated imports:", changes.ImportDetails)
		add("refactor: update dependencies", body, "imports")
	}

	// FORCE 3 SUGGESTIONS with detailed bodies
	baseName := cleanName(files[0])
	scope := getScope(files[0])
	ins, del := parseStat(stat)

	genericBody1 := fmt.Sprintf("- Modified %s\n- Updated %d files\n- %d insertions(+), %d deletions(-)",
		baseName, len(files), ins, del)
	genericBody2 := fmt.Sprintf("- Improve %s implementation\n- Update related logic", baseName)
	genericBody3 := "- Apply code changes\n- Update functionality"

	fallbacks := []Suggestion{
		{fmt.Sprintf("chore%s: update %s", scope, baseName), genericBody1, "fallback"},
		{fmt.Sprintf("refactor%s: improve %s", scope, baseName), genericBody2, "generic"},
		{fmt.Sprintf("feat%s: enhance functionality", scope), genericBody3, "generic"},
	}

	for _, g := range fallbacks {
		if len(suggestions) >= 3 {
			break
		}
		add(g.Subject, g.Body, g.Rule)
	}

	// Print suggestions
	fmt.Println()
	green.Println("Select a commit message:")
	for i, s := range suggestions {
		fmt.Printf("%d. %s\n", i+1, s.Subject)
		if s.Body!= "" {
			for _, line := range strings.Split(s.Body, "\n") {
				hiBlack.Printf(" %s\n", line)
			}
		}
		hiBlack.Printf(" %s\n", s.Rule)
		fmt.Println()
	}
	hiBlack.Println("q. quit")

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Choice [1-3/q]: ")
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

	selected := suggestions[num-1]
	msg := selected.Subject
	if selected.Body!= "" {
		msg = selected.Subject + "\n\n" + selected.Body
	}

	cmd := exec.Command("git", "commit", "-m", msg)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err!= nil {
		red.Printf("Commit failed: %v\n", err)
		os.Exit(1)
	}
	green.Printf("Committed: %s\n", selected.Subject)
}

// --- STRUCTS & HELPERS ---

type ChangeAnalysis struct {
	NewFunctions []string
	IsFix bool
	FixDetails []string
	IsRefactor bool
	ImportChanges bool
	ImportDetails []string
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

func getNewFiles(status string) []string {
	var files []string
	for _, line := range strings.Split(status, "\n") {
		if strings.HasPrefix(line, "A\t") {
			files = append(files, strings.TrimPrefix(line, "A\t"))
		}
	}
	return files
}

func getDeletedFiles(status string) []string {
	var files []string
	for _, line := range strings.Split(status, "\n") {
		if strings.HasPrefix(line, "D\t") {
			files = append(files, strings.TrimPrefix(line, "D\t"))
		}
	}
	return files
}

func parseDiffLines(diff string) ([]string, []string) {
	var added, removed []string
	for _, line := range strings.Split(diff, "\n") {
		if strings.HasPrefix(line, "+") &&!strings.HasPrefix(line, "+++") {
			added = append(added, strings.TrimSpace(strings.TrimPrefix(line, "+")))
		}
		if strings.HasPrefix(line, "-") &&!strings.HasPrefix(line, "---") {
			removed = append(removed, strings.TrimSpace(strings.TrimPrefix(line, "-")))
		}
	}
	return added, removed
}

func analyzeChangesDetailed(added, removed []string, files []string) ChangeAnalysis {
	analysis := ChangeAnalysis{}
	
	// Find new functions
	re := regexp.MustCompile(`(?i)^func\s+(\w+)|^function\s+(\w+)|^const\s+(\w+)\s*=.*=>|^def\s+(\w+)`)
	for _, line := range added {
		if m := re.FindStringSubmatch(line); len(m) > 1 {
			for i := 1; i < len(m); i++ {
				if m[i]!= "" {
					analysis.NewFunctions = append(analysis.NewFunctions, strings.ToLower(m[i]))
				}
			}
		}
	}

	// Check for fixes
	fixKeywords := []string{"fix", "bug", "error", "null", "undefined", "panic", "exception", "issue"}
	for _, line := range append(added, removed...) {
		l := strings.ToLower(line)
		for _, kw := range fixKeywords {
			if strings.Contains(l, kw) && len(line) < 100 {
				analysis.IsFix = true
				analysis.FixDetails = append(analysis.FixDetails, strings.TrimSpace(line))
				break
			}
		}
	}

	// Check imports
	for _, line := range append(added, removed...) {
		l := strings.TrimSpace(line)
		if strings.HasPrefix(l, "import ") || strings.HasPrefix(l, "require(") ||
		   strings.HasPrefix(l, "from ") || strings.HasPrefix(l, "use ") {
			analysis.ImportChanges = true
			analysis.ImportDetails = append(analysis.ImportDetails, l)
		}
	}

	// Refactor detection
	if len(removed) > len(added)*2 && len(removed) > 3 {
		analysis.IsRefactor = true
	}

	return analysis
}

func buildBodyFromFiles(prefix string, files []string) string {
	if len(files) == 0 {
		return ""
	}
	lines := []string{prefix}
	for i, f := range files {
		if i >= 5 {
			lines = append(lines, fmt.Sprintf("-... and %d more", len(files)-5))
			break
		}
		lines = append(lines, "- "+filepath.Base(f))
	}
	return strings.Join(lines, "\n")
}

func buildBodyFromItems(prefix string, items []string) string {
	if len(items) == 0 {
		return ""
	}
	lines := []string{prefix}
	seen := make(map[string]bool)
	for _, item := range items {
		if!seen[item] && len(lines) < 6 {
			lines = append(lines, "- "+item)
			seen[item] = true
		}
	}
	return strings.Join(lines, "\n")
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