package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"

	"github.com/gonutz/w32"
	"github.com/google/uuid"
	"github.com/zserge/lorca"
)

var sessionToken string = "su-" + uuid.New().String()

func runServer(ctx context.Context) *exec.Cmd {
	// cmd := exec.CommandContext(ctx, "focalboard-server.exe", "--monitorpid", strconv.FormatInt(int64(os.Getpid()), 10), "-single-user")
	cmd := exec.CommandContext(ctx, "focalboard-server.exe", "-single-user")
	// cmd := exec.CommandContext(ctx, "cmd.exe", "/C", "start", "./bin/focalboard-server.exe", "--monitorpid", strconv.FormatInt(int64(os.Getpid()), 10))
	// cmd := exec.CommandContext(ctx, "cmd.exe", "/C", "start", "./bin/focalboard-server.exe")

	// cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd.Env = []string{fmt.Sprintf("FOCALBOARD_SINGLE_USER_TOKEN=%s", sessionToken)}
	cmd.Stdout = os.Stdout
	go func() {
		err := cmd.Run()
		if err != nil {
			log.Println("Failed to start server")
			log.Fatal(err)
		}
		log.Printf("Just ran subprocess %d, exiting\n", cmd.Process.Pid)
	}()

	return cmd
}

func main() {
	// log.Printf("PID: %s", strconv.FormatInt(int64(os.Getpid()), 10))
	hideConsole()

	// Try to find Chrome if Lorca can't find it
	if len(lorca.ChromeExecutable()) == 0 {
		chromePath := locateChrome()
		log.Printf("chromePath: %s", chromePath)
		if len(chromePath) > 0 {
			os.Setenv("LORCACHROME", chromePath)
		} else {
			lorca.PromptDownload()
			log.Fatal("Chrome not installed")
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	cmd := runServer(ctx)

	defer func() {
		log.Println("Cleanup")
		cancel()
		if err := cmd.Process.Kill(); err != nil {
			log.Fatal("failed to kill server process: ", err)
		}

		if r := recover(); r != nil {
			log.Fatal("ERROR: ", r)
		}
	}()

	ui, err := lorca.New("", "", 1024, 768)
	if err != nil {
		log.Panic(err)
	}
	defer ui.Close()

	if err := ui.Load("http://localhost:8088"); err != nil {
		log.Panic(err)
	}

	script := fmt.Sprintf("localStorage.setItem('sessionId', '%s');", sessionToken)
	value := ui.Eval(script)
	if err := value.Err(); err != nil {
		log.Panic(err)
	}

	log.Printf("Started")
	<-ui.Done()

	log.Println("App Closed")
}

func hideConsole() {
	console := w32.GetConsoleWindow()
	if console != 0 {
		_, consoleProcID := w32.GetWindowThreadProcessId(console)
		if w32.GetCurrentProcessId() == consoleProcID {
			w32.ShowWindowAsync(console, w32.SW_HIDE)
		}
	}
}

// This duplicates the logic in Lorca, but adds Edge as an option for Windows, fallback to standard logic for other OSes
func locateChrome() string {
	var paths []string
	switch runtime.GOOS {
	// case "darwin":
	// 	paths = []string{
	// 		"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
	// 		"/Applications/Google Chrome Canary.app/Contents/MacOS/Google Chrome Canary",
	// 		"/Applications/Chromium.app/Contents/MacOS/Chromium",
	// 		"/usr/bin/google-chrome-stable",
	// 		"/usr/bin/google-chrome",
	// 		"/usr/bin/chromium",
	// 		"/usr/bin/chromium-browser",
	// 	}
	case "windows":
		paths = []string{
			os.Getenv("LocalAppData") + "/Google/Chrome/Application/chrome.exe",
			os.Getenv("ProgramFiles") + "/Google/Chrome/Application/chrome.exe",
			os.Getenv("ProgramFiles(x86)") + "/Google/Chrome/Application/chrome.exe",
			os.Getenv("LocalAppData") + "/Chromium/Application/chrome.exe",
			os.Getenv("ProgramFiles") + "/Chromium/Application/chrome.exe",
			os.Getenv("ProgramFiles(x86)") + "/Chromium/Application/chrome.exe",
			os.Getenv("ProgramFiles(x86)") + "/Microsoft/Edge/Application/msedge.exe",
		}
		// default:
		// 	paths = []string{
		// 		"/usr/bin/google-chrome-stable",
		// 		"/usr/bin/google-chrome",
		// 		"/usr/bin/chromium",
		// 		"/usr/bin/chromium-browser",
		// 		"/snap/bin/chromium",
		// 	}
	}

	for _, path := range paths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue
		}
		return path
	}

	return ""
}
