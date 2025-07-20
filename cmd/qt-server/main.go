/*
 * Copyright © 2025 UMU618 <umu618@hotmail.com>
 *
 * This file is part of UMU618/qt, licensed under GPLv3. See LICENSE for details.
 */
package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/UMU618/qt/internal/log"
	"github.com/UMU618/qt/internal/options"
	"github.com/UMU618/qt/internal/server"
	"github.com/UMU618/qt/internal/token"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	serOptions *options.ServerOptions
	secOptions *options.SecureOptions
	logOptions *log.Options
)

func buildCommand(basename string) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   basename,
		Short: "Start up the server side endpoint",
		Long: `Establish a fast & security tunnel,
make you can access remote TCP/UNIX application like local application.
	   
Find more qt information at:
	https://github.com/UMU618/qt/blob/master/README.md`,
		RunE: runCommand,
	}
	// Initialize the flags needed to start the server
	rootCmd.Flags().AddGoFlagSet(flag.CommandLine)
	serOptions.AddFlags(rootCmd.Flags())
	secOptions.AddFlags(rootCmd.Flags())
	options.AddConfigFlag(basename, rootCmd.Flags())
	logOptions.AddFlags(rootCmd.Flags())

	return rootCmd
}

func runCommand(cmd *cobra.Command, args []string) error {
	options.PrintWorkingDir()
	options.PrintFlags(cmd.Flags())
	options.PrintConfig()

	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return err
	}

	if err := viper.Unmarshal(serOptions); err != nil {
		return err
	}

	if err := viper.Unmarshal(secOptions); err != nil {
		return err
	}

	if err := viper.Unmarshal(logOptions); err != nil {
		return err
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// run server
	runFunc(ctx, serOptions, secOptions)
	return nil
}

func runFunc(ctx context.Context, so *options.ServerOptions, seco *options.SecureOptions) {
	log.Init(logOptions)
	defer log.Flush()

	keyFile := seco.KeyFile
	certFile := seco.CertFile
	tokenParserPlugin := so.TokenParserPlugin
	tokenParserKey := so.TokenParserKey

	var tlsConfig *tls.Config
	if keyFile == "" || certFile == "" {
		panic("should specify cert file and key file!")
	} else {
		tlsCert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			log.Errorw("Certificate file or private key file is invalid.", "error", err.Error())
			return
		}
		tlsConfig = &tls.Config{
			Certificates: []tls.Certificate{tlsCert},
			NextProtos:   []string{"qt"},
		}

	}

	// Start server
	s := &server.Server{
		BindEndpoint: so.BindEndpoint,
		TlsConfig:    tlsConfig,
		TokenParser:  loadTokenParserPlugin(tokenParserPlugin, tokenParserKey),
	}
	s.Start(ctx)
}

func loadTokenParserPlugin(plugin string, key string) token.TokenParserPlugin {
	switch strings.ToLower(plugin) {
	case "cleartext":
		return token.NewCleartextTokenParserPlugin(key)
	default:
		panic(fmt.Sprintf("Token parser plugin %s don't support", plugin))
	}
}

func main() {
	// Initialize the options needed to start the server
	serOptions = options.GetDefaultServerOptions()
	secOptions = options.GetDefaultSecureOptions()
	logOptions = log.NewOptions()

	rootCmd := buildCommand("qt-server")
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
