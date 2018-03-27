**GoDepo** is a dependency management tool for Golang.

# Rationale

This tool was written after using dep, govendor, godep and glide and none
of them fitted my requirements and requirements of company where I work.

Main reason of creating this tool was ability to control dependencies
without mangling with system tools, like SSH configuration.

# Functionality

See [doc/features.md](doc/features.md).

# Installation

It is enough to say:

```
go get github.com/pztrn/godepo/cmd/...
```

GoDepo was developed using latest available Golang version and it is
highly recommended that you'll use it if you compile from sources.

Otherwise just grab a binary from Releases section.

# Configuration

GoDepo can be configured with YAML-formatted file. For complete reference
of available options as well as example file take a look at
[doc/configuration.md](doc/configuration.md)

# Developing