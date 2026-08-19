package playit

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/charmbracelet/huh"
)

const apiBase = "https://api.playit.gg"

func PromptEnable() bool {
	var enabled bool
	_ = huh.NewConfirm().
		Title("Expose this server to the internet via playit.gg?").
		Description("Free, no port forwarding or router changes needed. Mainly useful if\n" +
			"you're running this on a home network - if this machine already has\n" +
			"a public IP (a VPS/dedicated server), you probably don't need this.\n" +
			"Requires a free playit.gg account and a one-time browser approval\n" +
			"(any device works - doesn't have to be this machine, so this is fine\n" +
			"on headless/SSH-only setups too).").
		Value(&enabled).
		Run()
	return enabled
}

func AgentInstalled() bool {
	_, err := exec.LookPath("playit")
	return err == nil
}

func PrintInstallInstructions() {
	fmt.Println("\nThe playit agent isn't installed on this system yet.")

	switch runtime.GOOS {
	case "linux":
		fmt.Println("Install it with:")
		fmt.Println(`  curl -SsL https://playit-cloud.github.io/ppa/key.gpg | gpg --dearmor | sudo tee /etc/apt/trusted.gpg.d/playit.gpg >/dev/null`)
		fmt.Println(`  echo "deb [signed-by=/etc/apt/trusted.gpg.d/playit.gpg] https://playit-cloud.github.io/ppa/data ./" | sudo tee /etc/apt/sources.list.d/playit-cloud.list`)
		fmt.Println(`  sudo apt update && sudo apt install playit`)
	default: // darwin, windows
		fmt.Println("Download the build for your system from:")
		fmt.Println("  https://github.com/playit-cloud/playit-agent/releases/latest")
	}

	fmt.Println("\nRun Deblock again after installing it to set up the tunnel.")
}

func Claim() (string, error) {
	code, err := generateCode()
	if err != nil {
		return "", fmt.Errorf("couldn't generate a claim code: %w", err)
	}

	if err := postJSON("/claim/setup", claimSetupBody(code), nil); err != nil {
		return "", fmt.Errorf("couldn't start the claim: %w", err)
	}

	claimURL := "https://playit.gg/claim/" + code
	fmt.Println("\nOpen this URL to link this server with your playit.gg account:")
	fmt.Println("  " + claimURL)
	fmt.Println("Trying to open it in a local browser now, but it doesn't have to be this")
	fmt.Println("machine - if you're on a headless server (SSH-only, no browser here),")
	fmt.Println("just copy that URL to your phone or laptop instead. Deblock will keep")
	fmt.Println("waiting either way.")
	openBrowser(claimURL)

	fmt.Println("\nWaiting for you to approve it (up to 2 minutes)...")

	const attempts = 40
	const interval = 3 * time.Second
	for i := 0; i < attempts; i++ {
		time.Sleep(interval)

		var resp claimSetupResponse
		if err := postJSON("/claim/setup", claimSetupBody(code), &resp); err != nil {
			continue
		}

		switch resp.Data {
		case "UserAccepted":
			return exchangeSecret(code)
		case "UserRejected":
			return "", fmt.Errorf("the claim was rejected in the browser")
		}
	}

	return "", fmt.Errorf("timed out waiting for approval in the browser")
}

func FinishSetup(secretKey string, port int) {
	fmt.Println("\nDeblock linked this server to your playit.gg account.")

	if runtime.GOOS == "linux" && runLinuxSetup(secretKey) {
		fmt.Println(okMsg())
	} else {
		printManualAgentSteps(secretKey)
	}

	dashboardURL := "https://playit.gg/account/agents"
	fmt.Printf("\nNow create a tunnel pointing at local port %d (protocol TCP, type \"Minecraft Java\")\n", port)
	fmt.Println("on this page:")
	fmt.Println("  " + dashboardURL)
	fmt.Println("Trying to open it in a local browser now, but same as before - any device")
	fmt.Println("works, this doesn't have to be the machine running the server.")
	openBrowser(dashboardURL)
}

func okMsg() string {
	return "The agent is configured and running."
}

func runLinuxSetup(secretKey string) bool {
	fmt.Println("Pointing the agent at this secret (you may be asked for your sudo password)...")

	write := exec.Command("sudo", "tee", "/etc/playit/playit.toml")
	write.Stdin = strings.NewReader(fmt.Sprintf("secret_key = %q\n", secretKey))
	write.Stdout = io.Discard
	write.Stderr = os.Stderr
	if err := write.Run(); err != nil {
		fmt.Println("Couldn't write the agent's config automatically:", err)
		return false
	}

	fmt.Println("Starting the agent...")
	start := exec.Command("sudo", "playit", "start")
	start.Stdin = os.Stdin
	start.Stdout = os.Stdout
	start.Stderr = os.Stderr
	if err := start.Run(); err != nil {
		fmt.Println("Couldn't start the agent automatically:", err)
		return false
	}

	return true
}

func printManualAgentSteps(secretKey string) {
	switch runtime.GOOS {
	case "linux":
		fmt.Println("Here's what to run by hand instead:")
		fmt.Printf("  echo 'secret_key = \"%s\"' | sudo tee /etc/playit/playit.toml\n", secretKey)
		fmt.Println("  sudo playit start")
	default:
		fmt.Println("Run the playit agent with this secret key (see the agent's own docs for")
		fmt.Println("how to provide it on your system):")
		fmt.Printf("  %s\n", secretKey)
	}
}

func exchangeSecret(code string) (string, error) {
	var resp claimExchangeResponse
	if err := postJSON("/claim/exchange", map[string]any{"code": code}, &resp); err != nil {
		return "", fmt.Errorf("couldn't exchange the claim for a secret: %w", err)
	}
	if resp.Data.SecretKey == "" {
		return "", fmt.Errorf("playit.gg didn't return a secret key")
	}
	return resp.Data.SecretKey, nil
}

func claimSetupBody(code string) map[string]any {
	return map[string]any{
		"code":       code,
		"agent_type": "self-managed",
		"version":    "deblock",
	}
}

type claimSetupResponse struct {
	Status string `json:"status"`
	Data   string `json:"data"`
}

type claimExchangeResponse struct {
	Status string `json:"status"`
	Data   struct {
		SecretKey string `json:"secret_key"`
	} `json:"data"`
}

func postJSON(path string, body any, out any) error {
	buf, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, apiBase+path, bytes.NewReader(buf))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(raw))
	}
	if out == nil {
		return nil
	}
	return json.Unmarshal(raw, out)
}

func generateCode() (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	code := make([]byte, 12)
	for i := range code {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		code[i] = charset[n.Int64()]
	}
	return string(code), nil
}

func openBrowser(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	_ = cmd.Start()
}
