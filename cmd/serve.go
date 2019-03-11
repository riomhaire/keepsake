// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dgrijalva/jwt-go"
	"github.com/ghodss/yaml"
	"github.com/riomhaire/keepsake/infrastructure/api"
	"github.com/riomhaire/keepsake/infrastructure/facades/storage"
	"github.com/riomhaire/keepsake/models"
	"github.com/riomhaire/keepsake/usecases"
	"github.com/spf13/cobra"
)



// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the OAuth2 Service on given port",
	Run: func(cmd *cobra.Command, args []string) {

		c := cmd.Flag("configuration").Value.String()
		log.Println("Starting Service Keepsake", VERSION, "... reading configuration", c)

		// read bootstap file
		data, err := ioutil.ReadFile(c)
		if err != nil {
			log.Fatal(err)
		}

		config := models.Configuration{}

		err = yaml.Unmarshal([]byte(data), &config)
		if err != nil {
			log.Fatal(err)
		}

		// Setup API
		tokenEncoderDecoder := usecases.NewTokenEncoderDecoder(jwt.SigningMethodHS256, config.MasterSecret, config.TimeToLiveSeconds)
		jwtEncoderDecoder := usecases.NewJWTEncoderDecoder(config.TimeToLiveSeconds)
		storageInteractor := storage.NewConfigurationStorageIntegrator(&config)

		rest := api.NewRestAPI(&config, tokenEncoderDecoder, jwtEncoderDecoder, storageInteractor)

		// Setup Shutdown
		p := make(chan os.Signal, 2)
		signal.Notify(p, os.Interrupt, syscall.SIGTERM)

		// Set up a way to cleanly shutdown / deregister
		go func() {
			<-p
			log.Println("Shutting Down")
			os.Exit(0)
		}()

		// OK start
		rest.Start()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	serveCmd.Flags().StringP("configuration", "c", "keepsake.yaml", "Where config file is located")
}
