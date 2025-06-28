package main

import (
    "context"
    "encoding/json"
    "fmt"
    "html"
    "io/ioutil"
    "os"
    "path/filepath"
    "regexp"
    "strconv"
    "strings"
    "time"

    "github.com/google/go-github/v57/github"
    "github.com/joho/godotenv"
    "golang.org/x/oauth2"
)

type DailySolution struct {
    Date       string `json:"date"`
    Target     int    `json:"target"`
    CSS        string `json:"css"`
    Bytes      int    `json:"bytes"`
    Difficulty string `json:"difficulty"`
    Screenshot string `json:"screenshot"`
    HasImage   bool   `json:"has_image"`
}

type MonthlyBattles struct {
    Month          string          `json:"month"`
    Year           int             `json:"year"`
    MonthName      string          `json:"month_name"`
    Solutions      []DailySolution `json:"solutions"`
    TotalSolutions int             `json:"total_solutions"`
}

func main() {
    // Load .env file FIRST (before anything else)
    if err := godotenv.Load(); err != nil {
        // Ignore error in production - .env file won't exist in GitHub Actions
    }

    // Get GitHub token from environment
    token := os.Getenv("GITHUB_TOKEN")
    if token == "" {
        fmt.Println("GITHUB_TOKEN environment variable is required")
        os.Exit(1)
    }

    // Initialize GitHub client
    ctx := context.Background()
    ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
    tc := oauth2.NewClient(ctx, ts)
    client := github.NewClient(tc)

    // Sync solutions
    if err := syncSolutions(ctx, client); err != nil {
        fmt.Printf("Error syncing solutions: %v\n", err)
        os.Exit(1)
    }

    fmt.Println("Successfully synced CSS Battle solutions!")
}

func syncSolutions(ctx context.Context, client *github.Client) error {
    owner := os.Getenv("GITHUB_OWNER")
    repo := os.Getenv("GITHUB_REPO")
    
    if owner == "" || repo == "" {
        return fmt.Errorf("GITHUB_OWNER and GITHUB_REPO environment variables are required")
    }
    
    // 1. Get daily-targets directory contents
    _, directoryContent, _, err := client.Repositories.GetContents(ctx, owner, repo, "daily-targets", nil)
    if err != nil {
        return fmt.Errorf("failed to get daily-targets directory: %v", err)
    }

    var allSolutions []DailySolution
    monthlyGroups := make(map[string][]DailySolution)

    fmt.Printf("Found %d items in daily-targets directory\n", len(directoryContent))

    // 2. Process each daily target directory
    for _, item := range directoryContent {
        if item.GetType() == "dir" {
            date := item.GetName() // Should be YYYY-MM-DD format
            if !isValidDateFormat(date) {
                fmt.Printf("Skipping invalid date format: %s\n", date)
                continue
            }

            fmt.Printf("Processing date: %s\n", date)
            solution, err := processDailyTarget(ctx, client, owner, repo, date)
            if err != nil {
                fmt.Printf("Warning: failed to process %s: %v\n", date, err)
                continue
            }

            if solution != nil {
                allSolutions = append(allSolutions, *solution)
                month := date[:7] // Extract YYYY-MM
                monthlyGroups[month] = append(monthlyGroups[month], *solution)
                fmt.Printf("✅ Processed %s: %d bytes\n", date, solution.Bytes)
            }
        }
    }

    fmt.Printf("Total solutions processed: %d\n", len(allSolutions))

    // 3. Generate monthly content files
    if err := generateMonthlyContent(monthlyGroups); err != nil {
        return fmt.Errorf("failed to generate monthly content: %v", err)
    }

    // 4. Generate API files
    if err := generateAPIFiles(allSolutions); err != nil {
        return fmt.Errorf("failed to generate API files: %v", err)
    }

    return nil
}

func processDailyTarget(ctx context.Context, client *github.Client, owner, repo, date string) (*DailySolution, error) {
    targetPath := fmt.Sprintf("daily-targets/%s", date)
    
    // Get directory contents
    _, contents, _, err := client.Repositories.GetContents(ctx, owner, repo, targetPath, nil)
    if err != nil {
        return nil, err
    }

    var htmlContent string
    var screenshotFiles []string
    var linkedinContent string

    // Find solutions.html, screenshots, and linkedin post
    for _, file := range contents {
        fileName := file.GetName()
        filePath := file.GetPath()
        
        if fileName == "solutions.html" {
            // Download HTML content
            fileContent, _, _, err := client.Repositories.GetContents(ctx, owner, repo, filePath, nil)
            if err != nil {
                return nil, fmt.Errorf("failed to get HTML content: %v", err)
            }
            content, err := fileContent.GetContent()
            if err != nil {
                return nil, fmt.Errorf("failed to decode HTML content: %v", err)
            }
            htmlContent = content
            
        } else if fileName == "linkedin-post.md" {
            // Get LinkedIn post for additional context (optional)
            fileContent, _, _, err := client.Repositories.GetContents(ctx, owner, repo, filePath, nil)
            if err == nil {
                if content, err := fileContent.GetContent(); err == nil {
                    linkedinContent = content
                }
            }
            
        } else if strings.HasSuffix(fileName, "-comparison.png") || 
                 strings.HasSuffix(fileName, "-comparison.jpg") ||
                 strings.HasSuffix(fileName, "-comparison.jpeg") {
            screenshotFiles = append(screenshotFiles, fileName)
            
            // Download screenshot
            if err := downloadScreenshot(ctx, client, owner, repo, filePath, date, fileName); err != nil {
                fmt.Printf("Warning: failed to download screenshot %s: %v\n", fileName, err)
            }
        }
    }

    if htmlContent == "" {
        return nil, fmt.Errorf("solutions.html not found")
    }

    // Parse the HTML to extract solution data
    solution := parseSolutionHTML(htmlContent, date)
    
    // Set screenshot info
    if len(screenshotFiles) > 0 {
        solution.Screenshot = screenshotFiles[0] // Use first screenshot
        solution.HasImage = true
    }
    
    // Try to extract target number from LinkedIn post if available
    if linkedinContent != "" && solution.Target == 0 {
        if target := extractTargetFromLinkedIn(linkedinContent); target > 0 {
            solution.Target = target
        }
    }

    return solution, nil
}

func parseSolutionHTML(htmlContent, date string) *DailySolution {
    // Debug: Print what we're parsing
    fmt.Printf("Parsing HTML for %s (length: %d chars)\n", date, len(htmlContent))
    fmt.Printf("First 200 chars: %s\n", htmlContent[:min(200, len(htmlContent))])
    
    // Extract target number from filename or infer from date
    target := extractTargetFromDate(date)
    
    // Remove HTML comments first
    contentRegex := regexp.MustCompile(`<!--[^>]*-->`)
    cleanContent := contentRegex.ReplaceAllString(htmlContent, "")
    cleanContent = strings.TrimSpace(cleanContent)
    
    var css string
    
    // Look for <style> tag anywhere in the content
    styleIndex := strings.Index(cleanContent, "<style>")
    if styleIndex != -1 {
        // Extract everything after <style> 
        css = cleanContent[styleIndex+7:] // 7 = len("<style>")
        
        // Look for closing </style> tag
        endIndex := strings.Index(css, "</style>")
        if endIndex != -1 {
            css = css[:endIndex]
        }
        
        // Clean up any remaining HTML tags and whitespace
        css = strings.TrimSpace(css)
        
        fmt.Printf("Found <style> tag, extracted CSS (length: %d): %s...\n", 
            len(css), css[:min(100, len(css))])
    } else {
        // Fallback: try to find CSS-like content
        // Look for patterns that start with CSS selectors
        cssPatterns := []string{
            "*{", "p{", "body{", "&{", "div{", "[", ".",
        }
        
        for _, pattern := range cssPatterns {
            if strings.Contains(cleanContent, pattern) {
                // Find the start of CSS-like content
                startIdx := strings.Index(cleanContent, pattern)
                if startIdx != -1 {
                    css = cleanContent[startIdx:]
                    css = strings.TrimSpace(css)
                    break
                }
            }
        }
        
        fmt.Printf("No <style> tag found, extracted CSS-like content (length: %d): %s...\n", 
            len(css), css[:min(100, len(css))])
    }
    
    // Final cleanup - remove any remaining HTML artifacts
    css = strings.ReplaceAll(css, "<p>", "")
    css = strings.ReplaceAll(css, "</p>", "")
    css = strings.ReplaceAll(css, "<a>", "")
    css = strings.ReplaceAll(css, "</a>", "")
    css = strings.ReplaceAll(css, "<div>", "")
    css = strings.ReplaceAll(css, "</div>", "")
    css = strings.TrimSpace(css)

    // Calculate bytes and difficulty
    bytes := len(css)
    difficulty := classifyDifficulty(bytes)
    
    fmt.Printf("Final result for %s: %d bytes, difficulty: %s\n", date, bytes, difficulty)
    if bytes == 0 {
        fmt.Printf("WARNING: No CSS extracted for %s. Raw content: %s\n", date, cleanContent[:min(200, len(cleanContent))])
    }

    return &DailySolution{
        Date:       date,
        Target:     target,
        CSS:        css,
        Bytes:      bytes,
        Difficulty: difficulty,
    }
}

// Helper function for min
func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}

// Helper to extract target from LinkedIn post
func extractTargetFromLinkedIn(linkedinContent string) int {
    targetRegex := regexp.MustCompile(`(?i)target\s*#?(\d+)`)
    matches := targetRegex.FindStringSubmatch(linkedinContent)
    if len(matches) > 1 {
        target, _ := strconv.Atoi(matches[1])
        return target
    }
    return 0
}

// Helper function to extract target number from date or other logic
func extractTargetFromDate(date string) int {
    // Option 1: Use day of year as target number
    t, err := time.Parse("2006-01-02", date)
    if err != nil {
        return 0
    }
    return t.YearDay()
}

func downloadScreenshot(ctx context.Context, client *github.Client, owner, repo, filePath, date, fileName string) error {
    // Get file content
    fileContent, _, _, err := client.Repositories.GetContents(ctx, owner, repo, filePath, nil)
    if err != nil {
        return err
    }

    // Decode base64 content
    content, err := fileContent.GetContent()
    if err != nil {
        return err
    }

    // Ensure directory exists
    destDir := fmt.Sprintf("../public/images/cssbattle/%s", date)
    if err := os.MkdirAll(destDir, 0755); err != nil {
        return err
    }

    // Write file
    destPath := filepath.Join(destDir, fileName)
    return ioutil.WriteFile(destPath, []byte(content), 0644)
}

func generateMonthlyContent(monthlyGroups map[string][]DailySolution) error {
    contentDir := "../src/content/daily-cssbattles"
    if err := os.MkdirAll(contentDir, 0755); err != nil {
        return err
    }

    fmt.Printf("Generating monthly content files for %d months\n", len(monthlyGroups))

    for month, solutions := range monthlyGroups {
        monthlyBattles := MonthlyBattles{
            Month:          month,
            Year:           parseYear(month),
            MonthName:      parseMonthName(month),
            Solutions:      solutions,
            TotalSolutions: len(solutions),
        }

        // Generate markdown content
        content := generateMarkdownContent(monthlyBattles)
        
        fileName := fmt.Sprintf("%s.md", month)
        filePath := filepath.Join(contentDir, fileName)
        
        if err := ioutil.WriteFile(filePath, []byte(content), 0644); err != nil {
            return err
        }
        
        fmt.Printf("✅ Generated %s with %d solutions\n", fileName, len(solutions))
    }

    return nil
}

func generateAPIFiles(solutions []DailySolution) error {
    dataDir := "../src/data"
    if err := os.MkdirAll(dataDir, 0755); err != nil {
        return err
    }

    fmt.Println("Generating data files...")

    // Generate solutions.json
    solutionsJSON, err := json.MarshalIndent(solutions, "", "  ")
    if err != nil {
        return err
    }
    if err := ioutil.WriteFile(filepath.Join(dataDir, "solutions.json"), solutionsJSON, 0644); err != nil {
        return err
    }
    fmt.Println("✅ Generated solutions.json")

    // Generate latest.json
    if len(solutions) > 0 {
        latest := solutions[len(solutions)-1]
        latestJSON, err := json.MarshalIndent(latest, "", "  ")
        if err != nil {
            return err
        }
        if err := ioutil.WriteFile(filepath.Join(dataDir, "latest.json"), latestJSON, 0644); err != nil {
            return err
        }
        fmt.Println("✅ Generated latest.json")
    }

    // Generate stats.json
    stats := map[string]interface{}{
        "total":    len(solutions),
        "avgBytes": calculateAverageBytes(solutions),
    }
    if len(solutions) > 0 {
        stats["latestDate"] = solutions[len(solutions)-1].Date
    }
    
    statsJSON, err := json.MarshalIndent(stats, "", "  ")
    if err != nil {
        return err
    }
    if err := ioutil.WriteFile(filepath.Join(dataDir, "stats.json"), statsJSON, 0644); err != nil {
        return err
    }
    fmt.Println("✅ Generated stats.json")

    return nil
}

// Helper functions
func isValidDateFormat(date string) bool {
    matched, _ := regexp.MatchString(`^\d{4}-\d{2}-\d{2}$`, date)
    return matched
}

func classifyDifficulty(bytes int) string {
    switch {
    case bytes >= 200:
        return "easy"
    case bytes >= 100:
        return "medium"
    default:
        return "hard"
    }
}

func parseYear(month string) int {
    year, _ := strconv.Atoi(month[:4])
    return year
}

func parseMonthName(month string) string {
    monthNames := []string{"", "January", "February", "March", "April", "May", "June",
        "July", "August", "September", "October", "November", "December"}
    monthNum, _ := strconv.Atoi(month[5:7])
    if monthNum >= 1 && monthNum <= 12 {
        return monthNames[monthNum]
    }
    return "Unknown"
}

func calculateAverageBytes(solutions []DailySolution) int {
    if len(solutions) == 0 {
        return 0
    }
    total := 0
    for _, solution := range solutions {
        total += solution.Bytes
    }
    return total / len(solutions)
}

func generateMarkdownContent(battles MonthlyBattles) string {
    // Escape CSS content to prevent rendering issues
    solutionsJSON := make([]map[string]interface{}, len(battles.Solutions))
    for i, solution := range battles.Solutions {
        solutionsJSON[i] = map[string]interface{}{
            "date":       solution.Date,
            "target":     solution.Target,
            "css":        html.EscapeString(solution.CSS), // Escape HTML/CSS
            "bytes":      solution.Bytes,
            "difficulty": solution.Difficulty,
            "screenshot": solution.Screenshot,
            "has_image":  solution.HasImage,
        }
    }
    
    solutionsJSONBytes, _ := json.Marshal(solutionsJSON)
    
    frontmatter := fmt.Sprintf(`---
month: "%s"
year: %d
monthName: "%s"
totalSolutions: %d
solutions: %s
---

# CSS Battle Solutions - %s %d

This month I completed %d CSS Battle challenges, practicing various CSS techniques and code golf optimization.

`, battles.Month, battles.Year, battles.MonthName, battles.TotalSolutions, 
        string(solutionsJSONBytes), battles.MonthName, battles.Year, battles.TotalSolutions)

    return frontmatter
}

func toJSONString(v interface{}) string {
    b, _ := json.Marshal(v)
    return string(b)
}
