# Copyright 2022 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

FROM ghcr.io/oracle/oraclelinux:8-slim
ARG VERSION
RUN microdnf -y upgrade && \
    microdnf -y module enable go-toolset:ol8 && \
    microdnf -y install go-toolset && \
    microdnf install git && \
    git clone https://github.com/rancher/local-path-provisioner && \
    cd local-path-provisioner && \
    git fetch && git checkout "${VERSION}" && \
    scripts/build && \
    mv bin/local-path-provisioner /usr/local/bin/local-path-provisioner \
    && GOBIN=/usr/local/bin go install github.com/google/go-licenses@latest \
    && GOARCH=$TARGETARCH go-licenses save --save_path=/_LICENSES .

FROM ghcr.io/oracle/oraclelinux:8-slim
COPY --from=0 /usr/local/bin/local-path-provisioner /usr/local/bin/local-path-provisioner
COPY --from=0 /_LICENSES/* /LICENSES/
COPY files/LICENSES/* /LICENSES/*
RUN chmod 0644 /LICENSES/*
ENTRYPOINT /usr/local/bin/local-path-provisioner
