# aoc-cli

My solution for wanting to download the data for particular days with minimal effort.

## Usage 

No flags are required. When called, the required information will be asked for.

```
-year, -y int
The year of advent of code to query.

-day, -d int
The day of the advent of code to query.

-filepath, -f string
The location to download the input file to (must exist)

-filename, -n string
The filename to download to

-stdout
If true will print the output to stdout rather than downloading to a file
```

## Config

To make life easier, access token and download location can be specified in `config.json` in either the directory called from or $config/aoc/

#### $config

- Linux `XDG_CONFIG_HOME` or `~/.config/`
- Windows `%AppData%`

