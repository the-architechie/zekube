# Zekube

A minimal Docker container orchestrator written in Go.

## Description

Zekube is a simplified container orchestrator designed for learning purposes. While Kubernetes is powerful, it often feels like a black box. This project is an implementation attempt to understand the inner workings of a typical orchestrator, focusing on the fundamental components like scheduling, worker nodes, and task management.

## Project Structure

- **Manager**: The central control plane that manages the state of the cluster, keeps track of workers, and handles task scheduling.
- **Worker**: Runs on each node in the cluster. It receives tasks from the manager and executes them using the Docker API.
- **Task**: The smallest unit of work, representing a container to be run.
- **Scheduler**: Responsible for deciding which worker node should run a particular task.
- **Node**: Represents a physical or virtual machine in the cluster.

## Goals

- Schedule containers on worker nodes efficiently.
- Start and stop containers using the Docker API.
- Manage a cluster of worker nodes and monitor their health.
- Provide a simplified interface for container orchestration.

## Prerequisites

- Go 1.26 or later
- Docker installed and running locally

## Getting Started

1. Clone the repository:

2. Build the project:

3. Run the main entry point (demonstrates container creation and removal):
   ```bash
   go run main.go
   ```
