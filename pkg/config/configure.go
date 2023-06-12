/*
Package config

Copyright © 2023 Gabriel Duke <gabeduke@gmail.com>

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
package config

import (
	"fmt"
	"github.com/gabeduke/wio-cli-go/pkg/util"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net"
	"net/url"
)

func configure(logger *log.Entry, cfgFile string) error {
	logger.Debug("configure called")

	// Prompt for user email address
	viper.Set(EMAIL, util.Prompt("Enter your email address: ", viper.GetString(EMAIL)))

	// Prompt for server address
	viper.Set(HOST, util.Prompt("Enter the server address (eg. https://wio.leetserve.com): ", viper.GetString(HOST)))

	// Prompt for server IP
	mip := util.Prompt("Enter the server IP address (leave blank to allow discovery): ", "")
	if mip == "" {
		host, err := url.Parse(viper.GetString(HOST))
		if err != nil {
			return errors.Errorf("Error parsing server address: %v", err)
		}
		hostAddr, err := net.LookupIP(host.Hostname())
		if err != nil {
			fmt.Println("Unknown host")
		} else {
			mip = hostAddr[0].String()
		}
	}
	viper.Set(HOST_IP, mip)

	logger.Debugf("Wio CLI Configuration: %v", viper.AllSettings())
	logger.WithField("file", viper.ConfigFileUsed()).Info("Wio CLI Configuration file updated")

	return viper.WriteConfigAs(viper.ConfigFileUsed())
}
