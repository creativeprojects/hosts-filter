[![Build](https://github.com/creativeprojects/hosts-filter/workflows/Build/badge.svg)](https://github.com/creativeprojects/hosts-filter/actions)

# Introduction

Compile and update your hosts file from various block lists. Works on various unixes, Linux, macOS X and Windows.

This tool simply download lists of bad domains (adware, spyware, phishing, scam, etc.) and adds an entry into your hosts file like:
```
0.0.0.0 bad.domain
```
which means your computer cannot find the malicious domain to load as it points to an invalid IP (`0.0.0.0`).

# Configuration

Put all the block list URL into a config file `config.yaml`:
```yaml
---
# hosts file to update (leave blank for the system default)
# hosts_file: /etc/hosts

# IP used for filtered domains: default is 0.0.0.0
# ip: 0.0.0.0

block_lists:
  -
    url: http://winhelp2002.mvps.org/hosts.txt
  -
    url: http://someonewhocares.org/hosts/hosts
  -
    url: http://pgl.yoyo.org/adservers/serverlist.php?hostformat=hosts&mimetype=plaintext

# list of domains to remove from the lists
# allow:
#   - github.com

# file containing a list of domains allowed (one per line)
# allow_from: allow.txt

```

# Usage

```
hosts-filter [flags]

hosts-filter flags:
  -c, --config string   configuration file (default "config.yaml")
  -h, --help            display this help
  -l, --log string      logs into a file instead of the console
      --no-ansi         disable ansi control characters (disable console colouring)
  -q, --quiet           display only warnings and errors
  -r, --remove          clear up the hosts file of all entries generated by hosts-filter
  -o, --stdout          don't save the hosts file in place, but send it to the standard output instead
  -v, --verbose         display some debugging information

```

## Compile block lists and update your hosts file

```
$ sudo hosts-filter

loaded "http://winhelp2002.mvps.org/hosts.txt": 8815 entries in total
loaded "http://someonewhocares.org/hosts/hosts": 22069 entries in total
loaded "http://pgl.yoyo.org/adservers/serverlist.php?hostformat=hosts&mimetype=plaintext": 24890 entries in total
```

## Remove the block lists from your hosts file

```
$ sudo hosts-filter -r
```


# Installation (macOS, Linux & other unixes)

Here's a simple script to download the binary automatically. It works on mac OS X, FreeBSD and Linux:

```
$ curl -sfL https://raw.githubusercontent.com/creativeprojects/hosts-filter/main/install.sh | sh
```

It should copy hosts-filter in a `bin` directory under your current directory.

If you need more control, you can save the shell script and run it manually:

```
$ curl -LO https://raw.githubusercontent.com/creativeprojects/hosts-filter/main/install.sh
$ chmod +x install.sh
$ sudo ./install.sh -b /usr/local/bin
```

It will install hosts-filter in `/usr/local/bin/`


## Installation for Windows using bash

You can use the same script if you're using bash in Windows (via WSL, git bash, etc.)

```
$ curl -LO https://raw.githubusercontent.com/creativeprojects/hosts-filter/main/install.sh
$ ./install.sh
```
It will create a `bin` directory under your current directory and place `hosts-filter.exe` in it.

## Manual installation (Windows)

- Download the package corresponding to your system and CPU from the [release page](https://github.com/creativeprojects/hosts-filter/releases)
- Once downloaded you need to open the archive and copy the binary file `hosts-filter` (or `hosts-filter.exe`) in your PATH.

# Using docker image ##

You can run hosts-filter inside a docker container. It is probably the easiest way to install hosts-filter and keep it updated.

By default, the hosts-filter container starts at `/hosts-filter`. So you can feed a configuration this way:

```
$ docker run -it --rm -v $PWD/examples:/hosts-filter creativeprojects/hosts-filter
```


# License

This work is licensed under the Creative Commons
Attribution-NonCommercial-ShareAlike License.
https://creativecommons.org/licenses/by-nc-sa/4.0/
