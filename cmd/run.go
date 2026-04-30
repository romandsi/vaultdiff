package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourorg/vaultdiff/internal/config"
	"github.com/yourorg/vaultdiff/internal/diff"
	"github.com/yourorg/vaultdiff/internal/output"
	"github.com/yourorg/vaultdiff/internal/vault"
)

var (
	cfgFile    string
	format     string
	maskValues bool
	failOnDiff bool
)

var rootCmd = &cobra.Command{
	Use:   "vaultdiff",
	Short: "Compare secrets across HashiCorp Vault environments",
	RunE:  runDiff,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&cfgFile, "config", "c", "vaultdiff.yaml", "path to config file")
	rootCmd.Flags().StringVarP(&format, "format", "f", "text", "output format: text, json, markdown")
	rootCmd.Flags().BoolVar(&maskValues, "mask", false, "mask secret values in output")
	rootCmd.Flags().BoolVar(&failOnDiff, "fail-on-diff", false, "exit with code 1 if differences are found")
}

func runDiff(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	if err := config.Validate(cfg); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	fmt, err := output.ParseFormat(format)
	if err != nil {
		return err
	}

	clients := make(map[string]*vault.Client, len(cfg.Environments))
	for _, env := range cfg.Environments {
		c, err := vault.NewClient(env.Address, env.Token)
		if err != nil {
			return fmt.Errorf("creating client for %s: %w", env.Name, err)
		}
		clients[env.Name] = c
	}

	var allResults []diff.Result
	envA := cfg.Environments[0]
	envB := cfg.Environments[1]

	for _, path := range cfg.Paths {
		secretsA, err := clients[envA.Name].ReadSecret(path)
		if err != nil {
			return fmt.Errorf("reading %s from %s: %w", path, envA.Name, err)
		}
		secretsB, err := clients[envB.Name].ReadSecret(path)
		if err != nil {
			return fmt.Errorf("reading %s from %s: %w", path, envB.Name, err)
		}
		results := diff.Compare(path, envA.Name, envB.Name, secretsA, secretsB)
		allResults = append(allResults, results...)
	}

	if err := output.Write(os.Stdout, fmt, allResults, maskValues); err != nil {
		return err
	}

	os.Exit(output.ResolveExitCode(allResults, failOnDiff))
	return nil
}
