package operations

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/noahc3/hatch-cli/src/log"
	"github.com/noahc3/hatch-cli/src/types"
	"github.com/noahc3/hatch-cli/src/utils"

	dtypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
)

func createDockerClient() (*client.Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}

	return cli, nil
}

func flattenEnvVariables(env map[string]string) []string {
	var envs []string

	for key, value := range env {
		envs = append(envs, fmt.Sprintf("%s=%s", key, value))
	}

	return envs
}

func pullDockerImage(client *client.Client, image string) error {
	log.Info("  Querying image '%s'...\n", image)

	ctx := context.Background()
	r, err := client.ImagePull(ctx, image, dtypes.ImagePullOptions{All: false})
	if err != nil {
		return err
	}

	defer r.Close()

	log.Info("  Pulling container image (this may take a while)...\n")

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		log.Debug(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func runInstallContainer(egg *types.Egg, client *client.Client) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conf := &container.Config{
		Hostname:     "crack_egg_installer",
		AttachStdout: true,
		AttachStderr: true,
		AttachStdin:  true,
		OpenStdin:    true,
		Tty:          true,
		Cmd:          []string{egg.Scripts.Installation.Entrypoint, "/mnt/install/install.sh"},
		Image:        egg.Scripts.Installation.Container,
		Env:          flattenEnvVariables(egg.CrackProps.EnvVariables),
		Labels: map[string]string{
			"Service":       "Crack (Pterodactyl Egg Installer)",
			"ContainerType": "crack_egg_installer",
		},
	}

	tmpDir := filepath.Dir(egg.CrackProps.InstallScriptPath)

	hostConf := &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Target:   "/mnt/server",
				Source:   egg.CrackProps.InstallDirPath,
				Type:     mount.TypeBind,
				ReadOnly: false,
			},
			{
				Target:   "/mnt/install",
				Source:   tmpDir,
				Type:     mount.TypeBind,
				ReadOnly: false,
			},
		},
		Tmpfs: map[string]string{
			"/tmp": "rw,exec,nosuid,size=100M",
		},
		DNS: []string{"1.1.1.1", "1.0.0.1"},
	}

	// defer func() {
	// 	err := os.RemoveAll(tmpDir)
	// 	if !os.IsNotExist(err) {
	// 		log.Error("  Failed to remove temporary directory: %s\n", err)
	// 	}
	// }()

	log.Info("  Creating container...\n")

	r, err := client.ContainerCreate(ctx, conf, hostConf, nil, nil, fmt.Sprintf("crack_egg_installer_%s", filepath.Base(tmpDir)))
	if err != nil {
		return err
	}

	log.Info("  Running container with install script...\n")
	err = client.ContainerStart(ctx, r.ID, container.StartOptions{})
	if err != nil {
		return err
	}

	sChan, eChan := client.ContainerWait(ctx, r.ID, container.WaitConditionNotRunning)
	select {
	case err := <-eChan:
		if err == nil {
			log.Info("  Install complete.\n")
		} else {
			return err
		}
	case <-sChan:
	}

	return nil
}

func exportScript(egg *types.Egg) (string, error) {
	temp, err := os.MkdirTemp("", "hatch-egg-*")
	if err != nil {
		return "", err
	}

	scriptPath := temp + "/install.sh"

	log.Info("  Exporting script to %s\n", scriptPath)

	scriptFile, err := os.OpenFile(scriptPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return "", err
	}

	defer scriptFile.Close()

	cleanScript := strings.ReplaceAll(egg.Scripts.Installation.Script, "\r\n", "\n")
	_, err = scriptFile.WriteString(cleanScript)
	if err != nil {
		return "", err
	}

	return scriptPath, nil
}

func Install(egg *types.Egg) error {
	//check if dir exists
	dir := egg.CrackProps.InstallDirPath
	outdirinfo, err := os.Stat(dir)
	if err != nil {
		return log.Error("Output directory does not exist\n")
	}

	if !outdirinfo.IsDir() {
		return log.Error("Output directory is not a directory\n")
	}

	vars, err := GetEnvVariables(egg)
	if err != nil {
		return err
	}

	egg.CrackProps.EnvVariables = vars

	log.Nl()

	log.Info("Installation summary:\n")
	log.Info("  Egg: %s by %s\n", egg.Name, egg.Author)
	log.Info("  Description: %s\n", egg.Description)
	log.Info("  Installation plan\n")
	log.Info("    Install container: %s\n", egg.Scripts.Installation.Container)
	log.Info("    Install entrypoint: %s\n", egg.Scripts.Installation.Entrypoint)
	log.Info("    Output directory: %s\n", dir)
	log.Info("  Container environment variables:\n")

	for _, variable := range egg.Variables {
		log.Info("    %s=%s\n", variable.EnvVariable, vars[variable.EnvVariable])
	}

	log.Info("\nIMPORTANT: Review the above summary before continuing!\n")
	res := utils.YesNoPrompt("Do you want to continue with the installation?")

	if !res {
		return log.Error("Installation aborted")
	}

	log.Info("Starting installation...\n")

	scriptPath, err := exportScript(egg)
	if err != nil {
		return log.Error("  Failed to export installation script: %s\n", err)
	}

	egg.CrackProps.InstallScriptPath = scriptPath

	cli, err := createDockerClient()
	if err != nil {
		return log.Error("  Failed to create Docker client: %s\n", err)
	}

	err = pullDockerImage(cli, egg.Scripts.Installation.Container)
	if err != nil {
		return log.Error("  Failed to pull Docker image: %s\n", err)
	}

	err = runInstallContainer(egg, cli)
	if err != nil {
		return log.Error("  Failed to run installation container: %s\n", err)
	}

	log.Info("Installation complete\n")

	return nil
}
