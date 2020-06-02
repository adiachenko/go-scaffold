package cmd

import (
	"fmt"
	"os"
	"strings"

	"adiachenko/go-scaffold/internal/platform/conf"
	"adiachenko/go-scaffold/internal/platform/jaeger"

	"github.com/Vinelab/go-reporting"
	"github.com/Vinelab/go-reporting/sentry"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var tracedCommands []string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go-scaffold",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if StringInSlice(cmd.Name(), tracedCommands) {
			span := jaeger.Trace.StartSpan(cmd.Name(), jaeger.Trace.EmptySpanContext())

			span.Tag("type", "cli")
			span.Tag("argv", strings.Join(args, "\n"))
		}
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if jaeger.Trace.RootSpan() != nil {
			jaeger.Trace.RootSpan().Finish()

			if err := jaeger.Trace.Close(); err != nil {
				logrus.Error(err)
			}
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(printActiveEnvironment, initSentry)

	tracedCommands = []string{}
}

func printActiveEnvironment() {
	logrus.Infof("Environment: %s", conf.AppENV)
}

func initSentry() {
	defer reporting.LogPanic()

	if conf.AppDebug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	_ = reporting.RegisterSentry(
		logrus.WarnLevel,
		sentry.Options{},
		sentry.TagInjector{
			Tags: func() map[string]string {
				return map[string]string{
					"uuid": jaeger.Trace.UUID(),
				}
			},
		},
	)
}

// StringInSlice checks whether list contains
func StringInSlice(str string, list []string) bool {
	for _, item := range list {
		if item == str {
			return true
		}
	}
	return false
}
