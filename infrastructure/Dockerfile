FROM python:3-slim-bullseye

WORKDIR /infrastructure/terraform

# install command.
RUN apt-get update && apt-get install -y less vim wget unzip jq nkf sudo

# install aws-cli
ARG AWSCLI_VERSION
RUN wget -nv https://awscli.amazonaws.com/awscli-exe-linux-x86_64-${AWSCLI_VERSION}.zip && \
  unzip ./awscli-exe-linux-x86_64-${AWSCLI_VERSION}.zip && ./aws/install && \
  rm -rf ./awscli-exe-linux-x86_64-${AWSCLI_VERSION}.zip ./aws

# install terraform
ARG TERRAFORM_VERSION
RUN wget -nv https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip && \
  unzip ./terraform_${TERRAFORM_VERSION}_linux_amd64.zip -d /usr/local/bin/ && \
  rm ./terraform_${TERRAFORM_VERSION}_linux_amd64.zip

# install tflint
ARG TFLINT_VERSION
RUN curl -s https://raw.githubusercontent.com/terraform-linters/tflint/master/install_linux.sh | bash

# COPY terraform related files
COPY . /infrastructure
