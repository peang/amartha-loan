# Apps

This repository contains two applications:

1. **Setup App**: The setup app is responsible for setting up the necessary dependencies and configurations for the main application.

2. **Main App**: The main app is the main application that utilizes the setup app to run. It can be started with a default port of 8080 or a custom port of your choice.

## Setup App

To set up the dependencies and configurations for the main application, you can execute the following command:

```/bin/setup```

This script will perform the necessary setup steps, including downloading any required datasets, setting up the environment variables, and installing the dependencies.

## Main App

To start the main application, you can execute the following command:

```./bin/run```

This script will start the main application with the default port of 8080.

Alternatively, you can specify a custom port by providing the desired port number as an argument:

```./bin/run 8082```


This will start the main application with the specified port number.

## Usage

To use the main application, you can access it in your web browser by navigating to `http://localhost:{PORT}`, where `{PORT}` is the port number specified when starting the application.

