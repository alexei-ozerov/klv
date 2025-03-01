# Kubernetes Log Viewer (KLV)

A golang TUI built with Bubble Tea to simplify log viewing.

![klv-demo](https://github.com/user-attachments/assets/65e6154c-ca70-4572-87bd-fff165d015f7)

## Todos

This is a very alpha version project with a lot of things to fix. Apologies in advance.

- [ ] Handle errors via an ErrorMsg which can be displayed to the user via a popup
- [ ] Remove all panic errors and replace them with a sensible message / fallback option
- [ ] Rework spaghetti code :3
- [ ] Implement channels or goroutines (need to decide which) to handle live log tailing
- [ ] Implement a search function for the log table
- [X] Complete first draft of README.md
- [X] Implement build pipeline

## Installation Instructions

### Prerequisites

You must have a kubeconfig file set in your default `.kube` folder.

### Installation

#### Linux

Navigate to the _Release_ section and pick the Linux release archive. It will contain a binary you can run from the
command line. Add it to your path to be able to execute it anywhere from your terminal.

#### MacOS

Navigate to the _Release_ section and pick the Darwin release archive. It will contain a binary you can run from the
command line. Add it to your path to be able to execute it anywhere from your terminal.

#### Windows

Navigate to the _Release_ section and pick the Windows release archive. It will contain a binary you can run from the
command line. Add it to your path to be able to execute it anywhere from your terminal.

## Usage Instructions

Make sure you have a valid and authenticated (if using GCloud Cli for example) Kubeconfig. Run the application. You will be presented with the first screen (the namespace / pod screen). Once you select a pod, you will be taken to the container / log screen.

### Controls

- j - scroll down
- k - scroll up
- r - reload current screen/
- enter, l - select table entry
- q - go back (if on container / log view), exit
