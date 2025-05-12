# Rejoinderoo

<p align="center"><img src="images/logo.png"></p>

Rejoinderoo creates a rejoinder (response to reviewers) based on a CSV or Excel file.
The generated document is a LaTeX or Typst file that can be compiled to PDF.
An example of a generated rejoinder document is shown in [example.pdf](./example.pdf).

<p align="center"><img src="images/screenshot.png"></p>

## Development

This project uses a Makefile to manage all build and test tasks.

```sh
# for help and overview of all tasks
make help

# to install all dependencies
make deps

# to build the program in the `bin` directory
make build

# to run the compiled program
./bin/rejoinderoo
```

## Usage

TODO: Add GIF

### Color coding of responses

The response boxes are color-coded based on the ID field,
which is the first selected field.
To determine different reviewers, the prefix of the ID field value is used until the first delimiter (`.`, `-`, or `:`).
E.g., `Rev1.3` becomes `Rev1` and `R1:3` becomes `R1`.

In the next step, a custom LaTeX color is created for each reviewer that can be adjusted.

`\colorlet{colorRev1}{blue!15!white}`

## Customization

To customize the generated LaTeX file, you can either adjust [template.tex](./template.tex) or replace it with your own file.
The script will replace the placeholder `%%%%%custom-command%%%%%` and `%%%%%ccomments%%%%%` in the [template.tex](./template.tex) file with the generated content.

## License

Copyright © 2023-2025 Andreas Bauer

This work (source code) is licensed under [MIT](./LICENSE).
