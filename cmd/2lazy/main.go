package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/goccy/go-yaml"
	"github.com/google/shlex"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

const lazyFileName = "2lazy.yml"

type lazyConfig struct {
	Quiet             bool
	StartInProjectDir bool   `yaml:"start_in_project_dir"`
	ProjectDir        string `yaml:"project_dir"`
	Commands          map[string]string
}

func findConfig() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	currentPath := path.Join(cwd, lazyFileName)

	for true {
		log.WithField("currentPath", currentPath).Debug("Testing path")

		if _, err := os.Stat(currentPath); err == nil {
			log.WithField("path", currentPath).Debug("Found config file")
			return currentPath, nil
		}
		// @TODO consider adding support for other operating systems
		if currentPath == path.Join("/", lazyFileName) {
			return "", fmt.Errorf("Unable to find %s file", lazyFileName)
		}

		currentPath = path.Join(path.Dir(currentPath), "..", lazyFileName)
	}

	return "", fmt.Errorf("Impossible confition occured")
}

func parseConfig() (lazyConfig, error) {
	var config lazyConfig

	configPath, err := findConfig()
	if err != nil {
		return config, err
	}

	configBytes, err := ioutil.ReadFile(configPath)
	if err != nil {
		return config, err
	}

	if err = yaml.Unmarshal(configBytes, &config); err != nil {
		return config, err
	}

	if config.ProjectDir == "" {
		config.ProjectDir = path.Dir(configPath)
	}

	return config, nil
}

func getAvailableAliases(config lazyConfig) []string {
	l := make([]string, len(config.Commands))

	i := 0
	for k := range config.Commands {
		l[i] = k
		i++
	}

	return l
}

func executeCommand(config lazyConfig, alias string, args []string) error {
	realCommand, ok := config.Commands[alias]
	if !ok {
		available := strings.Join(getAvailableAliases(config), ", ")
		return fmt.Errorf("Unknown alias %s. Available aliases: %s", alias, available)
	}

	cmdSplit, err := shlex.Split(realCommand)
	if err != nil {
		return err
	}

	fullArgs := append(cmdSplit[1:], args...)

	cmd := exec.Command(cmdSplit[0], fullArgs...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin

	if config.StartInProjectDir {
		cmd.Dir = config.ProjectDir
	}

	startTime := time.Now()

	if err := cmd.Run(); err != nil {
		exitError, ok := err.(*exec.ExitError)
		if ok {
			exitCode := exitError.ExitCode()

			if exitCode != 0 {
				log.WithField("exitCode", exitCode).Info("Process exited with an exit code not equal to zero")
			}

			return cli.Exit("", exitCode)
		}

		return err
	}

	endTime := time.Now()

	log.WithField("elapsedTime", endTime.Sub(startTime)).Info("Finished")

	return nil
}

func prepareApp() *cli.App {
	app := &cli.App{
		HideHelpCommand: true,
		Usage:           "when you're too lazy to type full command",
		ArgsUsage:       "[alias name] [alias arguments] [...]",

		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "debug",
				Usage: "whether to show debug messages",
			},
			&cli.BoolFlag{
				Name:  "quiet",
				Usage: "whether to hide info messages",
			},
		},
		Before: func(c *cli.Context) error {
			log.SetFormatter(&log.TextFormatter{
				DisableLevelTruncation: true,
				DisableTimestamp:       true,
			})

			if c.Bool("quiet") {
				log.SetLevel(log.WarnLevel)
			} else if c.Bool("debug") {
				log.SetLevel(log.DebugLevel)
			} else {
				log.SetLevel(log.InfoLevel)
			}

			return nil
		},
		Action: func(c *cli.Context) error {
			if !c.Args().Present() {
				cli.ShowAppHelp(c)
				return cli.Exit("", 1)
			}

			config, err := parseConfig()
			if err != nil {
				return err
			}

			if config.Quiet {
				log.SetLevel(log.WarnLevel)
			}

			if config.StartInProjectDir {
				log.WithField("projectDir", config.ProjectDir).Info("Command will be executed in project's directory")
			}

			if err = executeCommand(config, c.Args().First(), c.Args().Tail()); err != nil {
				return err
			}

			return nil
		},
	}

	return app
}

func main() {
	app := prepareApp()

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
