/*
Copyright 2019 The Kubernetes Authors.

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

package nodeimage

import (
	"io/ioutil"

	"sigs.k8s.io/kind/pkg/exec"
)

type containerdImporter struct {
	containerCmder exec.Cmder
}

func newContainerdImporter(containerCmder exec.Cmder) *containerdImporter {
	return &containerdImporter{
		containerCmder: containerCmder,
	}
}

func (c *containerdImporter) Prepare() error {
	if err := c.containerCmder.Command(
		"bash", "-c", "nohup crio > /dev/null 2>&1 &",
	).Run(); err != nil {
		return err
	}
	// TODO(bentheelder): some healthcheck?
	return nil
}

func (c *containerdImporter) End() error {
	return c.containerCmder.Command("pkill", "crio").Run()
}

func (c *containerdImporter) Pull(image, platform string) error {
	// TODO: this should exist with a --no-unpack and some way to operate quietly
	// without discarding output
	return c.containerCmder.Command(
		"podman", "image", "pull", image, "--arch", platform,
	).SetStdout(ioutil.Discard).SetStderr(ioutil.Discard).Run()
}

func (c *containerdImporter) LoadCommand() exec.Cmd {
	return c.containerCmder.Command(
		// TODO: ideally we do not need this in the future. we have fixed at least one image
		"podman", "import", "-",
	)
}

func (c *containerdImporter) ListImported() ([]string, error) {
	return exec.OutputLines(c.containerCmder.Command("podman", "images", "-a", "-q"))
}
