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
		red.Println("❌ Git not found in PATH. Install from https://git-scm.com")
		os.Exit(1)
	}

	nameOnly, _ := runGit("diff", "--staged", "--name-only")
	filesStr := strings.TrimSpace(nameOnly)
	if filesStr == "" {
		yellow.Println("⚠️ No staged files found. Run 'git add' first.")
		os.Exit(1)
	}
	files := strings.Split(filesStr, "\n")

	diff, _ := runGit("diff", "--staged")
	stat, _ := runGit("diff", "--staged", "--shortstat")
	status, _ := runGit("diff", "--staged", "--name-status")

	cyan.Println("🤖 Analyzing changes...")

	suggestions := []Suggestion{}
	seen := make(map[string]bool)

	add := func(msg, rule string) {
		if!seen[msg] && len(suggestions) < 3 {
			suggestions = append(suggestions, Suggestion{msg, rule})
			seen[msg] = true
		}
	}

	// --- RULE ENGINE START ---

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
	if strings.Contains(status, "A\t") {
		for _, line := range strings.Split(status, "\n") {
			if strings.HasPrefix(line, "A\t") {
				f := strings.TrimPrefix(line, "A\t")
				scope := getScope(f)
				name := strings.TrimSuffix(filepath.Base(f), filepath.Ext(f))
				add(fmt.Sprintf("feat%s: add %s", scope, name), "new file")
				break
			}
		}
	}

	// Rule 3: Delete file
	if strings.Contains(status, "D\t") {
		for _, line := range strings.Split(status, "\n") {
			if strings.HasPrefix(line, "D\t") {
				f := strings.TrimPrefix(line, "D\t")
				name := strings.TrimSuffix(filepath.Base(f), filepath.Ext(f))
				add(fmt.Sprintf("refactor: remove %s", name), "delete file")
				break
			}
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

	// Rule 8/9: package.json deps
	if contains(files, "package.json") {
		if strings.Contains(diff, `+ "`) &&!strings.Contains(diff, `- "`) {
			add("feat: add dependency", "package.json add")
		}
		if strings.Contains(diff, `- "`) &&!strings.Contains(diff, `+ "`) {
			add("chore: remove dependency", "package.json remove")
		}
	}

	// Rule 10:.env.example
	if contains(files, ".env.example") {
		add("chore: update env", ".env change")
	}

	// Rule 11: Folder based
	if folderType := getFolderType(files); folderType!= "" {
		name := strings.TrimSuffix(filepath.Base(files[0]), filepath.Ext(files[0]))
		add(fmt.Sprintf("%s: update %s", folderType, name), "folder based")
	}

	// --- RULE ENGINE END ---

	// FORCE 3 SUGGESTIONS: Agar 3 se kam hai to generic bharo
	baseName := strings.TrimSuffix(filepath.Base(files[0]), filepath.Ext(files[0]))
	scope := getScope(files[0])

	genericPool := []Suggestion{
		{fmt.Sprintf("chore%s: update %s", scope, baseName), "fallback"},
		{fmt.Sprintf("refactor%s: improve %s", scope, baseName), "generic"},
		{fmt.Sprintf("feat%s: enhance functionality", scope), "generic"},
		{fmt.Sprintf("fix%s: apply changes", scope), "generic"},
		{fmt.Sprintf("chore: update codebase"), "generic"},
		{fmt.Sprintf("style: format files"), "generic"},
	}

	for _, g := range genericPool {
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
		hiBlack.Printf(" └─ %s\n", s.Rule)
	}
	hiBlack.Println("q. quit")

	// Input
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
		red.Println("❌ Invalid input.")
		os.Exit(1)
	}

	msg := suggestions[num-1].Msg
	cmd := exec.Command("git", "commit", "-m", msg)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err!= nil {
		red.Printf("❌ Commit failed: %v\n", err)
		os.Exit(1)
	}
	green.Printf("✅ Committed: %s\n", msg)
}

func runGit(args...string) (string, error) {
	cmd := exec.Command("git", args...)
	out, err := cmd.Output()
	return string(out), err
}

func getScope(file string) string {
	parts := strings.Split(filepath.ToSlash(file), "/")
	if len(parts) > 1 && parts[0]!= "." {
		return fmt.Sprintf("(%s)", parts[0])
	}
	return ""
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

func getFolderType(files []string) string {
	if len(files) == 0 {
		return ""
	}
	f := strings.ToLower(files[0])
	if strings.Contains(f, "api/") || strings.Contains(f, "routes/") || strings.Contains(f, "controller") {
		return "feat"
	}
	if strings.Contains(f, "util") || strings.Contains(f, "helper") || strings.Contains(f, "lib/") {
		return "refactor"
	}
	return ""
}