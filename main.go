// Copyright 2018 The Loopix-Messaging Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"github.com/nymtech/loopix-messaging/client"
	"github.com/nymtech/loopix-messaging/config"
	"github.com/nymtech/loopix-messaging/logging"
	"github.com/nymtech/loopix-messaging/pki"
	"github.com/nymtech/loopix-messaging/server"
	"github.com/nymtech/loopix-messaging/sphinx"

	"flag"
	"fmt"

	"github.com/nymtech/loopix-messaging/helpers"

	"github.com/golang/protobuf/proto"
)

var logLocal = logging.PackageLogger()

const (
	PKI_DIR = "pki/database.db"
)

func pkiPreSetting(pkiDir string) error {
	db, err := pki.OpenDatabase(pkiDir, "sqlite3")
	if err != nil {
		return err
	}
	defer db.Close()

	params := make(map[string]string)
	params["Id"] = "TEXT"
	params["Typ"] = "TEXT"
	params["Config"] = "BLOB"

	err = pki.CreateTable(db, "Pki", params)
	if err != nil {
		return err
	}

	return nil
}

//func FakeAdding(c *client.Client) {
//	logLocal.Info("Adding simulated traffic of a client")
//	for {
//		sphinxPacket, err := c.EncodeMessage("hello world", c.Config)
//		if err != nil {
//		}
//		packet, err := config.WrapWithFlag("\xc6", sphinxPacket)
//		if err != nil {
//			logLocal.Info("Something went wrong")
//		}
//		c.OutQueue <- packet
//		time.Sleep(10 * time.Second)
//	}
//}

func main() {

	typ := flag.String("typ", "", "A type of entity we want to run")
	id := flag.String("id", "", "Id of the entity we want to run")
	host := flag.String("host", "", "The host on which the entity is running")
	port := flag.String("port", "", "The port on which the entity is running")
	providerId := flag.String("provider", "", "The port on which the entity is running")
	flag.Parse()

	err := pkiPreSetting(PKI_DIR)
	if err != nil {
		panic(err)
	}

	ip, err := helpers.GetLocalIP()
	if err != nil {
		panic(err)
	}

	// even though we're just overwriting the passed value of host, the startup script(s) still rely on its existence
	// and hence can't, for time being, be removed
	host = &ip

	switch *typ {
	case "client":
		// DEPERCATED
		logLocal.Warn("Client startup using this entry point is deprecated. Please use daemon/client/client.go instead")
		db, err := pki.OpenDatabase(PKI_DIR, "sqlite3")

		if err != nil {
			panic(err)
		}

		row := db.QueryRow("SELECT Config FROM Pki WHERE Id = ? AND Typ = ?", providerId, "Provider")

		var results []byte
		err = row.Scan(&results)
		if err != nil {
			fmt.Println(err)
		}
		var providerInfo config.MixConfig
		if err := proto.Unmarshal(results, &providerInfo); err != nil {
			panic(err)
		}

		pubC, privC, err := sphinx.GenerateKeyPair()
		if err != nil {
			panic(err)
		}

		client, err := client.NewClient(*id, *host, *port, pubC, privC, PKI_DIR, providerInfo)
		if err != nil {
			panic(err)
		}

		err = client.Start()
		if err != nil {
			panic(err)
		}

	case "mix":
		pubM, privM, err := sphinx.GenerateKeyPair()
		if err != nil {
			panic(err)
		}

		mixServer, err := server.NewMixServer(*id, *host, *port, pubM, privM, PKI_DIR)
		if err != nil {
			panic(err)
		}

		err = mixServer.Start()
		if err != nil {
			panic(err)
		}
	case "provider":
		pubP, privP, err := sphinx.GenerateKeyPair()
		if err != nil {
			panic(err)
		}

		providerServer, err := server.NewProviderServer(*id, *host, *port, pubP, privP, PKI_DIR)
		if err != nil {
			panic(err)
		}

		err = providerServer.Start()
		if err != nil {
			panic(err)
		}
	}
}
