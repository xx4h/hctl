<!-- markdownlint-disable no-empty-links -->

[![Lint Code Base](https://github.com/xx4h/hctl/actions/workflows/linter-full.yml/badge.svg)](https://github.com/xx4h/hctl/actions/workflows/linter-full.yml)
[![SLOC](https://tokei.rs/b1/github/xx4h/hctl?category=code&style=flat)](#)
[![Number of programming languages used](https://img.shields.io/github/languages/count/xx4h/hctl)](#)
[![Top programming languages used](https://img.shields.io/github/languages/top/xx4h/hctl)](#)
[![Open issues](https://img.shields.io/github/issues/xx4h/hctl)](https://github.com/xx4h/hctl/issues)
[![Closed issues](https://img.shields.io/github/issues-closed/xx4h/hctl?color=success)](https://github.com/xx4h/hctl/issues?q=is%3Aissue+is%3Aclosed)
[![Closed PRs](https://img.shields.io/github/issues-pr-closed/xx4h/hctl?color=success)](https://github.com/xx4h/hctl/pulls?q=is%3Apr+is%3Aclosed)
[![License](https://img.shields.io/badge/license-Apache--2.0-blue)](LICENSE.txt)
[![Latest tag](https://img.shields.io/github/v/tag/xx4h/hctl)](https://github.com/xx4h/hctl/tags)
<br>

<!-- markdownlint-enable no-empty-links -->

# hctl

hctl is a tool to control your Home Assistant (and maybe more in the future) devices from the command line

## Features

- Support for Home Assistant
- List all Domains & Domain-Services
- Turn on/off, or toggle all capable devices
- Completion for `bash`, `zsh`, `fish` and `powershell`, auto completing all capable devices
- Control over short and long names
- Fuzzy matching your devices so you can keep it short

## Install

### Go

```bash
go install github.com/xx4h/hctl@latest
```

## Configuration

Run the init command

```bash
hctl init
```

or copy the example config from this project

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

To really benefit from all features, ensure you've loaded the bash completion

```bash
# Bash (e.g. your ~/.bashrc)
type hctl >/dev/null 2>&1 && source <(hctl completion bash)
```

## Usage

```bash
# Turn on all lights on Floor 1
hctl on floor1

# Toggle a switch called "some-switch"
hctl toggle some_switch

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

## Roadmap

- [ ] Add more actions (like `press` e.g. Buttons, `trigger` e.g. Automations, or `lock` and `unlock` a Lock)
- [ ] Add `config` command to actively set config options in the config file
- [ ] Add output/feedback on actions (e.g. use pterm)
- [ ] Allow multiple devices on actions
- [ ] Improve output and add filters to `list` (e.g. use pterm)
- [ ] Add optional positional for `list entities`, following the same logic as in `toggle`, `on` and `off` (e.g. matching short names and fuzzy matching)
- [ ] Add possibility to add local mappings for devices in config
- [ ] Add install methods (native, asdf, ...)
