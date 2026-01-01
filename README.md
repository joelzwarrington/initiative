# initiative

[![Latest Release](https://img.shields.io/github/release/joelzwarrington/initiative.svg)](https://github.com/joelzwarrington/initiative/releases)
[![GoDoc](https://godoc.org/github.com/joelzwarrington/initiative?status.svg)](https://godoc.org/github.com/joelzwarrington/initiative)
[![Build Status](https://github.com/joelzwarrington/initiative/workflows/ci/badge.svg)](https://github.com/joelzwarrington/initiative/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/joelzwarrington/initiative)](https://goreportcard.com/report/github.com/joelzwarrington/initiative)

A powerful CLI for running **Dungeons & Dragons** encounters. Track initiative, health, effects, and turns with an intuitive terminal interface.

**initiative** is built with [Cobra](https://github.com/spf13/cobra) and the [Charm](https://charm.land/) ecosystem: [Bubble Tea](https://github.com/charmbracelet/bubbletea), [Bubbles](https://github.com/charmbracelet/bubbles), [Lipgloss](https://github.com/charmbracelet/lipgloss), [Huh?](https://github.com/charmbracelet/huh), and [more](#built-with).

## Getting started

1. Install the application using one of the below methods
2. Run `initiative` in your terminal!

_Review our [wiki](https://github.com/joelzwarrington/initiative/wiki) for additional documentation, and submit any questions to our [Q&A discussion category](https://github.com/joelzwarrington/initiative/discussions/new?category=q-a)_.

### Script (recommended)

The installation script will use the built artifacts from the releases based on your detected operating system and architecture.

```bash
curl -sSL https://raw.githubusercontent.com/joelzwarrington/initiative/main/scripts/install.sh | bash
```

### Artifacts

Download the artifact from [releases](https://github.com/joelzwarrington/initiative/releases) for your specific operating system and architecture.

### Go

The binaries can be built from source using the `go install` utility.

```bash
go install github.com/joelzwarrington/initiative@latest
```

## Encounters

<img src="examples/encounters.gif" width="640" alt="Setting up an encounter">

Setup encounters with characters and monsters by providing initiative, quantities.

## Characters

<img src="examples/characters.gif" width="640" alt="Managing characters">

Organize your party members and NPCs, which can be added to encounters quickly.

## Initiative

<!-- todo: implement tape -->
<img src="examples/initiative.gif" width="640" alt="Tracking initiative">

Manage turn order and initiative tracking during combat encounters.

## Health

<img src="examples/health.gif" width="640" alt="Tracking monster health">

Track monster hit points throughout the encounter.

## Saved game files

<img src="examples/saves.gif" width="640" alt="Game file persistence">

Your party data saves automatically and persists across sessions. Manage multiple campaigns with custom file paths using the `--game` flag.

## Sources

<img src="examples/sources.gif" width="640" alt="Using custom sources">

The [System Reference Document](https://www.dndbeyond.com/srd) comes built-in with all core monsters. Expand your creature library with custom sources using the `--source` flag.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) to get started.

## Feedback

We'd love to hear from you!

- **Bugs**: [Submit a bug report](https://github.com/joelzwarrington/initiative/issues/new?template=bug_report.yml) for bugs or problems
- **Ideas**: [Share your ideas](https://github.com/joelzwarrington/initiative/discussions/new?category=ideas) for new features or improvements
- **Questions**: [Ask questions](https://github.com/joelzwarrington/initiative/discussions/new?category=q-a) about usage or functionality

## Built With

- [[Cobra](https://github.com/spf13/cobra)] library for creating powerful modern CLI applications
- [[Bubble Tea](https://github.com/charmbracelet/bubbletea)] powerful terminal user interface framework
- [[Skeleton](https://github.com/joelzwarrington/skeleton)] tab framework for Bubble Tea applications
- [[Bubbles](https://github.com/charmbracelet/bubbles)] components for Bubble Tea applications
- [[Lipgloss](https://github.com/charmbracelet/lipgloss)] beautiful terminal styling
- [[Huh?](https://github.com/charmbracelet/huh)] interactive prompts and forms

## License

[MIT](https://github.com/joelzwarrington/initiative/raw/main/LICENSE)

In addition, this work also includes material from the System Reference Document 5.2.1 ("SRD 5.2.1") by Wizards of the Coast LLC, available at https://www.dndbeyond.com/srd. The SRD 5.2.1 is licensed under the Creative Commons Attribution 4.0 International License, available at https://creativecommons.org/licenses/by/4.0/legalcode.

---

Built with ♥ by [@joelzwarrington](https://github.com/joelzwarrington)
