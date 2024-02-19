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

###### Energy Testbed

###### Device Configuration

#### Directory Structure

```zsh
.
├── config
├── config.yaml
├── diode.go
├── diode_test.go
├── docker-compose.yaml
├── Dockerfile
├── go.mod
├── go.sum
├── Makefile
├── mqtt
├── Pipfile
├── Pipfile.lock
├── README.md
├── sample
└── utility

4 directories, 11 files

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

#### Risk Matrix ([5x5](https://safetyculture.com/topics/risk-assessment/5x5-risk-matrix/))

| ↔ Probability <br> Impact ↕ | **Insignificant** | **Minor** | **Significant** | **Major** | **Severe** |
| --------------------------- | ----------------- | --------- | --------------- | --------- | ---------- |
| **Almost Certain**          | R01               | R02       | R03             | R04       | R05        |
| **Likely**                  | R06               | R07       | R08             | R09       | R10        |
| **Moderate**                | R11               | R12       | R13             | R14       | R15        |
| **Unlikely**                | R16               | R17       | R18             | R19       | R20        |
| **Rare**                    | R21               | R22       | R23             | R24       | R25        |

#### Experimental Design

- [data/logbook.ipynb](data/logbook.ipynb)

[^1]: https://csrc.nist.gov/glossary/term/tactics_techniques_and_procedures
