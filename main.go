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
	"io"
	"net/http"

	"cloud.google.com/go/storage"

	exporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	otelhttp "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	otel "go.opentelemetry.io/otel"
	attribute "go.opentelemetry.io/otel/attribute"
	propagation "go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"github.com/helloworlddan/run" // <--- Loads of useful stuff for running on Cloud Run
)

func main() {
	// Say hi when the server starts booting
	run.Noticef(
		nil,
		"Hi 🦫 Let's start the service '%s' in project '%s'",
		run.Name(),
		run.ProjectID(),
	)

	ctx := context.Background()

	// Instantiate clients early and eagerly
	httpClient := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
	gsClient, err := storage.NewClient(ctx)
	if err != nil {
		run.Fatal(nil, err)
	}

	// Configure otel for Cloud Trace
	exporter, err := exporter.New(exporter.WithProjectID(run.ProjectID()))
	if err != nil {
		run.Fatal(nil, err)
	}
	provider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithSampler(sdktrace.TraceIDRatioBased(0.1)),
	)
	otel.SetTracerProvider(provider)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)
	tracer := provider.Tracer(run.Name())

	// Logging severities
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

		// Drop some useful info when debugging using a type formatter
		run.Debugf(nil, "I am running in region '%s'", run.Region())

		// Print all incoming headers
		for name, values := range r.Header {
			for _, value := range values {
				run.Debugf(r, "[HEADER] '%s': '%s'", name, value)
			}
		}

		fmt.Fprintln(w, "What's up log? 🪵")
	})

	// Propagation among services
	// NOTE: might need otelhttp
	http.HandleFunc("/service-to-service", func(w http.ResponseWriter, r *http.Request) {
		run.Debug(r, "I am calling someone else!")

		ctx := r.Context()
		remoteService := "https://whats-up-log-549074658641.europe-north2.run.app"
		req, err := http.NewRequestWithContext(ctx, "GET", remoteService, nil)
		if err != nil {
			run.Error(r, err)
			return
		}

		resp, err := httpClient.Do(req)
		if err != nil {
			run.Error(r, err)
			return
		}
		defer resp.Body.Close()

		fmt.Fprintln(w, resp.Status)
	})

	// Propagation to *.googleapis.com
	http.HandleFunc("/google-service", func(w http.ResponseWriter, r *http.Request) {
		run.Debug(r, "I am calling Cloud Storage!")

		ctx := r.Context()
		reader, err := gsClient.Bucket("gcs-whats-up-log").Object("pikachu.jpg").NewReader(ctx)
		if err != nil {
			run.Error(r, err)
			return
		}
		defer reader.Close()

		_, err = io.Copy(w, reader)
		if err != nil {
			run.Error(r, err)
			return
		}
	})

	// With OTEL instrumentation
	http.HandleFunc("/otel-instrumentation", func(w http.ResponseWriter, r *http.Request) {
		run.Debug(r, "I am doing some work!")

		ctx := r.Context()
		ctx, span := tracer.Start(ctx, "/otel-instrumentation")
		span.SetAttributes(
			attribute.KeyValue{
				Key:   attribute.Key("service.name"),
				Value: attribute.StringValue("brite-demo"),
			},
		)
		defer span.End()

		// Do something with the new child context
		reader, err := gsClient.Bucket("gcs-whats-up-log").Object("pikachu.png").NewReader(ctx)
		if err != nil {
			run.Error(r, err)
			return
		}
		defer reader.Close()

		fmt.Fprintln(w, "work completed")
	})

	err = run.ServeHTTP(func(_ context.Context) {
		// Catch SIGTERM and say goodbye after shutdown
		gsClient.Close()
		run.Notice(nil, "Goodbye 👋")
	}, nil)
	if err != nil {
		run.Error(nil, err)
	}
}
