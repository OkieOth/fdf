
# fdf

[![CI](https://github.com/OkieOth/fdf/actions/workflows/ci.yml/badge.svg?branch=main&event=push)](https://github.com/OkieOth/fdf/actions/workflows/ci.yml)
[![go report card](https://goreportcard.com/badge/github.com/OkieOth/fdf)](https://goreportcard.com/report/github.com/OkieOth/fdf)

Simple tool to find files that are double on a local machine.



```bash
# build the binary from the source
go build -o fdf main.go

# showing help
./fdf --help

# find all duplicate files in your home directory, with use of only one CPU
./fdf list --cpus 1 --source ~/

# find all duplicate files from folder one in folder two, with use of all available CPUs
./fdf list --cpus 1 --source ~/one --searchRoot ~/two
```
