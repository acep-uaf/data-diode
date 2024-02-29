# Data Diode

Scripts for verifying TCP passthrough functionality.

## Development Instructions

> [!NOTE]
> This project utilizes [`go`](https://go.dev/) for module management.
> You can find installation instructions for `1.22.0` via https://go.dev/doc/install.

- Clone repository: `gh repo clone acep-uaf/data-diode`
- Source navigation: `cd data-diode`
- Build binary: `make`
- CLI: `./diode [options...]`

#### Architecture Diagram

#### Directory Structure

```zsh
.
├── config
├── config.yaml
├── diode.go
├── diode_test.go
├── docker-compose.yaml
├── Dockerfile
├── docs
├── go.mod
├── go.sum
├── Makefile
├── mqtt
├── Pipfile
├── Pipfile.lock
├── README.md
├── sample
└── utility

5 directories, 11 files

```

## User Stories

#### Scenario Planning

1. Power Plant Operator
1. Information Security Auditor
1. Energy Awareness Application Developer
1. Community Member

#### Threat Model[^1]

- [ ] Tactics
- [ ] Techniques
- [ ] Procedures

## System Benchmarking

#### Experimental Design

- [docs/SOP.md](docs/SOP.md)

###### Device Configuration

[^1]: https://csrc.nist.gov/glossary/term/tactics_techniques_and_procedures
