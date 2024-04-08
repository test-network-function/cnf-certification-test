// Copyright (C) 2022-2024 Red Hat, Inc.
//
// This program is free software; you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation; either version 2 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along
// with this program; if not, write to the Free Software Foundation, Inc.,
// 51 Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA.

package provider

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	defaultLog "log"
	"strings"

	"github.com/go-logr/logr"
	"github.com/go-logr/stdr"
	"github.com/test-network-function/cnf-certification-test/internal/log"
	corev1 "k8s.io/api/core/v1"

	"github.com/redhat-openshift-ecosystem/openshift-preflight/artifacts"
	plibRuntime "github.com/redhat-openshift-ecosystem/openshift-preflight/certification"
	plibContainer "github.com/redhat-openshift-ecosystem/openshift-preflight/container"
)

var (
	// Certain tests that have been known to fail because of injected containers (such as Istio) that fail certain tests.
	ignoredContainerNames = []string{"istio-proxy"}
)

// Tag and Digest should not be populated at the same time. Digest takes precedence if both are populated
type ContainerImageIdentifier struct {
	// Repository is the name of the image that you want to check if exists in the RedHat catalog
	Repository string `yaml:"repository" json:"repository"`

	// Registry is the name of the registry `docker.io` of the container
	// This is valid for container only and required field
	Registry string `yaml:"registry" json:"registry"`

	// Tag is the optional image tag. "latest" is implied if not specified
	Tag string `yaml:"tag" json:"tag"`

	// Digest is the image digest following the "@" in a URL, e.g. image@sha256:45b23dee08af5e43a7fea6c4cf9c25ccf269ee113168c19722f87876677c5cb2
	Digest string `yaml:"digest" json:"digest"`
}

type Container struct {
	*corev1.Container
	Status                   corev1.ContainerStatus
	Namespace                string
	Podname                  string
	NodeName                 string
	Runtime                  string
	UID                      string
	ContainerImageIdentifier ContainerImageIdentifier
	PreflightResults         PreflightResultsDB
}

func NewContainer() *Container {
	return &Container{
		Container: &corev1.Container{}, // initialize the corev1.Container object
	}
}

func (c *Container) GetUID() (string, error) {
	split := strings.Split(c.Status.ContainerID, "://")
	uid := ""
	if len(split) > 0 {
		uid = split[len(split)-1]
	}
	if uid == "" {
		log.Debug(fmt.Sprintf("could not find uid of %s/%s/%s\n", c.Namespace, c.Podname, c.Name))
		return "", errors.New("cannot determine container UID")
	}
	log.Debug(fmt.Sprintf("uid of %s/%s/%s=%s\n", c.Namespace, c.Podname, c.Name, uid))
	return uid, nil
}

func (c *Container) SetPreflightResults(preflightImageCache map[string]PreflightResultsDB, env *TestEnvironment) error {
	log.Info("Running Preflight container test for container %q with image %q", c, c.Image)

	// Short circuit if the image already exists in the cache
	if _, exists := preflightImageCache[c.Image]; exists {
		log.Info("Container image %q exists in the cache. Skipping this run.", c.Image)
		c.PreflightResults = preflightImageCache[c.Image]
		return nil
	}

	opts := []plibContainer.Option{}
	opts = append(opts, plibContainer.WithDockerConfigJSONFromFile(env.GetDockerConfigFile()))
	if env.IsPreflightInsecureAllowed() {
		log.Info("Insecure connections are being allowed to Preflight")
		opts = append(opts, plibContainer.WithInsecureConnection())
	}

	// Create artifacts handler
	artifactsWriter, err := artifacts.NewMapWriter()
	if err != nil {
		return err
	}
	ctx := artifacts.ContextWithWriter(context.TODO(), artifactsWriter)

	// Add logger output to the context
	logbytes := bytes.NewBuffer([]byte{})
	checklogger := defaultLog.Default()
	checklogger.SetOutput(logbytes)
	logger := stdr.New(checklogger)
	ctx = logr.NewContext(ctx, logger)

	check := plibContainer.NewCheck(c.Image, opts...)

	results, runtimeErr := check.Run(ctx)
	if runtimeErr != nil {
		_, checks, err := check.List(ctx)
		if err != nil {
			return fmt.Errorf("could not get preflight container test list")
		}

		results.TestedImage = c.Image
		for _, c := range checks {
			results.PassedOverall = false
			result := plibRuntime.Result{Check: c, ElapsedTime: 0}
			results.Errors = append(results.Errors, *result.WithError(runtimeErr))
		}
	}

	// Take all of the Preflight logs and stick them into our log.
	log.Info(logbytes.String())

	// Store the Preflight test results into the container's PreflightResults var and into the cache.
	resultsDB := GetPreflightResultsDB(&results)
	c.PreflightResults = resultsDB
	preflightImageCache[c.Image] = resultsDB
	return nil
}

func (c *Container) StringLong() string {
	return fmt.Sprintf("node: %s ns: %s podName: %s containerName: %s containerUID: %s containerRuntime: %s",
		c.NodeName,
		c.Namespace,
		c.Podname,
		c.Name,
		c.Status.ContainerID,
		c.Runtime,
	)
}
func (c *Container) String() string {
	return fmt.Sprintf("container: %s pod: %s ns: %s",
		c.Name,
		c.Podname,
		c.Namespace,
	)
}

func (c *Container) HasIgnoredContainerName() bool {
	for _, ign := range ignoredContainerNames {
		if c.IsIstioProxy() || strings.Contains(c.Name, ign) {
			return true
		}
	}
	return false
}

func (c *Container) IsIstioProxy() bool {
	return c.Name == "istio-proxy" //nolint:goconst
}

func (c *Container) HasExecProbes() bool {
	return c.LivenessProbe != nil && c.LivenessProbe.Exec != nil ||
		c.ReadinessProbe != nil && c.ReadinessProbe.Exec != nil ||
		c.StartupProbe != nil && c.StartupProbe.Exec != nil
}

func (c *Container) IsTagEmpty() bool {
	return c.ContainerImageIdentifier.Tag == ""
}

func (c *Container) IsReadOnlyRootFilesystem(logger *log.Logger) bool {
	logger.Info("Testing Container %q", c)
	if c.Container.SecurityContext == nil || c.Container.SecurityContext.ReadOnlyRootFilesystem == nil {
		return false
	}
	return *c.Container.SecurityContext.ReadOnlyRootFilesystem
}

func (c *Container) IsContainerRunAsNonRoot() bool {
	if c.SecurityContext != nil && c.SecurityContext.RunAsNonRoot != nil {
		return *c.SecurityContext.RunAsNonRoot
	}
	return false
}
