# Overview

This file will act as a decision log

## Wrapping in Minio Client

There is some wrapping in the Minio Client in order to both apply unit and integratioon tests.

## PUT behaviour

Put behaviour will follow the rules of [RFC 7231](https://datatracker.ietf.org/doc/html/rfc7231#section-4.3.4), therefore it will create or replace the resource.
