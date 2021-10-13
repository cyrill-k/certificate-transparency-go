// Copyright 2021 Google LLC. All Rights Reserved.
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

// Package witness is designed to make sure the STHs of CT logs are consistent
// and store/serve/sign them if so.  It is expected that a separate feeder
// component would be responsible for the actual interaction with logs.
package main

import (
	"context"
	"flag"
	"io/ioutil"

	"github.com/golang/glog"
	"github.com/google/certificate-transparency-go/witness/witness/cmd/witness/impl"
	"gopkg.in/yaml.v2"
)

var (
	listenAddr = flag.String("listen", ":8000", "address:port to listen for requests on")
	dbFile     = flag.String("db_file", ":memory:", "path to a file to be used as sqlite3 storage for STHs, e.g. /tmp/chkpts.db")
	configFile = flag.String("config_file", "example_config.yaml", "path to a YAML config file that specifies the logs followed by this witness")
	witnessSK  = flag.String("private_key", "", "private signing key for the witness")
)

func main() {
	flag.Parse()

	if *witnessSK == "" {
		glog.Exit("--private_key must not be empty")
	}

	if *configFile == "" {
		glog.Exit("--config_file must not be empty")
	}
	fileData, err := ioutil.ReadFile(*configFile)
	if err != nil {
		glog.Exitf("Failed to read from config file: %v", err)
	}
	var lc impl.LogConfig
	if err := yaml.Unmarshal(fileData, &lc); err != nil {
		glog.Exitf("Failed to parse config file as proper YAML: %v", err)
	}

	ctx := context.Background()
	if err := impl.Main(ctx, impl.ServerOpts{
		ListenAddr: *listenAddr,
		DBFile:     *dbFile,
		PrivKey:    *witnessSK,
		Config:     lc,
	}); err != nil {
		glog.Exitf("Error running witness: %v", err)
	}
}