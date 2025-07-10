# Quick Start Guide

## Introduction

The Terraform provider for SAP Cloud Connector enables you to automate the provisioning, management, and configuration of resources on [SAP Cloud Connector](https://help.sap.com/docs/connectivity/sap-btp-connectivity-cf/cloud-connector). By leveraging this provider, you can simplify and streamline the deployment and maintenance of [subaccount](https://help.sap.com/docs/connectivity/sap-btp-connectivity-cf/managing-subaccounts), [system mapping](https://help.sap.com/docs/connectivity/sap-btp-connectivity-cf/configure-access-control-http#loioe7d4927dbb571014af7ef6ebd6cc3511__expose), [system mapping resources](https://help.sap.com/docs/connectivity/sap-btp-connectivity-cf/configure-access-control-http#loioe7d4927dbb571014af7ef6ebd6cc3511__fly), [domain mappings](https://help.sap.com/docs/connectivity/sap-btp-connectivity-cf/configure-domain-mappings-for-cookies) and [service channels](https://help.sap.com/docs/connectivity/sap-btp-connectivity-cf/using-service-channels) over a cloud connector instance.

## Prerequisites

To follow along with this tutorial, ensure you have already installed and configrued a [SAP Cloud Connector Instance](https://help.sap.com/docs/connectivity/sap-btp-connectivity-cf/cloud-connector-initial-configuration) and Terraform installed on your machine. You can download it from the official [Terraform website](https://developer.hashicorp.com/terraform/downloads).

## Authentication

In order to run the scripts, you need the credentials of an admin on the instance. Terraform Provider for SAP Cloud Connector supports the following authentication methods:

1. [Basic Authentication](./basic_auth.md) 
2. [X.509 Certificate Authentication](cert_auth.md)

Refer to the link corresponding to the chosen authentication method.

## Documentation

Terraform Provider for SAP Cloud Connector [Documentation](https://registry.terraform.io/providers/SAP/scc/latest/docs)
