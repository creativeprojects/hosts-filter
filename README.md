[![Build](https://github.com/creativeprojects/hosts-filter/workflows/Build/badge.svg)](https://github.com/creativeprojects/go-selfupdate/actions)

# Introduction

Compile and update your hosts file from various block lists. Works on various unixes, Linux, macOS X and Windows.

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

# License

This work is licensed under the Creative Commons
Attribution-NonCommercial-ShareAlike License.
https://creativecommons.org/licenses/by-nc-sa/4.0/
