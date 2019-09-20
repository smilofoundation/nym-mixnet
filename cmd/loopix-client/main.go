// Copyright 2019 The Loopix-Messaging Authors
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
	cmd "github.com/nymtech/loopix-messaging/cmd/loopix-client/commands"
	"github.com/tav/golly/optparse"
)

func main() {
	var logo = `
  _                      _      
 | |    ___   ___  _ __ (_)_  __
 | |   / _ \ / _ \| '_ \| \ \/ /
 | |___ (_) | (_) | |_) | |>  < 
 |_____\___/ \___/| .__/|_/_/\_\
		  |_|            (client)
		  
		  `
	cmds := map[string]func([]string, string){
		"run":  cmd.RunCmd,
		"init": cmd.InitCmd,
	}
	info := map[string]string{
		"run": "Run a persistent Loopix client process",
	}
	optparse.Commands("loopix-client", "0.0.2", cmds, info, logo)
}
