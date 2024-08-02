// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/helloworlddan/run" // <--- Loads of useful stuff for running on Cloud Run
)

func main() {
	// Say hi when the server starts listening
	run.Noticef(nil, "Hi 🦫 Let's start the service '%s' in project '%s'", run.Name(), run.ProjectID())

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Let's log in some different severities
		run.Default(r, "I got called!")
		run.Debug(r, "I got called!")
		run.Info(r, "I got called!")
		run.Notice(r, "I got called!")
		run.Warning(r, "I got called!")
		run.Critical(r, "I got called!")
		run.Alert(r, "I got called!")
		run.Emergency(r, "I got called!")

		// Drop some useful info when debugging usxing a type formatter
		run.Debugf(nil, "I am running in region '%s'", run.Region())

		fmt.Fprintln(w, "What's up log? 🪵")
	})

	err := run.ServeHTTP(func(_ context.Context) {
		// Catch SIGTERM and say goodbye after shutdown
		run.Notice(nil, "Goodbye 👋")
	}, nil)
	if err != nil {
		run.Error(nil, err)
	}
}
