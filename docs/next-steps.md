# Overview

This document will document the software next steps.

## API version

This product should have a versioned API in order to have more flexibility when providing the API to clients.

## CI/CD

The application should have a pipeline running tests and deploying frequently to ensure it's health, it also should be assessed if nightly runs for tests would prove valuable, I would say no if this is deployed frequently, but yes otherwise.

## Monitoring

There should exist some monitoring such as alerts and dashboards to check for the application's health. Examples of these alerts could be like 1 5xx error, a dashboard with the activity of the system would also be valuable to assess the architecture, should this be a microservice or a lambda.

## Testing

Testing should be improved with testing on the load balancer not done for time sake.

## Architecture

A diagram with the architecture of the application would be very valuable to onboard new people, as it stands and no interactions exist it is no needed.

## Logging

There should be added logs to whole application.