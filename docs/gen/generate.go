package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/docker/docker/cli/command"
	"github.com/docker/docker/cli/command/commands"
	"github.com/docker/docker/pkg/term"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"github.com/spf13/pflag"
)

const descriptionSourcePath = "man/src/"

func generateManPages(opts *options) error {
	cmd, err := loadCommands(opts.source)
	if err != nil {
		return err
	}

	header := &doc.GenManHeader{
		Title:   "DOCKER",
		Section: "1",
		Source:  "Docker Community",
	}

	return doc.GenManTreeFromOpts(cmd, doc.GenManTreeOptions{
		Header:           header,
		Path:             opts.manTarget,
		CommandSeparator: "-",
	})
}

func generateCliReferenceDocs(opts *options) error {
	cmd, err := loadCommands(opts.source)
	if err != nil {
		return err
	}

	return doc.GenMarkdownTree(cmd, opts.cliTarget)
}

func loadCommands(source string) (*cobra.Command, error) {
	stdin, stdout, stderr := term.StdStreams()
	dockerCli := command.NewDockerCli(stdin, stdout, stderr)
	cmd := &cobra.Command{Use: "docker"}
	commands.AddCommands(cmd, dockerCli)
	source = filepath.Join(source, descriptionSourcePath)
	if err := loadLongDescription(cmd, source); err != nil {
		return cmd, err
	}
	cmd.DisableAutoGenTag = true
	return cmd, nil
}

func loadLongDescription(cmd *cobra.Command, path string) error {
	for _, cmd := range cmd.Commands() {
		if cmd.Name() == "" {
			continue
		}
		fullpath := filepath.Join(path, cmd.Name()+".md")

		if cmd.HasSubCommands() {
			loadLongDescription(cmd, filepath.Join(path, cmd.Name()))
		}

		if _, err := os.Stat(fullpath); err != nil {
			log.Printf("WARN: %s does not exist, skipping\n", fullpath)
			continue
		}

		content, err := ioutil.ReadFile(fullpath)
		if err != nil {
			return err
		}
		cmd.Long = string(content)
	}
	return nil
}

type options struct {
	source    string
	manTarget string
	cliTarget string
}

func parseArgs() (*options, error) {
	opts := &options{}
	cwd, _ := os.Getwd()
	flags := pflag.NewFlagSet(os.Args[0], pflag.ContinueOnError)
	flags.StringVar(&opts.source, "root", cwd, "Path to project root")
	flags.StringVar(&opts.manTarget, "man", "", "Target path for generated man pages")
	flags.StringVar(&opts.cliTarget, "cli", "", "Target path for generated cli reference")
	err := flags.Parse(os.Args[1:])
	return opts, err
}

func main() {
	opts, err := parseArgs()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
	fmt.Printf("Project root: %s\n", opts.source)
	if opts.manTarget != "" {
		fmt.Printf("Generating man pages into %s\n", opts.manTarget)
		if err := generateManPages(opts); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to generate man pages: %s\n", err.Error())
			os.Exit(2)
		}
	}
	if opts.cliTarget != "" {
		fmt.Printf("Generating cli reference into %s\n", opts.cliTarget)
		if err := generateCliReferenceDocs(opts); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to generate cli reference: %s\n", err.Error())
			os.Exit(2)
		}
	}
	if opts.manTarget == "" && opts.cliTarget == "" {
		fmt.Fprintln(os.Stderr, "Nothing to do")
		os.Exit(1)
	}
}
