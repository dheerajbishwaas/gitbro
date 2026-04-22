package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
)

type Suggestion struct {
	Subject string
	Body string
	Rule string
}

type GeminiRequest struct {
	Contents []GeminiContent `json:"contents"`
}

type GeminiContent struct {
	Parts []GeminiPart `json:"parts"`
}

type GeminiPart struct {
	Text string `json:"text"`
}

type GeminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

func main() {
	green := color.New(color.FgGreen)
	yellow := color.New(color.FgYellow)
	red := color.New(color.FgRed)
	cyan := color.New(color.FgCyan)
	hiBlack := color.New(color.FgHiBlack)

	if _, err := exec.LookPath("git"); err!= nil {
		red.Println("Git not found in PATH")
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

	projectName := getProjectName()
	branch, _ := runGit("rev-parse", "--abbrev-ref", "HEAD")
	
	suggestions := []Suggestion{}
	apiKey := os.Getenv("GEMINI_API_KEY")
	
	if apiKey!= "" {
		cyan.Println("Calling Gemini AI (free)...")
		aiSuggestions, err := getGeminiSuggestions(diff, stat, files, projectName, branch, apiKey)
		if err == nil && len(aiSuggestions) > 0 {
			suggestions = aiSuggestions
			hiBlack.Println("AI suggestions ready")
		} else {
			yellow.Println("AI failed, using rule engine")
		}
	} else {
		yellow.Println("GEMINI_API_KEY not set. Using rule engine.")
		yellow.Println("Get free key: https://aistudio.google.com/app/apikey")
	}

	if len(suggestions) == 0 {
		suggestions = getRuleBasedSuggestions(files, diff, stat, status)
	}

	fmt.Println()
	green.Println("Select a commit message:")
	for i, s := range suggestions {
		fmt.Printf("%d. %s\n", i+1, s.Subject)
		if s.Body!= "" {
			for _, line := range strings.Split(s.Body, "\n") {
				if strings.TrimSpace(line)!= "" {
					hiBlack.Printf(" %s\n", line)
				}
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

func getGeminiSuggestions(diff, stat string, files []string, project, branch, apiKey string) ([]Suggestion, error) {
	if len(diff) > 8000 {
		diff = diff[:8000] + "\n... [truncated]"
	}

	prompt := fmt.Sprintf(`You are a git commit message generator. Based on the git diff, generate 3 conventional commit messages.

Project: %s
Branch: %s
Files: %s
Stats: %s

Diff:
%s

Rules:
1. Use conventional commits: feat, fix, refactor, chore, docs, test, style
2. Subject line max 72 chars, lowercase, no period
3. Add scope like (auth) if files are in subdirectories
4. Body: 2-4 bullet points explaining WHAT and WHY, max 100 chars per line
5. Be specific - mention function names from diff
6. Output ONLY valid JSON array, no markdown

Format: [{"subject":"feat(auth): add login","body":"- Add handleLogin function\n- Validate credentials"},{"subject":"...","body":"..."},{"subject":"...","body":"..."}]`,
		project, branch, strings.Join(files, ", "), strings.TrimSpace(stat), diff)

	reqBody := GeminiRequest{
		Contents: []GeminiContent{{
			Parts: []GeminiPart{{Text: prompt}},
		}},
	}

	jsonData, _ := json.Marshal(reqBody)
	url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash-latest:generateContent?key=" + apiKey
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err!= nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode!= 200 {
		return nil, fmt.Errorf("gemini API error: %s", string(body))
	}

	var geminiResp GeminiResponse
	if err := json.Unmarshal(body, &geminiResp); err!= nil {
		return nil, err
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("empty response")
	}

	cleanText := strings.TrimSpace(geminiResp.Candidates[0].Content.Parts[0].Text)
	cleanText = strings.TrimPrefix(cleanText, "```json")
	cleanText = strings.TrimPrefix(cleanText, "```")
	cleanText = strings.TrimSuffix(cleanText, "```")
	cleanText = strings.TrimSpace(cleanText)

	var rawSugs []map[string]string
	if err := json.Unmarshal([]byte(cleanText), &rawSugs); err!= nil {
		return nil, err
	}

	var suggestions []Suggestion
	for _, rs := range rawSugs {
		if len(suggestions) >= 3 {
			break
		}
		suggestions = append(suggestions, Suggestion{
			Subject: rs["subject"],
			Body: rs["body"],
			Rule: "gemini ai",
		})
	}

	return suggestions, nil
}

func getRuleBasedSuggestions(files []string, diff, stat, status string) []Suggestion {
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

	addedLines, removedLines := parseDiffLines(diff)
	changes := analyzeChangesDetailed(addedLines, removedLines)

	if len(changes.NewFunctions) > 0 {
		scope := getScope(files[0])
		body := buildBodyFromItems("Added functions:", changes.NewFunctions)
		add(fmt.Sprintf("feat%s: add %s", scope, changes.NewFunctions[0]), body, "new function")
	}

	if changes.IsFix {
		scope := getScope(files[0])
		body := buildBodyFromItems("Fixed:", changes.FixDetails)
		add(fmt.Sprintf("fix%s: resolve issues", scope), body, "bug fix")
	}

	baseName := cleanName(files[0])
	scope := getScope(files[0])
	ins, del := parseStat(stat)
	genericBody := fmt.Sprintf("- Modified %s\n- %d files changed\n- %d insertions(+), %d deletions(-)",
		baseName, len(files), ins, del)

	fallbacks := []Suggestion{
		{fmt.Sprintf("chore%s: update %s", scope, baseName), genericBody, "fallback"},
		{fmt.Sprintf("refactor%s: improve %s", scope, baseName), "- Improve code structure\n- Update implementation", "generic"},
		{fmt.Sprintf("feat%s: enhance functionality", scope), "- Apply code changes\n- Update features", "generic"},
	}

	for _, g := range fallbacks {
		if len(suggestions) >= 3 {
			break
		}
		add(g.Subject, g.Body, g.Rule)
	}

	return suggestions
}

func getProjectName() string {
	dir, _ := os.Getwd()
	return filepath.Base(dir)
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

type ChangeAnalysis struct {
	NewFunctions []string
	IsFix bool
	FixDetails []string
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

func analyzeChangesDetailed(added, removed []string) ChangeAnalysis {
	analysis := ChangeAnalysis{}
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

	fixKeywords := []string{"fix", "bug", "error", "null", "undefined", "panic", "exception"}
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
	return analysis
}

func buildBodyFromItems(prefix string, items []string) string {
	if len(items) == 0 {
		return ""
	}
	lines := []string{prefix}
	seen := make(map[string]bool)
	for _, item := range items {
		if!seen[item] && len(lines) < 5 && item!= "" {
			lines = append(lines, "- "+item)
			seen[item] = true
		}
	}
	return strings.Join(lines, "\n")
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