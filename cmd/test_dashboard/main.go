package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

// TestEvent represents a line of JSON output from go test -json
type TestEvent struct {
	Time    time.Time `json:"Time"`
	Action  string    `json:"Action"`
	Package string    `json:"Package"`
	Test    string    `json:"Test"`
	Elapsed float64   `json:"Elapsed"`
	Output  string    `json:"Output"`
}

// TestResult represents the aggregated results for a specific test case
type TestResult struct {
	Name    string   `json:"name"`
	Package string   `json:"package"`
	Status  string   `json:"status"` // "pass", "fail", "skip"
	Elapsed float64  `json:"elapsed"`
	Logs    []string `json:"logs"`
}

// TestSummary is the payload returned to the UI
type TestSummary struct {
	Total     int          `json:"total"`
	Passed    int          `json:"passed"`
	Failed    int          `json:"failed"`
	Duration  float64      `json:"duration"`
	Results   []TestResult `json:"results"`
	Timestamp string       `json:"timestamp"`
}

func main() {
	port := "8085"
	log.Printf("Starting Test Dashboard Server on http://localhost:%s", port)

	http.HandleFunc("/", handleHome)
	http.HandleFunc("/api/run", handleRunTests)

	// Automatically attempt to open the browser
	go func() {
		time.Sleep(1 * time.Second)
		openBrowser("http://localhost:" + port)
	}()

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func openBrowser(url string) {
	var cmd *exec.Cmd
	// Simple multiplatform browser opening
	if strings.Contains(os.Getenv("OS"), "Windows") || os.PathSeparator == '\\' {
		cmd = exec.Command("cmd", "/c", "start", url)
	} else {
		cmd = exec.Command("open", url)
	}
	_ = cmd.Start()
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("index").Parse(htmlTemplate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = tmpl.Execute(w, nil)
}

func handleRunTests(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Run go test -json
	cmd := exec.Command("go", "test", "./test/...", "-v", "-json")
	var outBuf bytes.Buffer
	cmd.Stdout = &outBuf
	_ = cmd.Run() // Ignore exit code since failed tests return exit code 1

	scanner := bufio.NewScanner(&outBuf)
	testMap := make(map[string]*TestResult)
	var orderedNames []string

	var totalDuration float64
	var totalPassed, totalFailed int

	for scanner.Scan() {
		var event TestEvent
		if err := json.Unmarshal(scanner.Bytes(), &event); err != nil {
			continue
		}

		// We only care about specific tests, not package-level summaries
		if event.Test == "" {
			continue
		}

		res, exists := testMap[event.Test]
		if !exists {
			res = &TestResult{
				Name:    event.Test,
				Package: event.Package,
				Status:  "running",
				Logs:    []string{},
			}
			testMap[event.Test] = res
			orderedNames = append(orderedNames, event.Test)
		}

		switch event.Action {
		case "output":
			cleanOut := strings.TrimSpace(event.Output)
			if cleanOut != "" {
				res.Logs = append(res.Logs, cleanOut)
			}
		case "pass":
			res.Status = "pass"
			res.Elapsed = event.Elapsed
			totalDuration += event.Elapsed
			totalPassed++
		case "fail":
			res.Status = "fail"
			res.Elapsed = event.Elapsed
			totalDuration += event.Elapsed
			totalFailed++
		case "skip":
			res.Status = "skip"
		}
	}

	var results []TestResult
	for _, name := range orderedNames {
		results = append(results, *testMap[name])
	}

	summary := TestSummary{
		Total:     totalPassed + totalFailed,
		Passed:    totalPassed,
		Failed:    totalFailed,
		Duration:  totalDuration,
		Results:   results,
		Timestamp: time.Now().Format("2006-01-02 15:04:05"),
	}

	json.NewEncoder(w).Encode(summary)
}

const htmlTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Go Test Dashboard</title>
    <!-- Google Fonts -->
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    <link href="https://fonts.googleapis.com/css2?family=Outfit:wght@300;400;500;600;700;800&family=JetBrains+Mono:wght@400;500&display=swap" rel="stylesheet">
    
    <style>
        :root {
            --bg-primary: #0a0e17;
            --bg-secondary: rgba(18, 26, 44, 0.6);
            --border-color: rgba(255, 255, 255, 0.08);
            --text-main: #f3f4f6;
            --text-muted: #9ca3af;
            
            --success-glow: rgba(16, 185, 129, 0.15);
            --success-color: #10b981;
            --fail-glow: rgba(239, 68, 68, 0.15);
            --fail-color: #ef4444;
            --accent-glow: rgba(59, 130, 246, 0.15);
            --accent-color: #3b82f6;
        }

        * {
            box-sizing: border-box;
            margin: 0;
            padding: 0;
        }

        body {
            font-family: 'Outfit', sans-serif;
            background-color: var(--bg-primary);
            color: var(--text-main);
            min-height: 100vh;
            display: flex;
            flex-direction: column;
            overflow-x: hidden;
            background-image: 
                radial-gradient(circle at 10% 20%, rgba(59, 130, 246, 0.08) 0%, transparent 40%),
                radial-gradient(circle at 90% 80%, rgba(16, 185, 129, 0.05) 0%, transparent 40%);
        }

        header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            padding: 1.5rem 2.5rem;
            border-bottom: 1px solid var(--border-color);
            background: rgba(10, 14, 23, 0.8);
            backdrop-filter: blur(12px);
            position: sticky;
            top: 0;
            z-index: 10;
        }

        .logo-section h1 {
            font-size: 1.5rem;
            font-weight: 800;
            background: linear-gradient(135deg, #3b82f6 0%, #10b981 100%);
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
            letter-spacing: -0.5px;
        }

        .logo-section p {
            font-size: 0.85rem;
            color: var(--text-muted);
            margin-top: 2px;
        }

        .actions {
            display: flex;
            gap: 12px;
        }

        button {
            font-family: inherit;
            font-size: 0.9rem;
            font-weight: 600;
            padding: 0.65rem 1.25rem;
            border-radius: 8px;
            cursor: pointer;
            transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
            display: inline-flex;
            align-items: center;
            gap: 8px;
        }

        .btn-primary {
            background: linear-gradient(135deg, #3b82f6 0%, #2563eb 100%);
            color: white;
            border: none;
            box-shadow: 0 4px 14px rgba(59, 130, 246, 0.3);
        }

        .btn-primary:hover {
            transform: translateY(-2px);
            box-shadow: 0 6px 20px rgba(59, 130, 246, 0.4);
        }

        .btn-primary:active {
            transform: translateY(0);
        }

        .container {
            max-width: 1200px;
            width: 100%;
            margin: 0 auto;
            padding: 2.5rem;
            flex-grow: 1;
        }

        /* Stats Cards */
        .stats-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
            gap: 1.5rem;
            margin-bottom: 2.5rem;
        }

        .stat-card {
            background: var(--bg-secondary);
            border: 1px solid var(--border-color);
            border-radius: 12px;
            padding: 1.5rem;
            position: relative;
            overflow: hidden;
            backdrop-filter: blur(8px);
            transition: transform 0.3s ease;
        }

        .stat-card::before {
            content: '';
            position: absolute;
            top: 0;
            left: 0;
            width: 4px;
            height: 100%;
        }

        .stat-card.total::before { background-color: var(--accent-color); }
        .stat-card.passed::before { background-color: var(--success-color); }
        .stat-card.failed::before { background-color: var(--fail-color); }
        .stat-card.duration::before { background-color: #f59e0b; }

        .stat-label {
            font-size: 0.85rem;
            color: var(--text-muted);
            text-transform: uppercase;
            letter-spacing: 0.5px;
            font-weight: 500;
        }

        .stat-val {
            font-size: 2rem;
            font-weight: 700;
            margin-top: 0.5rem;
        }

        /* Filter bar */
        .filter-bar {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 1.5rem;
            gap: 1rem;
        }

        .filters {
            display: flex;
            background: var(--bg-secondary);
            border: 1px solid var(--border-color);
            border-radius: 8px;
            padding: 4px;
        }

        .filter-btn {
            background: transparent;
            border: none;
            color: var(--text-muted);
            padding: 0.5rem 1rem;
            font-size: 0.85rem;
            font-weight: 600;
            border-radius: 6px;
        }

        .filter-btn.active {
            background: rgba(255, 255, 255, 0.08);
            color: var(--text-main);
        }

        .search-box {
            background: var(--bg-secondary);
            border: 1px solid var(--border-color);
            border-radius: 8px;
            padding: 0.5rem 1rem;
            color: var(--text-main);
            font-family: inherit;
            font-size: 0.9rem;
            outline: none;
            min-width: 250px;
            transition: border-color 0.3s;
        }

        .search-box:focus {
            border-color: var(--accent-color);
        }

        /* Test List */
        .test-list {
            display: flex;
            flex-direction: column;
            gap: 0.75rem;
        }

        .test-item {
            background: var(--bg-secondary);
            border: 1px solid var(--border-color);
            border-radius: 10px;
            overflow: hidden;
            transition: all 0.2s ease;
        }

        .test-item:hover {
            border-color: rgba(255, 255, 255, 0.15);
            background: rgba(18, 26, 44, 0.8);
        }

        .test-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            padding: 1.25rem;
            cursor: pointer;
        }

        .test-title-group {
            display: flex;
            align-items: center;
            gap: 12px;
        }

        .status-badge {
            font-size: 0.75rem;
            font-weight: 700;
            text-transform: uppercase;
            padding: 0.25rem 0.6rem;
            border-radius: 50px;
            letter-spacing: 0.5px;
        }

        .status-badge.pass {
            background: var(--success-glow);
            color: var(--success-color);
            border: 1px solid rgba(16, 185, 129, 0.3);
        }

        .status-badge.fail {
            background: var(--fail-glow);
            color: var(--fail-color);
            border: 1px solid rgba(239, 68, 68, 0.3);
        }

        .test-name {
            font-weight: 600;
            font-size: 1rem;
        }

        .test-package {
            font-size: 0.8rem;
            color: var(--text-muted);
            font-family: 'JetBrains Mono', monospace;
        }

        .test-meta {
            display: flex;
            align-items: center;
            gap: 1rem;
        }

        .test-duration {
            font-size: 0.85rem;
            color: var(--text-muted);
            font-family: 'JetBrains Mono', monospace;
        }

        .chevron {
            width: 16px;
            height: 16px;
            fill: none;
            stroke: currentColor;
            stroke-width: 2;
            transition: transform 0.3s ease;
        }

        .test-item.expanded .chevron {
            transform: rotate(180deg);
        }

        .test-details {
            max-height: 0;
            overflow: hidden;
            transition: max-height 0.3s ease-out;
            background: rgba(5, 7, 12, 0.4);
            border-top: 0px solid var(--border-color);
        }

        .test-item.expanded .test-details {
            max-height: 500px;
            border-top: 1px solid var(--border-color);
            overflow-y: auto;
        }

        .log-output {
            padding: 1.25rem;
            font-family: 'JetBrains Mono', monospace;
            font-size: 0.85rem;
            line-height: 1.6;
            color: #d1d5db;
            white-space: pre-wrap;
            background: #05070c;
        }

        .log-line.error-highlight {
            color: var(--fail-color);
            font-weight: 600;
        }

        /* Loading / empty state */
        .empty-state {
            display: flex;
            flex-direction: column;
            align-items: center;
            justify-content: center;
            padding: 5rem 2.5rem;
            text-align: center;
        }

        .empty-state h3 {
            margin-bottom: 0.5rem;
            font-weight: 600;
        }

        .empty-state p {
            color: var(--text-muted);
            font-size: 0.9rem;
            margin-bottom: 1.5rem;
        }

        /* Spinner animation */
        .spinner {
            width: 18px;
            height: 18px;
            border: 2px solid rgba(255, 255, 255, 0.2);
            border-top-color: white;
            border-radius: 50%;
            animation: spin 0.8s linear infinite;
            display: none;
        }

        @keyframes spin {
            to { transform: rotate(360deg); }
        }

        .running .spinner {
            display: inline-block;
        }
        .running .btn-text {
            display: none;
        }
    </style>
</head>
<body>
    <header>
        <div class="logo-section">
            <h1>Backend Suite Dashboard</h1>
            <p>Interactive Go Test Suite Runner</p>
        </div>
        <div class="actions">
            <button class="btn-primary" id="runBtn" onclick="runTests()">
                <div class="spinner"></div>
                <span class="btn-text">Run Suite</span>
            </button>
        </div>
    </header>

    <div class="container">
        <div class="stats-grid">
            <div class="stat-card total">
                <div class="stat-label">Total Tests</div>
                <div class="stat-val" id="statTotal">-</div>
            </div>
            <div class="stat-card passed">
                <div class="stat-label">Passed</div>
                <div class="stat-val" id="statPassed" style="color: var(--success-color)">-</div>
            </div>
            <div class="stat-card failed">
                <div class="stat-label">Failed</div>
                <div class="stat-val" id="statFailed" style="color: var(--fail-color)">-</div>
            </div>
            <div class="stat-card duration">
                <div class="stat-label">Total Time</div>
                <div class="stat-val" id="statDuration">-</div>
            </div>
        </div>

        <div class="filter-bar">
            <div class="filters">
                <button class="filter-btn active" onclick="setFilter('all', this)">All</button>
                <button class="filter-btn" onclick="setFilter('pass', this)">Passed</button>
                <button class="filter-btn" onclick="setFilter('fail', this)">Failed</button>
            </div>
            <input type="text" class="search-box" id="searchBox" placeholder="Search tests..." oninput="filterList()">
        </div>

        <div class="test-list" id="testList">
            <div class="empty-state">
                <h3>No Results Yet</h3>
                <p>Click "Run Suite" above to run the test cases.</p>
            </div>
        </div>
    </div>

    <script>
        let currentData = null;
        let activeFilter = 'all';

        async function runTests() {
            const btn = document.getElementById('runBtn');
            btn.classList.add('running');
            btn.disabled = true;

            try {
                const response = await fetch('/api/run');
                currentData = await response.json();
                renderResults(currentData);
            } catch (error) {
                console.error("Error running tests:", error);
                alert("Failed to run tests.");
            } finally {
                btn.classList.remove('running');
                btn.disabled = false;
            }
        }

        function renderResults(data) {
            document.getElementById('statTotal').innerText = data.total;
            document.getElementById('statPassed').innerText = data.passed;
            document.getElementById('statFailed').innerText = data.failed;
            document.getElementById('statDuration').innerText = data.duration.toFixed(3) + 's';

            const listContainer = document.getElementById('testList');
            listContainer.innerHTML = '';

            if (!data.results || data.results.length === 0) {
                listContainer.innerHTML = '<div class="empty-state"><h3>No test results returned</h3></div>';
                return;
            }

            data.results.forEach((test, index) => {
                const item = document.createElement('div');
                item.className = 'test-item';
                item.setAttribute('data-status', test.status);
                item.setAttribute('data-name', test.name.toLowerCase());

                const logLines = test.logs.map(line => {
                    const isErr = line.includes("FAIL:") || line.includes("panic:") || line.includes("ERROR");
                    return '<div class="log-line ' + (isErr ? 'error-highlight' : '') + '">' + escapeHtml(line) + '</div>';
                }).join('\n');

                item.innerHTML = 
                    '<div class="test-header" onclick="toggleExpand(' + index + ')">' +
                        '<div class="test-title-group">' +
                            '<span class="status-badge ' + test.status + '">' + test.status + '</span>' +
                            '<div>' +
                                '<div class="test-name">' + test.name + '</div>' +
                                '<div class="test-package">' + test.package + '</div>' +
                            '</div>' +
                        '</div>' +
                        '<div class="test-meta">' +
                            '<span class="test-duration">' + test.elapsed.toFixed(4) + 's</span>' +
                            '<svg class="chevron" viewBox="0 0 24 24"><path d="M6 9l6 6 6-6"/></svg>' +
                        '</div>' +
                    '</div>' +
                    '<div class="test-details" id="details-' + index + '">' +
                        '<pre class="log-output">' + (logLines || "No logs captured.") + '</pre>' +
                    '</div>';

                listContainer.appendChild(item);
            });

            filterList();
        }

        function toggleExpand(index) {
            const item = document.querySelectorAll('.test-item')[index];
            item.classList.toggle('expanded');
        }

        function setFilter(status, btnElement) {
            activeFilter = status;
            document.querySelectorAll('.filter-btn').forEach(btn => btn.classList.remove('active'));
            btnElement.classList.add('active');
            filterList();
        }

        function filterList() {
            const query = document.getElementById('searchBox').value.toLowerCase();
            const items = document.querySelectorAll('.test-item');

            items.forEach(item => {
                const status = item.getAttribute('data-status');
                const name = item.getAttribute('data-name');

                const matchesStatus = (activeFilter === 'all') || (status === activeFilter);
                const matchesSearch = name.includes(query);

                if (matchesStatus && matchesSearch) {
                    item.style.display = 'block';
                } else {
                    item.style.display = 'none';
                }
            });
        }

        function escapeHtml(text) {
            return text
                .replace(/&/g, "&amp;")
                .replace(/</g, "&lt;")
                .replace(/>/g, "&gt;")
                .replace(/"/g, "&quot;")
                .replace(/'/g, "&#039;");
        }

        // Auto-run when the page loads
        window.addEventListener('DOMContentLoaded', runTests);
    </script>
</body>
</html>
`
