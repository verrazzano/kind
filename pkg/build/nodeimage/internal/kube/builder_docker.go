/*
Copyright 2018 The Kubernetes Authors.

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

package kube

import (
	"os"
	"path/filepath"
	"sigs.k8s.io/kind/pkg/log"
)

// TODO(bentheelder): plumb through arch

// dockerBuilder implements Bits for a local docker-ized make / bash build
type dockerBuilder struct {
	kubeRoot string
	arch     string
	logger   log.Logger
}

var _ Builder = &dockerBuilder{}
const outPath = "/images/base/_output/"

// NewDockerBuilder returns a new Bits backed by the docker-ized build,
// given kubeRoot, the path to the kubernetes source directory
func NewDockerBuilder(logger log.Logger, kubeRoot, arch string) (Builder, error) {
	return &dockerBuilder{
		kubeRoot: kubeRoot,
		arch:     arch,
		logger:   logger,
	}, nil
}

func (b *dockerBuilder) OCNEBuild() (Bits, error) {
	cwd, err := os.Getwd()
	k8sVersion, ok := os.LookupEnv("K8S-VERSION")
	if !ok {
		k8sVersion = "v1.26.6"
	}
	if err != nil {
		return nil, err
	}
	return &bits{
		binaryPaths: []string{},
		imagePaths: []string{

			filepath.Join(cwd, outPath, "kube-apiserver.tar"),
			filepath.Join(cwd, outPath, "kube-controller-manager.tar"),
			filepath.Join(cwd, outPath, "kube-proxy.tar"),
			filepath.Join(cwd, outPath, "kube-scheduler.tar"),
			filepath.Join(cwd, outPath, "etcd.tar"),
			filepath.Join(cwd, outPath, "pause.tar"),
		},
		version: k8sVersion,
	}, nil
}
