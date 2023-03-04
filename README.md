# Rejoinderoo

<p align="center"><img src="logo.png"></p>

Rejoinderoo creates a rejoinder (response to reviewers) LaTeX document based on a CSV file.

## Prerequisites

Rejoinderoo depends on Pick to provide the selection interface for
data fields of the CSV file.

`pip3 install pick`

## Usage

Step 1) Prepare the response to reviewers in a spreadsheet and export it as a CSV file,
like the [input.csv](./input.csv).

At least three columns are required to be able to parse the CSV file.
Have a look at  as an example of

Step 2) (optional) Copy the exported CSV file in this folder and rename it to `input.csv`.

Step 3) Run `main.py` without any arguments to use `input.csv` as input,
or specify the target file as an argument.

```sh
# use input.csv
python3 main.py

# use a specific file
python3 main.py my_file.csv
```

Step 4) The created LaTeX file has the name `out.tex`

## Customization

To customize the generated LaTeX file, you can either adjust [template.tex](./template.tex) or replace it with your own file.
The script will replace the placeholder `%%%%%custom-command%%%%%` and `%%%%%ccomments%%%%%` in the [template.tex](./template.tex) file with the generated content.

## License

Copyright (c) 2023 Andreas Bauer

This work (source code) is licensed under  [MIT](./LICENSE).
