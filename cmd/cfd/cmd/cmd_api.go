/*
Copyright © 2020 @fernandezvara

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
	"context"
	"io/ioutil"
	"regexp"

	"github.com/fernandezvara/certsfor/db/store"
	"github.com/fernandezvara/certsfor/internal/api"
	"github.com/fernandezvara/certsfor/internal/manager"
	"github.com/fernandezvara/certsfor/internal/service"
	"github.com/fernandezvara/certsfor/pkg/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// apiCmd represents the bootstrap command
var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Starts as REST API service allowing remote cfd clients.",
	Long:  `Starts as REST API service allowing remote cfd clients.`,
	Run:   apiFunc,
}

func init() {
	startCmd.AddCommand(apiCmd)
}

func isUUID(uuid string) bool {
	return regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$").MatchString(uuid)
}

func fileBytes(filename string) ([]byte, error) {

	if filename == "" {
		return []byte{}, nil
	}

	return ioutil.ReadFile(filename)

}

func apiFunc(cmd *cobra.Command, args []string) {

	var (
		sto               store.Store
		srv               *service.Service
		a                 *api.API
		cert, key, cacert []byte
		err               error
	)

	sto, err = store.Open(context.Background(), viper.GetString(configDBType), viper.GetString(configDBConnectionString))
	er(err)
	srv = service.NewAsServer(sto, Version)

	if isUUID(viper.GetString(configTLSCA)) {
		// get cert from DB
		var (
			ca  *manager.CA
			crt client.Certificate
		)

		ca, err = srv.CAGet(viper.GetString(configTLSCA))
		er(err)

		crt, err = srv.CertificateGet(context.Background(), viper.GetString(configTLSCA), viper.GetString(configTLSCertificate), global.remaining)
		er(err)

		cacert = ca.CACertificateBytes()
		cert = crt.Certificate
		key = crt.Key

	} else {
		// ca certificate is a file?
		cacert, err = fileBytes(viper.GetString(configTLSCA))
		er(err)
		cert, err = fileBytes(viper.GetString(configTLSCertificate))
		er(err)
		key, err = fileBytes(viper.GetString(configTLSKey))
		er(err)

	}

	a = api.New(srv, Version)
	err = a.Start(viper.GetString(configAPIAddr), cert, key, cacert, viper.GetStringSlice(configAPIAccessLog), viper.GetStringSlice(configAPIErrorLog), viper.GetBool(configAPIDebugLog))
	er(err)

}
