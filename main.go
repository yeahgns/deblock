package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"

	"mcserverwizard/internal/download"
	"mcserverwizard/internal/loaders"
	"mcserverwizard/internal/props"
	"mcserverwizard/internal/startscript"
)

var (
	titleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("42"))
	errStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("196"))
	okStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("42"))
	dimStyle   = lipgloss.NewStyle().Faint(true)
)

const banner = `
██████╗ ███████╗██████╗ ██╗      ██████╗  ██████╗██╗  ██╗
██╔══██╗██╔════╝██╔══██╗██║     ██╔═══██╗██╔════╝██║ ██╔╝
██║  ██║█████╗  ██████╔╝██║     ██║   ██║██║     █████╔╝
██║  ██║██╔══╝  ██╔══██╗██║     ██║   ██║██║     ██╔═██╗
██████╔╝███████╗██████╔╝███████╗╚██████╔╝╚██████╗██║  ██╗
╚═════╝ ╚══════╝╚═════╝ ╚══════╝ ╚═════╝  ╚═════╝╚═╝  ╚═╝
        Minecraft server, unblocked. — by yeahgns
`

var validNameRe = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

func main() {
	fmt.Println(titleStyle.Render(banner))

	cwd, err := os.Getwd()
	fatalIf(err)

	name := "my-minecraft-server"
	fatalIf(huh.NewInput().
		Title("Server name (It turns into a folder with that name, created right here)").
		Value(&name).
		Validate(validateName).
		Run())

	serverDir := filepath.Join(cwd, name)
	propsPath := filepath.Join(serverDir, "server.properties")

	if _, err := os.Stat(propsPath); err == nil {
		handleExistingServer(serverDir, propsPath)
		return
	}

	runFreshInstall(serverDir)
}

func validateName(s string) error {
	if s == "" {
		return fmt.Errorf("The name field cannot be left blank.")
	}
	if !validNameRe.MatchString(s) {
		return fmt.Errorf("Use only letters, numbers, hyphens, and underscores.")
	}
	return nil
}

func runFreshInstall(serverDir string) {
	// For now, only Vanilla is available in the setup menu
	loader := loaders.Vanilla
	fmt.Println(dimStyle.Render("\nLoader: Vanilla (Paper e Fabric in the next versions)"))

	version := chooseVersion(loader)

	fmt.Printf("\nLooking for the server .jar file (%s %s)...\n", loader, version)
	jar, err := loaders.Resolve(loader, version)
	fatalIf(err)
	fmt.Println(okStyle.Render("  Found: " + jar.Label))

	cfg := props.Default()
	memory := "2G"
	askServerConfig(&cfg, &memory)

	if !confirmEULA() {
		fmt.Println(errStyle.Render("EULA not accepted – Installation cancelled."))
		return
	}

	fatalIf(os.MkdirAll(serverDir, 0o755))

	jarPath := filepath.Join(serverDir, "server.jar")
	fmt.Printf("\nDownloading %s...\n", jar.Label)
	fatalIf(download.File(jar.URL, jarPath))

	fatalIf(os.WriteFile(filepath.Join(serverDir, "eula.txt"), []byte("eula=true\n"), 0o644))
	fatalIf(props.Write(filepath.Join(serverDir, "server.properties"), cfg))
	fatalIf(startscript.Write(serverDir, memory))

	fmt.Println(okStyle.Render("\nServer ready in: ") + serverDir)
	fmt.Println(dimStyle.Render("  (Run the wizard again in that same folder to reconfigure without downloading everything all over again.)"))

	offerToRun(serverDir, memory)
}

func handleExistingServer(serverDir, propsPath string) {
	fmt.Println(dimStyle.Render("\nThere is already a server configured in " + serverDir))

	var action string
	fatalIf(huh.NewSelect[string]().
		Title("What do you want to do?").
		Options(
			huh.NewOption("Edit settings (server.properties)", "reconfigure"),
			huh.NewOption("Reinstall from scratch (deletes .jar file and downloads it again)", "reinstall"),
			huh.NewOption("Start the existing server", "start"),
		).
		Value(&action).
		Run())

	switch action {
	case "start":
		offerToRun(serverDir, "2G")
		return
	case "reinstall":
		var sure bool
		fatalIf(huh.NewConfirm().
			Title("This will delete the current .jar and server.properties files. Continue?").
			Value(&sure).
			Run())
		if !sure {
			fmt.Println("Aborted.")
			return
		}
		runFreshInstall(serverDir)
		return
	default: // reconfigure
		cfg, _ := props.Read(propsPath)
		memory := "2G"
		askServerConfig(&cfg, &memory)
		fatalIf(props.Write(propsPath, cfg))
		fatalIf(startscript.Write(serverDir, memory))
		fmt.Println(okStyle.Render("\nConfiguration updated in: ") + propsPath)
		offerToRun(serverDir, memory)
	}
}

func chooseVersion(loader string) string {
	releases, err := loaders.Releases(loader, 15)
	if err != nil {
		fmt.Println(errStyle.Render("Unable to automatically retrieve the list of versions: " + err.Error()))
		var version string
		fatalIf(huh.NewInput().
			Title("Enter the Minecraft version manually (e.g., 1.21.1)").
			Value(&version).
			Run())
		return version
	}

	options := make([]huh.Option[string], 0, len(releases)+1)
	options = append(options, huh.NewOption(releases[0]+"  (latest)", releases[0]))
	for _, r := range releases[1:] {
		options = append(options, huh.NewOption(r, r))
	}
	options = append(options, huh.NewOption("Another version (type manually)", "__custom__"))

	var version string
	fatalIf(huh.NewSelect[string]().
		Title("Which version of Minecraft?").
		Options(options...).
		Value(&version).
		Run())

	if version == "__custom__" {
		fatalIf(huh.NewInput().
			Title("Enter the Minecraft version (e.g., 1.21.1)").
			Value(&version).
			Run())
	}
	return version
}

func askServerConfig(cfg *props.Config, memory *string) {
	maxPlayersStr := strconv.Itoa(cfg.MaxPlayers)
	portStr := strconv.Itoa(cfg.Port)

	fatalIf(huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("MOTD (message displayed in the server list)").Value(&cfg.MOTD),
			huh.NewInput().Title("Maximum number of players").Value(&maxPlayersStr).Validate(validatePositiveInt),
			huh.NewInput().Title("Server port").Value(&portStr).Validate(validatePositiveInt),
			huh.NewSelect[string]().Title("Difficulty").
				Options(
					huh.NewOption("Peaceful", "peaceful"),
					huh.NewOption("Easy", "easy"),
					huh.NewOption("Normal", "normal"),
					huh.NewOption("Hard", "hard"),
				).Value(&cfg.Difficulty),
			huh.NewSelect[string]().Title("Gamemode").
				Options(
					huh.NewOption("Survival", "survival"),
					huh.NewOption("Creative", "creative"),
					huh.NewOption("Adventure", "adventure"),
					huh.NewOption("Spectator", "spectator"),
				).Value(&cfg.Gamemode),
			huh.NewConfirm().Title("Enable whitelist?").Value(&cfg.Whitelist),
			huh.NewConfirm().Title("Require original Microsoft/Mojang account (online-mode)?").Value(&cfg.OnlineMode),
			huh.NewInput().Title("Ammount of memory allocated to the server (e.g. 2G, 4G)").Value(memory),
		),
	).Run())

	if n, err := strconv.Atoi(maxPlayersStr); err == nil {
		cfg.MaxPlayers = n
	}
	if n, err := strconv.Atoi(portStr); err == nil {
		cfg.Port = n
	}
}

func validatePositiveInt(s string) error {
	n, err := strconv.Atoi(s)
	if err != nil || n <= 0 {
		return fmt.Errorf("Enter a valid number greater than zero.")
	}
	return nil
}

func confirmEULA() bool {
	fmt.Println(dimStyle.Render("\nTo run a Minecraft server, you need to accept Mojang's EULA:"))
	fmt.Println(dimStyle.Render("  https://aka.ms/MinecraftEULA"))

	var accepted bool
	fatalIf(huh.NewConfirm().
		Title("Have you read and do you accept the terms of Mojang's EULA?").
		Value(&accepted).
		Run())
	return accepted
}

func offerToRun(serverDir, memory string) {
	var runNow bool
	fatalIf(huh.NewConfirm().
		Title("Do you want to start the server now?").
		Value(&runNow).
		Run())

	if !runNow {
		fmt.Println(dimStyle.Render("\nTo run it later: go to the folder and execute ./start.sh (Linux/Mac) or start.bat (Windows)."))
		return
	}

	if _, err := exec.LookPath("java"); err != nil {
		fmt.Println(errStyle.Render("\nI didn't find Java installed on this machine."))
		fmt.Println(dimStyle.Render("Install a JDK (e.g., https://adoptium.net) and then run ./start.sh or start.bat in the server folder."))
		return
	}

	fmt.Println(okStyle.Render("\nStarting the server (Ctrl+C to stop)...\n"))
	cmd := exec.Command("java", fmt.Sprintf("-Xmx%s", memory), fmt.Sprintf("-Xms%s", memory), "-jar", "server.jar", "nogui")
	cmd.Dir = serverDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		fmt.Println(errStyle.Render("The server terminated with an error: " + err.Error()))
	}
}

func fatalIf(err error) {
	if err != nil {
		fmt.Println(errStyle.Render("Erro: " + err.Error()))
		os.Exit(1)
	}
}
