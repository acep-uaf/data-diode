# Data Diode

Scripts for verifying TCP passthrough functionality.

## Development Instructions

> [!TIP]
> This project utilizes [`go`](https://go.dev/) for module management.
> You can find installation instructions via [relevant operating system](https://go.dev/doc/install) documentation.

- Clone repository: `gh repo clone acep-uaf/data-diode`
- Source navigation: `cd data-diode`
- Build binary: [`make`](https://www.gnu.org/software/make/)
- CLI: `./diode [options...]`

### Branch Management

- `main` → production ready environment.
- `dev` → testing changes to be merged into `main`.

### Directory Structure

```zsh
.
├── config
├── docker-compose.yaml
├── Dockerfile
├── docs
├── go.mod
├── go.sum
├── insights
├── main.go
├── Makefile
├── README.md
├── sample
└── utility

5 directories, 7 files
```

### Architecture Diagram

```mermaid
graph LR
    A("Subscribe (MQTT)") -->|TCP Client|B(Data Diode) -->|TCP Server|C("Publish (MQTT)")

```

> [!NOTE]
> Operational Technology (OT) vs. Information Technology (IT) system boundaries.

## User Stories

### Scenario Planning

1. Power Plant Operator
1. Information Security Auditor
1. Energy Awareness Application Developer
1. Community Member

### Threat Model[^1]

- [ ] Tactics
- [ ] Techniques
- [ ] Procedures

### System Benchmarking

### Experimental Design

> Research Question: Can you replicate an OT (private) system state published via the MQTT protocol to an IT (public) environment through a data diode?

#### [Device Configuration](docs/SOP.md)

[^1]: <https://csrc.nist.gov/glossary/term/tactics_techniques_and_procedures>
