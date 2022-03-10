# Chainlink Integration Framework

[![Go Report Card](https://goreportcard.com/badge/github.com/smartcontractkit/integrations-framework)](https://goreportcard.com/report/github.com/smartcontractkit/integrations-framework)
[![Go Reference](https://pkg.go.dev/badge/github.com/smartcontractkit/integrations-framework.svg)](https://pkg.go.dev/github.com/smartcontractkit/integrations-framework)
![Tests](https://github.com/smartcontractkit/integrations-framework/actions/workflows/test.yaml/badge.svg)
![Lint](https://github.com/smartcontractkit/integrations-framework/actions/workflows/lint.yaml/badge.svg)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

The Chainlink Integration Framework is a blockchain development framework written in Go. Its primary purpose is to help chainlink developers create extensive integration, e2e, performance, and chaos tests to ensure the stability of the chainlink project. It can also be helpful to those who just want to use chainlink oracles in their projects to help test their contracts, or even for those that aren't using chainlink. This project makes ample use of the [helmenv](https://github.com/smartcontractkit/helmenv) package to launch resources and testing apparatus.

See the [docs](https://smartcontractkit.github.io/integrations-framework/) or our [go reference](https://pkg.go.dev/github.com/smartcontractkit/integrations-framework) page for more detailed info and examples. If you just want a quick overview, keep reading.

## Setup

In order to use this framework, you must have a connection to an actively running Kubernetes cluster. If you don't have one handy, check out [minikube](https://minikube.sigs.k8s.io/docs/start/) which should work fine for smaller tests, but will likely need to be allocated more power, or you'll need to use a more powerful cluster in general to run tests that require lots of services, like OCR. Others have had more luck with [K3D](https://k3d.io/v5.3.0/) when running more resource-intense tests.

## Usage

Here's a simple example on deploying and interacting with a basic storage contract using this framework and [Ginkgo](https://github.com/onsi/ginkgo), a BDD testing framework we've come to really enjoy. You can use another testing framework, including Go's default testing if you prefer otherwise.

See our [suite/smoke](suite/smoke) directory for quite a few examples of the framework's usage.

## Chainlink Values

If you would like to change the Chainlink values that are used for environments, you can use the `framework.yaml` file,
or set environment variables that are all caps versions of the values found in the config file.

```yaml
# Specify the image and version of the chainlink image you want to run tests against. Leave blank for default.

chainlink_image:      # Image of chainlink node
chainlink_version:    # Version of the image on the chainlink node
chainlink_env_values: # Environment values to pass onto the chainlink nodes

# Specify the image and version of the simulated geth image you want to run tests against. Leave blank for default.
# Has no effect when running tests on networks other than the simulated geth instances.
geth_image:   # Image of the simulated geth to use
geth_version: # Version of the geth image
geth_args:    # List of CLI arguments to pass to simulated geth image. WARNING

# Experimental: Let the remote test runner email you when the test is done.
remote_runner_email_server: smtp.gmail.com # Example if you use GMAIL
remote_runner_email_address: # Username or email address to use
remote_runner_email_pass: # Password to email account. REMEMBER TO DELETE
# Experimental: Let the remote test runner pass a message to a slack webhook when the test is done
# See https://slack.com/help/articles/115005265063-Incoming-webhooks-for-Slack
remote_slack_webhook: https://hooks.slack.com/services/XXX
```

### WARNING

Values passed into `geth_args` will fully REPLACE all existing defaults we use in our launch. This enables freedom from defaults, but you should most definitely look at all the [current defaults](https://github.com/smartcontractkit/helmenv/blob/master/charts/geth/values.yaml#L16) we usually use and replace them as necessary.

## Soak Tests

Soak testing support is in fairly preliminary support at the moment. Head over to the [soak test explainer](./SOAK_RUNNER.md) to learn more.
