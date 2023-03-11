package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/vitorfhc/queryxss/pkg/httpclient"
	"github.com/vitorfhc/queryxss/pkg/reflections"
)

type cliOptions struct {
	Debug         bool
	Silent        bool
	File          string
	Stdin         bool
	Headers       []string
	RateLimit     uint
	NoColor       bool
	MinLength     uint
	AllowInsecure bool
}

var options cliOptions

const banner = `
██████████████████████████████████████████████████
█─▄▄▄─█▄─██─▄█▄─▄▄─█▄─▄▄▀█▄─█─▄█▄─▀─▄█─▄▄▄▄█─▄▄▄▄█
█─██▀─██─██─███─▄█▀██─▄─▄██▄─▄███▀─▀██▄▄▄▄─█▄▄▄▄─█
▀───▄▄▀▀▄▄▄▄▀▀▄▄▄▄▄▀▄▄▀▄▄▀▀▄▄▄▀▀▄▄█▄▄▀▄▄▄▄▄▀▄▄▄▄▄▀

`

func init() {
	options = cliOptions{
		Debug:   false,
		Silent:  false,
		File:    "",
		Stdin:   true,
		Headers: []string{},
	}
}

func main() {
	execute()
}

var rootCmd = &cobra.Command{
	Use:   "queryxss",
	Short: "QueryXSS finds reflected values in the HTTP response.",
	Long:  ``,
	Run:   run,
}

func run(cmd *cobra.Command, args []string) {
	c := make(chan os.Signal, 1)
	defer close(c)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	cmdCtx, cmdCancel := context.WithCancel(context.Background())

	go func() {
		logrus.Debug("starting interrupt handler")
		s := <-c
		if s == nil {
			return
		}
		cmdCancel()
		logrus.Errorf("interrupt received: waiting 10 seconds before forcing exit")
		time.Sleep(10 * time.Second)
		logrus.Error("forcing exit")
		os.Exit(1)
	}()

	logrus.Debug("validating flags")
	err := validateFlags()
	if err != nil {
		logrus.Errorf("error validating flags: %v", err)
		os.Exit(1)
	}

	logrus.Debug("executing flags")
	execFlags()

	if !options.Silent {
		fmt.Print(banner)
	}

	var inputScanner *bufio.Scanner

	if options.Stdin {
		logrus.Info("reading from stdin")
		inputScanner = bufio.NewScanner(os.Stdin)
	} else {
		logrus.Infof("reading from file %q", options.File)
		file, err := os.Open(options.File)
		if err != nil {
			logrus.Errorf("error opening file: %v", err)
			os.Exit(1)
		}
		defer file.Close()
		inputScanner = bufio.NewScanner(file)
	}

	headers := transformHeaders()

	scanners := []reflections.ScanFunc{
		reflections.SimpleScan,
		reflections.ReplaceValuesHtmlCharsScan,
	}

	client := httpclient.NewHttpClient()
	client.AddLimiter(cmdCtx, options.RateLimit)
	client.AddHeaders(headers)

	if options.AllowInsecure {
		err := client.AllowInsecure()
		if err != nil {
			logrus.Errorf("error allowing insecure connections: %v", err)
			os.Exit(1)
		}
	}

	for inputScanner.Scan() {
		logrus.Debugf("scanning: %q", inputScanner.Text())
		input := inputScanner.Text()
		input, err = reflections.AddScheme(input)
		if err != nil {
			logrus.Errorf("error adding scheme to %q: %v", input, err)
			continue
		}
		for _, scanner := range scanners {
			if cmdCtx.Err() != nil {
				logrus.Debug("context cancelled")
				break
			}
			result, err := scanner(client, input, options.MinLength)
			if err != nil {
				logrus.Errorf("error scanning %q: %v", input, err)
				continue
			}
			for _, r := range result {
				msg := reflections.ReflectionToString(r, options.NoColor)
				fmt.Println(msg)
			}
		}
	}
}

func execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVarP(&options.AllowInsecure, "allow-insecure", "k", false, "Allow insecure connections")
	rootCmd.Flags().BoolVarP(&options.Debug, "debug", "d", false, "Enable debug mode")
	rootCmd.Flags().BoolVarP(&options.Silent, "silent", "s", false, "Outputs only errors and the results")
	rootCmd.Flags().StringVarP(&options.File, "file", "f", "", "File with URLs to scan")
	rootCmd.Flags().UintVarP(&options.MinLength, "min-length", "m", 3, "Minimum value's length to scan for reflections")
	rootCmd.Flags().UintVarP(&options.RateLimit, "rate-limit", "r", 25, "Number of requests per second")
	rootCmd.Flags().BoolVarP(&options.NoColor, "no-color", "n", false, "Disable color output")
	rootCmd.Flags().StringArrayVarP(&options.Headers, "header", "H", []string{}, `Headers to send with the request (specify multiple times)
Example: -H 'X-Forwarded-For: 127.0.0.1' -H 'X-Random: 1234'`)
}

func validateFlags() error {
	if options.Debug && options.Silent {
		return fmt.Errorf("debug and silent flags cannot be used together")
	}

	if options.File != "" {
		options.Stdin = false
		_, err := os.Stat(options.File)
		if os.IsNotExist(err) {
			return fmt.Errorf("file %q does not exist", options.File)
		}
	}

	return nil
}

func execFlags() {
	if options.Debug {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debug("debug mode enabled")
	}

	if options.Silent {
		logrus.SetLevel(logrus.ErrorLevel)
		logrus.Debug("silent mode enabled")
	}
}

func transformHeaders() map[string][]string {
	headers := make(map[string][]string)
	for _, header := range options.Headers {
		parts := strings.SplitN(header, ":", 2)
		if len(parts) != 2 {
			logrus.Errorf("invalid header: %q", header)
			os.Exit(1)
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		_, ok := headers[key]
		if !ok {
			headers[key] = []string{}
		}
		headers[key] = append(headers[key], strings.TrimSpace(value))
	}
	return headers
}
