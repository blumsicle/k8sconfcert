/*
Copyright © 2024 Brian Blumberg <blumsicle@icloud.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	Name    string
	Version string
	Commit  string
)

var rootCmd = &cobra.Command{
	Use:     Name,
	Short:   "Extract out certificate authority data from kubeconfig",
	Version: Version + " " + Commit,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		logLevelStr, _ := cmd.Flags().GetString("log-level")
		logLevel, err := zerolog.ParseLevel(logLevelStr)
		if err != nil {
			return err
		}
		zerolog.SetGlobalLevel(logLevel)
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		kubeconfig, _ := cmd.Flags().GetString("kubeconfig")
		output, _ := cmd.Flags().GetString("output")

		kubeconfig = os.ExpandEnv(kubeconfig)
		output = os.ExpandEnv(output)

		log.Debug().Str("kubeconfig", kubeconfig).Send()
		log.Debug().Str("output", output).Send()

		config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return err
		}

		fd, err := os.Create(output)
		if err != nil {
			return err
		}
		defer fd.Close()

		_, err = fd.Write(config.CAData)
		if err != nil {
			return err
		}

		return nil
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP("log-level", "l", "info", "log level")
	rootCmd.Flags().StringP("kubeconfig", "k", "$HOME/.kube/config", "kubeconfig file")
	rootCmd.Flags().StringP("output", "o", "./ca.crt", "output file")
}
