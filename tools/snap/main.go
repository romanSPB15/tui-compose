package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"image/png"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/romanSPB15/tui-compose/v3/cell"
	"github.com/romanSPB15/tui-compose/v3/tools/snap/render"
)

func main() {
	pkg := flag.String("pkg", ".", "пакет для сборки")
	out := flag.String("o", "screenshot.png", "выходной файл")
	w := flag.Int("w", 120, "ширина TUI")
	h := flag.Int("h", 30, "высота TUI")
	bin := flag.String("bin", "", "готовый бинарник вместо сборки")
	all := flag.String("all", "", "обработать все подпакеты, напр. ./examples/...")
	flag.Parse()

	if *all != "" {
		runAll(*all, *w, *h, *out)
		return
	}

	if err := capture(*pkg, *bin, *out, *w, *h); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func runAll(pattern string, w, h int, out string) {
	root := strings.TrimSuffix(pattern, "/...")
	entries, err := os.ReadDir(root)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		pkgPath := root + "/" + e.Name()
		if _, err := os.Stat(filepath.Join(pkgPath, "go.mod")); err == nil {
			continue
		}
		target := filepath.Join(pkgPath, out)
		fmt.Println("→", pkgPath)
		if err := capture(pkgPath, "", target, w, h); err != nil {
			fmt.Fprintln(os.Stderr, "  ", err)
		}
	}
}

func capture(pkg, bin, out string, w, h int) error {
	exe := bin
	if exe == "" {
		exe = fmt.Sprintf("./capture-app-%d", time.Now().UnixNano())
		if runtime.GOOS == "windows" {
			exe += ".exe"
		}
		defer os.Remove(exe)

		buildPath := pkg
		if !strings.HasPrefix(buildPath, ".") && !filepath.IsAbs(buildPath) {
			buildPath = "./" + buildPath
		}
		build := exec.Command("go", "build", "-tags", "capture", "-o", exe, buildPath)
		build.Stderr = os.Stderr
		if err := build.Run(); err != nil {
			return fmt.Errorf("build: %w", err)
		}
	}

	cmd := exec.Command(exe)
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("TUI_WIDTH=%d", w),
		fmt.Sprintf("TUI_HEIGHT=%d", h),
	)
	cmd.Stderr = os.Stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	defer func() {
		_ = cmd.Process.Kill()
		_ = cmd.Wait()
	}()

	var cells [][]cell.Cell
	sc := bufio.NewScanner(stdout)
	sc.Buffer(make([]byte, 1<<20), 1<<24)
	for sc.Scan() {
		var frame [][]cell.Cell
		if err := json.Unmarshal(sc.Bytes(), &frame); err == nil {
			cells = frame
			break
		}
	}
	if cells == nil {
		return fmt.Errorf("no frame captured")
	}

	img := render.Frame(cells)
	f, err := os.Create(out)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		return err
	}
	fmt.Println("saved", out)
	return nil
}
