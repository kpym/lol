# lol a LaTeX online compiler CLI tool

lol is a small command line interface (CLI) that sends local files to distant server (https://latexonline.cc/ or [latex.ytotech.com](https://github.com/YtoTech/latex-on-http)) for LaTeX compilation and save the resulting `pdf`.

## Usage

To compile a single `main.tex` file to `main.pdf` using `latexonline.cc`:
```
> ./lol main.tex
```

To compile `main.tex` including `png` images in `imgs` folder with `xelatex` using `latex.ytotech.com`:
```
> ./lol -s ytotech -c xelatex main.tex imgs/*.png
```

A help message is provided:
```
> ./lol -h
lol (version: ---)
LaTeX online compiler. More info at www.github.com/kpym/lol.

Available options:
  -s, --service string    Service can be laton or ytotex.
  -c, --compiler string   One of pdflatex,xelatex or lualatex.
                          For ytotex platex, uplatex and context are also available.
                           (default "pdflatex")
  -f, --force             Do not use the laton cache. Force compile. Ignored by ytotech.
  -b, --biblio string     Can be bibtex or biber for ytotex. Not used by laton.
  -o, --output string     The name of the pdf file. If empty, same as the main tex file.
  -m, --main string       The main tex file to compile.
  -q, --quiet             Prevent any output.
  -v, --verbose           Print info and errors. No debug info is printed.
      --debug             Print everithing (debug info included).

Examples:
> lol main.tex
> lol  -s ytotech -c xelatex main.tex
> lol main.tex personal.sty images/img*.pdf
> cat main.tex | lol -c lualatex -o out.pdf
```

## Installation

### Precompiled executables

You can download the executable for your platform from the [Releases](https://github.com/kpym/lol/releases).

### Compile it yourself

#### Using Go

```
$ go get github.com/kpym/lol
```

#### Using goreleaser

After cloning this repo you can compile the sources with [goreleaser](https://github.com/goreleaser/goreleaser/) for all available platforms:

```
git clone https://github.com/kpym/lol.git .
goreleaser --snapshot --skip-publish --rm-dist
```

You will find the resulting binaries in the `dist/` sub-folder.

## Configuration

As `lol` use [viper](https://github.com/spf13/viper) the parameters can be provided not only by flags but also be read from config file (`lol.yaml`, `lol.toml`, `lol.json`...) or/and from environment variables (starting with `LOL_`).

### Using config file

You can provide all default values for flags in a config `lol` file in the current folder.
For example if your project needs `xelatex` and use `imgs/logo.png` you can save the following `lol.yaml` in the current folder
```yaml
Compiler: xelatex
Patterns:
  - imgs/logo.png
```

### Using environment variables

If you wan to provide global default values you can set an environment variable.
For example if you want by default to use `ytotech` service you can set `LOL_SERVICE=ytotech`.


## License

[MIT](LICENSE) for this code _(but all used libraries may have different licences)_.
