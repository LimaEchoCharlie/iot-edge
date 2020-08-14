/*
 * Copyright 2020 ForgeRock AS
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"flag"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"time"
)

func main() {
	container := flag.String("container", "AM container name", "The name of the AM container")
	timeout := flag.Int("timeout", 600, "Maximum number of seconds to wait")
	flag.Parse()

	fmt.Println("\nWaiting for AM to start within container", *container)

	timer := time.NewTimer(time.Duration(*timeout) * time.Second)
	defer timer.Stop()

	var output []byte
	var err error
	var started bool

	for !started {
		select {
		case <-timer.C:
			log.Fatalf("\nTimeout fired.\nLast logs\n\n%s", string(output))
		default:
			time.Sleep(500 * time.Millisecond)
			cmd := exec.Command("docker", "logs", *container, "--tail", "25")
			output, err = cmd.CombinedOutput()

			if err != nil {
				log.Fatalf("\nLog read error, %s\nOutput\n\n%s", err, string(output))
			}

			started, err = regexp.Match("Server startup", output)
			if err != nil {
				log.Fatalf("\nRegex error; %s", err)
			}
			fmt.Print(".")
		}
	}
	fmt.Println("\nAM Started")
}
