# Overview

This file will act as a decision log

## Wrapping in Minio Client

There is some wrapping in the Minio Client in order to both apply unit and integratioon tests.

## PUT behaviour

Put behaviour will follow the rules of [RFC 7231](https://datatracker.ietf.org/doc/html/rfc7231#section-4.3.4), therefore it will create or replace the resource.

## Consistent hashing

I've chosen the simplest consistent hashing algorithm just to showcase what should be done in terms of load balancing, I would use k8s to leverage it for me if this was an application going to production.

## Git log

Git history is not the cleanest, I can clean the commits but it will take me some time so I've lefted that as it is.