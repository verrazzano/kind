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

ARG BASE_IMAGE=ghcr.io/oracle/oraclelinux:8-slim
FROM $BASE_IMAGE as build

# copy in static files
# all scripts are 0755 (rwx r-x r-x)
COPY --chmod=0755 files/usr/local/bin/* /usr/local/bin/

# all configs are 0644 (rw- r-- r--)
COPY files/etc/* /etc/
COPY files/etc/containers/* /etc/containers/
COPY files/etc/crio/crio.conf.d/* /etc/crio/crio.conf.d/
COPY files/etc/default/* /etc/default/
COPY files/etc/sysctl.d/* /etc/sysctl.d/
COPY files/etc/systemd/system/* /etc/systemd/system/
COPY files/etc/systemd/system/kubelet.service.d/* /etc/systemd/system/kubelet.service.d/

RUN chmod 0644 /etc/containers/ /etc/crio/crio.conf.d/ /etc/default/ /etc/sysctl.d/ /etc/systemd/system/ \
               /etc/systemd/system/kubelet.service.d/
RUN  microdnf update -y \
     && microdnf install -y  yum  yum-utils \
     && microdnf clean all \
     && rm -rf /var/cache/yum/*


RUN echo "Adding yum repos ..." \
    && yum-config-manager --add-repo https://yum.oracle.com/repo/OracleLinux/OL8/olcne15/x86_64/ \
    && yum-config-manager --add-repo https://yum.oracle.com/repo/OracleLinux/OL8/developer/x86_64/ \
    && yum-config-manager --add-repo https://yum.oracle.com/repo/OracleLinux/OL8/addons/x86_64/ \
    && yum-config-manager --add-repo https://yum.oracle.com/repo/OracleLinux/OL8/UEKR6/x86_64/ \
    && yum-config-manager --add-repo http://yum.oracle.com/repo/OracleLinux/OL8/kvm/appstream/x86_64

RUN echo "Installing Packages ..." \
    && dnf -y install hostname which vi systemd iptables fuse-overlayfs \
    && dnf clean all \
    && rm -rf /var/cache/yum/*


RUN echo "Installing OCNE binaries ..." \
    && dnf install -y  https://yum.oracle.com/repo/OracleLinux/OL8/olcne15/x86_64/getPackage/conmon-2.0.26-2.el8.x86_64.rpm  \
    && dnf install -y  cri-tools cri-o olcne-utils  \
    && dnf install -y kubeadm-1.23.7-1.el8 kubectl-1.23.7-1.el8 kubelet-1.23.7-1.el8 kubernetes-cni-1.0.1-2.el8 kubernetes-cni-plugins-1.0.1-3.el8 \
    && systemctl enable crio.service \
    && dnf clean all \
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


RUN echo "Enabling kubelet ... " \
    && systemctl enable kubelet.service


RUN echo "Ensuring /etc/kubernetes/manifests" \
    && mkdir -p /etc/kubernetes/manifests
#
RUN echo "Adjusting systemd-tmpfiles timer" \
    && sed -i /usr/lib/systemd/system/systemd-tmpfiles-clean.timer -e 's#OnBootSec=.*#OnBootSec=1min#'

RUN mkdir /license
COPY LICENSE README.md THIRD_PARTY_LICENSES.txt /license/

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