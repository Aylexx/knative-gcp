/*
Copyright 2020 Google LLC

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

package ingress

import (
	"context"

	"cloud.google.com/go/pubsub"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/knative-gcp/pkg/broker/config"
	"knative.dev/eventing/pkg/kncloudevents"
)

// HandlerOption is the option to configure ingress handler.
type HandlerOption func(*handler) error

// WithPort specifies the port number that ingress listens on. It will create an HttpMessageReceiver for that port.
func WithPort(port int) HandlerOption {
	return func(h *handler) error {
		h.httpReceiver = kncloudevents.NewHttpMessageReceiver(port)
		return nil
	}
}

// WithProjectID creates a pubsub client for the given project ID to communicate with pubsusb decouple topics.
func WithProjectID(id string) HandlerOption {
	return func(h *handler) error {
		ctx := context.Background()
		p, err := pubsub.NewClient(ctx, id)
		if err != nil {
			return err
		}
		h.decouple, err = NewMultiTopicDecoupleSink(context.Background(), WithPubsubClient(p))
		return err
	}
}

// MultiTopicDecoupleSinkOption is the option to configure multiTopicDecoupleSink.
type MultiTopicDecoupleSinkOption func(sink *multiTopicDecoupleSinkOptions)

type multiTopicDecoupleSinkOptions struct {
	client       cloudevents.Client
	pubsub       *pubsub.Client
	brokerConfig config.ReadonlyTargets
}

// WithClient specifies an eventing client to use.
func WithClient(client cloudevents.Client) MultiTopicDecoupleSinkOption {
	return func(opts *multiTopicDecoupleSinkOptions) {
		opts.client = client
	}
}

// WithPubsubClient specifies the pubsub client to use.
func WithPubsubClient(ps *pubsub.Client) MultiTopicDecoupleSinkOption {
	return func(opts *multiTopicDecoupleSinkOptions) {
		opts.pubsub = ps
	}
}

// WithBrokerConfig specifies the broker config. It can be created by reading a configmap mount.
func WithBrokerConfig(brokerConfig config.ReadonlyTargets) MultiTopicDecoupleSinkOption {
	return func(opts *multiTopicDecoupleSinkOptions) {
		opts.brokerConfig = brokerConfig
	}
}
