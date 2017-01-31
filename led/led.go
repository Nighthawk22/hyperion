// Copyright 2017 Hyperion Team
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not
// use this file except in compliance with the License.  You may obtain a copy
// of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.  See the
// License for the specific language governing permissions and limitations under
// the License.
package led

import "os/exec"

const ledCommand = "/usr/local/bin/pigs"

// ChangeLED changes the given port to the given Brigthness. Will be run with pigs.
func ChangeLED(portNumber string, brightness string) error {
	err := exec.Command(ledCommand, "p", portNumber, brightness).Run()
	if err != nil {
		return err
	}
	return nil
}
