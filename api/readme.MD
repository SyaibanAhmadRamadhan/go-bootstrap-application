# API Directory

This directory contains the OpenAPI specification of the API both for
our public-facing API and the external API that we consume.

## 1. Directory Structure

The directory structure of the API directory is as follows:

```
/
├── proto This directory contains the protocol buffer definitions for interoperability with other internal services.
├── provider This directory contains the OpenAPI specification of the API that we consume, if they don't have any Go SDK.
├── api.yml This is our main public-facing API specification.
```