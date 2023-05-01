# Copyright 2018 The Kubernetes Authors.
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

# kind node base image

ARG BASE_IMAGE=ghcr.io/oracle/oraclelinux:8

FROM $BASE_IMAGE as build
ARG PAUSE_IMAGE_TAG
ARG ETCD_IMAGE_TAG
ARG DNS_IMAGE_TAG
ARG K8S_VERSION

# copy in static files
# all scripts are 0755 (rwx r-x r-x)
COPY --chmod=0755 files/usr/local/bin/* /usr/local/bin/

# all configs are 0644 (rw- r-- r--)
COPY --chmod=0644 files/etc/* /etc/
COPY --chmod=0644 files/etc/crio/crio.conf /etc/crio/crio.conf
COPY --chmod=0644 files/etc/default/* /etc/default/
COPY --chmod=0644 files/etc/sysctl.d/* /etc/sysctl.d/
COPY --chmod=0644 files/etc/systemd/system/* /etc/systemd/system/
COPY --chmod=0644 files/etc/systemd/system/kubelet.service.d/* /etc/systemd/system/kubelet.service.d/


RUN  sed -i "s/container-registry.oracle.com\/olcne\/pause:pause_image_tag/container-registry.oracle.com\/olcne\/pause:${PAUSE_IMAGE_TAG}/g" /etc/crio/crio.conf \
     && sed -i "s/k8s_version/${K8S_VERSION}/g" /etc/kubeadm.yaml \
     && sed -i "s/etcd_image_tag/${ETCD_IMAGE_TAG}/g" /etc/kubeadm.yaml \
     && sed -i "s/dns_image_tag/${DNS_IMAGE_TAG}/g" /etc/kubeadm.yaml \
     && yum update -y \
     && yum clean all \
     && rm -rf /var/cache/yum/*

RUN echo "Adding yum repos ..." \
    && dnf install -y oracle-olcne-release-el8 \
    && dnf config-manager --enable ol8_olcne16 ol8_olcne15 ol8_addons ol8_baseos_latest ol8_appstream ol8_UEKR6 \
    && yum clean all \
    && rm -rf /var/cache/yum/*

RUN echo "Installing OCNE binaries ..." \
    && dnf install -y kubeadm-${K8S_VERSION}-1.el8 kubectl-${K8S_VERSION}-1.el8 kubelet-${K8S_VERSION}-1.el8 \
    && dnf install -y oraclelinux-developer-release-el8 python36-oci-cli olcnectl olcne-api-server olcne-utils \
    && yum -y install hostname which vi iptables fuse-overlayfs \
    && dnf -y reinstall systemd \
    && yum clean all \
    && rm -rf /var/cache/yum/*


RUN echo "Prep for systemd ..." \
    && (cd /lib/systemd/system/sysinit.target.wants/; for i in *; do [ $i == \
    systemd-tmpfiles-setup.service ] || rm -f $i; done); \
    rm -f /lib/systemd/system/multi-user.target.wants/*;\
    rm -f /etc/systemd/system/*.wants/*;\
    rm -f /lib/systemd/system/local-fs.target.wants/*; \
    rm -f /lib/systemd/system/sockets.target.wants/*udev*; \
    rm -f /lib/systemd/system/sockets.target.wants/*initctl*; \
    rm -f /lib/systemd/system/basic.target.wants/*;\
    rm -f /lib/systemd/system/anaconda.target.wants/*;


RUN echo "Enabling kubelet and crio... " \
    && systemctl enable crio.service \
    && systemctl enable kubelet.service


RUN echo "Ensuring /etc/kubernetes/manifests" \
    && mkdir -p /etc/kubernetes/manifests
#
RUN echo "Adjusting systemd-tmpfiles timer" \
    && sed -i /usr/lib/systemd/system/systemd-tmpfiles-clean.timer -e 's#OnBootSec=.*#OnBootSec=1min#'

RUN mkdir /license
#COPY LICENSE README.md THIRD_PARTY_LICENSES.txt /license/

# squash
FROM scratch
COPY --from=build / /
#
# tell systemd that it is in docker (it will check for the container env)
# https://systemd.io/CONTAINER_INTERFACE/
ENV container docker
# systemd exits on SIGRTMIN+3, not SIGTERM (which re-executes it)
# https://bugzilla.redhat.com/show_bug.cgi?id=1201657
STOPSIGNAL SIGRTMIN+3

# NOTE: this is *only* for documentation, the entrypoint is overridden later
ENTRYPOINT [ "/usr/local/bin/entrypoint", "/sbin/init" ]