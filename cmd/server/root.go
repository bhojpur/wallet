package cmd

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import (
	"fmt"
	"log"
	"os"

	"github.com/bhojpur/wallet/pkg/config"
	"github.com/bhojpur/wallet/pkg/registry"
	"github.com/bhojpur/wallet/pkg/routing"
	"github.com/bhojpur/wallet/pkg/storage/postgres"

	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var verbose bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "wallet",
	Short: "Bhojpur Wallet is a digital wallet processing engine powered by Kubernetes",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if verbose {
			logger.SetLevel(logger.DebugLevel)
			logger.Debug("verbose logging enabled")
		}
	},

	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	// read yaml config file. Dont pass path to read
	// from default path
	yamlConfig := config.ReadYaml("")
	config := config.GetConfig(*yamlConfig)

	database, err := postgres.NewDatabase(config)
	if err != nil {
		log.Printf("database err %s", err)
		os.Exit(1)
	}

	// run migrations; update tables
	postgres.Migrate(database)

	channels := registry.NewChannels()
	domain := registry.NewDomain(config, database, channels)

	// create the fiber server.
	server := routing.Router(domain, config) // add endpoints

	// listen and serve
	port := fmt.Sprintf(":%v", 6700)
	log.Fatal(server.Listen(port))
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "en/disable verbose logging")
}
