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

## License

[MIT](LICENSE) for this code _(but all used libraries may have different licence)_.
