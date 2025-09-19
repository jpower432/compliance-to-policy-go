/*
Copyright 2023 IBM Corporation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	hplugin "github.com/hashicorp/go-plugin"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	olog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.34.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/oscal-compass/compliance-to-policy-go/v2/cmd/kyverno-plugin/server"
	"github.com/oscal-compass/compliance-to-policy-go/v2/plugin"
)

var (
	serviceName = semconv.ServiceNameKey.String("kyverno")
	shutdown    func(ctx context.Context) error
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Failed to execute command: %v", err)
	}
}

var rootCmd = &cobra.Command{
	Use:   "kyverno",
	Short: "A brief description of your Kyverno JSON plugin.",
	Long:  `A longer description of your Kyverno JSON plugin, which can span multiple lines.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		endpoint, _ := cmd.Flags().GetString("endpoint")
		skipTLS, _ := cmd.Flags().GetBool("skip-tls")
		skipTLSVerify, _ := cmd.Flags().GetBool("skip-tls-verify")

		conn, err := newClient(endpoint, skipTLS, skipTLSVerify)
		if err != nil {
			return fmt.Errorf("failed to create gRPC connection to collector: %w", err)
		}
		shutdown, err = otelSDKSetup(context.Background(), conn)
		if err != nil {
			return fmt.Errorf("failed to set up OpenTelemetry SDK: %v", err)
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		defer func() {
			if err := shutdown(context.Background()); err != nil {
				log.Fatalf("failed to shut down OpenTelemetry SDK: %v", err)
			}
		}()

		kyvernoPlugin, err := server.NewPlugin()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		plugins := map[string]hplugin.Plugin{
			plugin.PVPPluginName: &plugin.PVPPlugin{Impl: kyvernoPlugin},
		}
		config := plugin.ServeConfig{
			PluginSet: plugins,
			Logger:    server.Logger(),
		}
		plugin.Register(config)
	},
}

func init() {
	plugin.SetBase(rootCmd)
}

// otelSDKSetup completes setup of the Otel SDK with providers.
func otelSDKSetup(ctx context.Context, conn *grpc.ClientConn) (func(context.Context) error, error) {
	var shutdownFuncs []func(context.Context) error
	shutDown := func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}
		shutdownFuncs = nil
		return err
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			serviceName,
		),
	)
	if err != nil {
		return nil, err
	}

	// --- Start of Tracing Setup ---
	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, err
	}

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tracerProvider)
	shutdownFuncs = append(shutdownFuncs, tracerProvider.Shutdown)

	// And here, we set a global propagator. This is what handles injecting
	// context into gRPC metadata.
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	metricExporter, err := otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, err
	}

	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExporter, sdkmetric.WithInterval(3*time.Second))), sdkmetric.WithResource(res),
	)
	otel.SetMeterProvider(meterProvider)

	logExporter, err := otlploggrpc.New(ctx, otlploggrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, err
	}

	logProcessor := olog.NewSimpleProcessor(logExporter)
	logProvider := olog.NewLoggerProvider(olog.WithProcessor(logProcessor), olog.WithResource(res))

	// Register the provider as the global logger provider.
	global.SetLoggerProvider(logProvider)

	shutdownFuncs = append(shutdownFuncs, logProvider.Shutdown, meterProvider.Shutdown)

	return shutDown, nil
}

func newClient(otelEndpoint string, skipTLS, skipTLSVerify bool) (*grpc.ClientConn, error) {
	var creds credentials.TransportCredentials
	if skipTLS {
		creds = insecure.NewCredentials()
	} else {
		sysPool, err := x509.SystemCertPool()
		if err != nil {
			return nil, fmt.Errorf("failed to get system cert: %w", err)
		}
		// By default, skip TLS verify is false.
		creds = credentials.NewTLS(&tls.Config{RootCAs: sysPool, InsecureSkipVerify: skipTLSVerify}) /* #nosec G402  */ //pragma: allowlist secret
	}
	return grpc.NewClient(otelEndpoint, grpc.WithTransportCredentials(creds))
}
