package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var rowRe = regexp.MustCompile(`\|\s*\[` + "`" + `([^` + "`" + `]+)` + "`" + `\]\([^)]*examples/([^)/]+)[^)]*\)\s*\|\s*([^|]+?)\s*\|`)

const tmpl = `# Пример %s

%s

## 🚀 Запуск

` + "```bash" + `
go run ./examples/%s
` + "```" + `

## 📸 Скриншот

![Screenshot](screenshot.png)

## 📝 Реализация

[main.go](main.go).
`

func main() {
	root := flag.String("root", "examples", "корень с примерами")
	dry := flag.Bool("dry", false, "не писать файлы, только показать")
	flag.Parse()

	readmePath := filepath.Join(*root, "README.md")
	f, err := os.Open(readmePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer f.Close()

	seen := map[string]bool{}
	count := 0

	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 1<<20), 1<<20)
	for sc.Scan() {
		m := rowRe.FindStringSubmatch(sc.Text())
		if m == nil {
			continue
		}
		name := m[2]
		desc := strings.TrimSpace(m[3])
		if seen[name] {
			continue
		}
		seen[name] = true
		count++

		dir := filepath.Join(*root, name)
		if _, err := os.Stat(dir); err != nil {
			fmt.Fprintf(os.Stderr, "skip %s: %v\n", name, err)
			continue
		}

		content := fmt.Sprintf(tmpl, name, desc, name)
		out := filepath.Join(dir, "README.md")

		if *dry {
			fmt.Printf("=== %s ===\n%s\n", out, content)
			continue
		}
		if err := os.WriteFile(out, []byte(content), 0o644); err != nil {
			fmt.Fprintf(os.Stderr, "write %s: %v\n", out, err)
			continue
		}
		fmt.Println("wrote", out)
	}

	if err := sc.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("done: %d\n", count)
}
