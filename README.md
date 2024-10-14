# hctl

<p align="center">
  <img alt="GitHub stars" src="https://img.shields.io/github/stars/xx4h/hctl">
  <img alt="GitHub forks" src="https://img.shields.io/github/forks/xx4h/hctl">
</p>

<!-- markdownlint-disable no-empty-links -->

[![Lint Code Base](https://github.com/xx4h/hctl/actions/workflows/linter-full.yml/badge.svg)](https://github.com/xx4h/hctl/actions/workflows/linter-full.yml)
[![Test Code Base](https://github.com/xx4h/hctl/actions/workflows/test-full.yml/badge.svg)](https://github.com/xx4h/hctl/actions/workflows/test-full.yml)
[![GitHub Downloads (all assets, all releases)](https://img.shields.io/github/downloads/xx4h/hctl/total)](https://github.com/xx4h/hctl/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/xx4h/hctl?)](https://goreportcard.com/report/github.com/xx4h/hctl)
[![codebeat badge](https://codebeat.co/badges/21ee1b92-b94c-4425-a600-b01dd4b1c045)](https://codebeat.co/projects/github-com-xx4h-hctl-main)
[![SLOC](https://tokei.rs/b1/github/xx4h/hctl?category=code&style=flat)](#)
[![Number of programming languages used](https://img.shields.io/github/languages/count/xx4h/hctl)](#)
[![Top programming languages used](https://img.shields.io/github/languages/top/xx4h/hctl)](#)
[![License](https://img.shields.io/badge/license-Apache--2.0-blue)](LICENSE)
[![Latest tag](https://img.shields.io/github/v/tag/xx4h/hctl)](https://github.com/xx4h/hctl/tags)
[![Closed issues](https://img.shields.io/github/issues-closed/xx4h/hctl?color=success)](https://github.com/xx4h/hctl/issues?q=is%3Aissue+is%3Aclosed)
[![Closed PRs](https://img.shields.io/github/issues-pr-closed/xx4h/hctl?color=success)](https://github.com/xx4h/hctl/pulls?q=is%3Apr+is%3Aclosed)
<br>

<!-- markdownlint-enable no-empty-links -->

hctl is a tool to control your Home Assistant devices from the command line

I needed a tool to quickly control my devices from the command line, focusing on easy to use and short commands to toggle or turn on/off lights, switches or even automations, play a mp3 from my local system, or change the volume of a media player.
And here we are!

## Features

<p align="center"><img alt="hctl showcase demo" src="/assets/demo.gif?raw=true"/></p>

- Support for Home Assistant
- Turn on/off, or toggle all capable devices
- Play local and remote music files
- Set volume on media players
- List all Domains & Domain-Services
- Completion for `bash`, `zsh`, `fish` and `powershell`, auto completing all capable devices
- Control over short and long names
- Fuzzy matching your devices so you can keep it short

## Install

### asdf

```bash
asdf plugin add hctl https://github.com/xx4h/asdf-hctl.git
asdf global hctl latest
```

for more information see [asdf-hctl](https://github.com/xx4h/asdf-hctl)

### Go

```bash
# version will be the latest tag, but will show version "dev"
go install github.com/xx4h/hctl@latest
```

### Release binary

Download the latest release binary from the [Release Page](https://github.com/xx4h/hctl/releases/latest) and extract it

### Build & Install from Source

```bash
git clone https://github.com/xx4h/hctl.git && cd hctl
make build && make local-install # intalls to ~/.local/bin/hctl
```

## Configuration

### Wizard

Run the init command

```bash
hctl init
```

### Manually

Copy the example config from this project

```yaml
# Configure Hub
hub:
  type: hass
  url: https://home-assistant.example.com/api
  token: YourToken
```

ensure the folder does already exist and edit with your favorite editor

```bash
mkdir -p ~/.config/hctl
$EDITOR ~/.config/hctl/hctl.yaml
```

## Completion

To really benefit from all features, ensure you've loaded the shell completion

```bash
# For bash (e.g. in your ~/.bashrc)
type hctl >/dev/null 2>&1 && source <(hctl completion bash)
```

For more information on how to setup completion for `bash`, `zsh`, `fish` and `PowerShell`, see `hctl completion -h`

**Optional**
Shorten command to a minimum

```bash
# this should at least work for bash and zsh
alias h='hctl'
source <(hctl completion bash | sed -e 's/hctl/h/g')

# afterwards toggling `switch.livingroom_warp` (with `Short Names` and `Fuzzy Matching` enabled) can be used like this
h t lw
```

## Usage

```bash
# Turn on all lights on Floor 1
hctl on floor1

# Toggle a switch called "some-switch"
hctl toggle some_switch

# Play a local music file
hctl play myplayer ~/path/to/some.mp3
```

### Completion Short Names

Home Assistant names its entities `domain.name`, like `light.some_light`.

```bash
# Imagine having the following devices/entities
light.livingroom_main
light.livingroom_corner
light.livingroom_other
switch.livingroom_warp

# Completion with Short Names feature enabled will auto complete them like
# e.g. if you want to turn off a switch you remeber starting with "sp"
hctl off li<TAB>
hclt off livingroom_<TAB><TAB>
livingroom_main     livingroom_corner      livingroom_other

# Without Short Names feature enabled they will be completed like
hctl off li<TAB>
hclt off light.<TAB><TAB>
light.livingroom_main     light.livingroom_corner      light.livingroom_other
```

Completion Short Names can be disabled with:

```yaml
completion:
  short_names: false
```

### Fuzzy Matching

```bash
# Imagine having the following devices
light.livingroom_main
light.livingroom_corner
light.livingroom_other
switch.livingroom_warp

# Turn on device with fuzzy matching (matching "switch.livingroom_warp")
hctl on lw
```

Fuzzy Matching is enabled by default.
Fuzzy Matching can be turned off in the config with:

```yaml
handling:
  fuzz: false
```

## What's Next / Roadmap

- [ ] Add more actions (like `press` e.g. Buttons, `trigger` e.g. Automations, or `lock` and `unlock` a Lock)
- [ ] Add output/feedback on actions (e.g. use pterm)
- [ ] Allow multiple devices on actions
- [ ] Add optional positional for `list entities`, following the same logic as in `toggle`, `on` and `off` (e.g. matching short names and fuzzy matching)
- [ ] Add possibility to add local mappings for devices in config
- [ ] Add install methods (native, asdf, ...)
