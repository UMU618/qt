/*
 * Copyright © 2025 UMU618 <umu618@hotmail.com>
 *
 * This file is part of UMU618/qt, licensed under GPLv3. See LICENSE for details.
 */
package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/UMU618/qt/internal/client"
	"github.com/UMU618/qt/internal/log"
	"github.com/UMU618/qt/internal/options"
	"github.com/UMU618/qt/internal/token"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	clientOptions *options.ClientOptions
	secOptions    *options.ClientSecureOptions
	logOptions    *log.Options
)

func buildCommand(basename string) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   basename,
		Short: "Start up the client side endpoint",
		Long: `Establish a fast & security tunnel,
make you can access remote TCP/UNIX application like local application.
	   
Find more qt information at:
	https://github.com/UMU618/qt/blob/master/README.md`,
		RunE: runCommand,
	}
	// Initialize the flags needed to start the server
	rootCmd.Flags().AddGoFlagSet(flag.CommandLine)
	clientOptions.AddFlags(rootCmd.Flags())
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

	if err := viper.Unmarshal(clientOptions); err != nil {
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
	runFunc(ctx, clientOptions, secOptions)
	return nil
}

func runFunc(ctx context.Context, co *options.ClientOptions, seco *options.ClientSecureOptions) {
	log.Init(logOptions)
	defer log.Flush()

	localSocket := co.ListenOn
	serverEndpointSocket := co.ServerEndpointSocket
	tokenPlugin := co.TokenPlugin
	tokenSource := co.TokenSource
	serverCertFile := seco.ServerCertFile
	if serverCertFile == "" {
		panic("server-cert-file must be specified!")
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
			fmt.Println("UMU")
			return nil
		},
		NextProtos: []string{"qt"},
	}

	caPemBlock, err := os.ReadFile(serverCertFile)
	if err != nil {
		log.Errorw("Failed to read server cert file.", "error", err.Error())
	}
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(caPemBlock)
	tlsConfig.RootCAs = certPool

	// Start client endpoint
	c := client.ClientEndpoint{
		LocalSocket:          localSocket,
		ServerEndpointSocket: serverEndpointSocket,
		TokenSource:          loadTokenSourcePlugin(tokenPlugin, tokenSource),
		TlsConfig:            tlsConfig,
	}
	c.Start(ctx)
}

func loadTokenSourcePlugin(plugin string, source string) token.TokenSourcePlugin {
	switch strings.ToLower(plugin) {
	case "fixed":
		return token.NewFixedTokenPlugin(source)
	case "file":
		return token.NewFileTokenSourcePlugin(source)
	case "http":
		return token.NewHttpTokenPlugin(source)
	default:
		panic(fmt.Sprintf("The token source plugin %s is invalid", plugin))
	}
}

func main() {
	// Initialize the options needed to start the server
	clientOptions = options.GetDefaultClientOptions()
	secOptions = options.GetDefaultClientSecureOptions()
	logOptions = log.NewOptions()

	rootCmd := buildCommand("qt-client")
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
