import argparse
import csv
import os
import sys
from pathlib import Path

from pick import pick

CSV_FILE = 'input.csv'
TEMPLATE_FILE = 'template.tex'
MIN_FIELDS = 3
PLACEHOLDER_CMD = '%%%%%custom-command%%%%%'
PLACEHOLDER_CCOMMENT = '%%%%%ccomments%%%%%'
COLOR_REV_DEFAULT = 'colorRevDefault'


def custom_tex_command(fieldnames: list[str]) -> str:
    amount = len(fieldnames)
    colors = '\\colorlet{' + COLOR_REV_DEFAULT + '}{black!15!white}\n\n'
    cmd = (
        colors
        + f'''\\newcommand{{\\ccomment}}[{amount + 1}]{{
    \\begin{{tcolorbox}}[title=#2, breakable, colback=white, coltitle=black, colbacktitle=#1]
    \\textbf{{{fieldnames[1]}:}} #3 \\tcblower
    \\textbf{{{fieldnames[2]}:}} #4
    '''
    )

    for i in range(3, amount):
        cmd += f'  \\\\ \\textbf{{{fieldnames[i]}:}} #{i+2} \n'

    cmd += '\\end{tcolorbox}\n}\n'
    return cmd


def color_command(rev_ids: list[str]) -> str:
    cmd = ''
    for rev_id in rev_ids:
        cmd += f'\\colorlet{{color{rev_id}}}{{black!15!white}}\n'
    return cmd


def to_latex(s: str) -> str:
    replacements = {
        '&': '\\&',
        '%': '\\%',
        '$': '\\$',
        '#': '\\#',
        '_': '\\_',
        '{': '\\{',
        '}': '\\}',
        '~': '\\textasciitilde',
        '^': '\\textasciicircum',
        '\\': '\\textbackslash',
        '|': '\\textbar',
        '<': '\\textless',
        '>': '\\textgreater',
        '[': '{[}',
        ']': '{]}',
    }

    tex = "".join(replacements.get(c, c) for c in s)

    tex = tex.replace('\n-', '\n\\\\-')
    tex = tex.replace('\n -', '\n\\\\ -')

    return tex


def extract_rev_id(full_id: str) -> str:
    return full_id.split(".")[0].split("-")[0].split(":")[0]


def create_box(row, fieldnames: list[str], reviewer_ids: set) -> str:
    box_cmd = '\n\n\\ccomment'

    rev_id = extract_rev_id(row[fieldnames[0]])
    if rev_id in reviewer_ids:
        box_cmd += f'{{\ncolor{rev_id}\n}}'
    else:
        box_cmd += f'{{\n{COLOR_REV_DEFAULT}\n}}'

    for field in fieldnames:
        tex = to_latex(row[field])
        tex = tex.rstrip()
        box_cmd += f'{{ % {field}\n{tex}\n}}'
    return box_cmd


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description='desc')
    parser.add_argument('comment_file', type=Path, help='The comment csv file')
    parser.add_argument('tex_file', type=Path, help='The TeX output file')

    args = parser.parse_args()
    box_cmd = ''
    cmd = ''
    color_cmd = ''

    with open(args.comment_file, 'r', encoding='utf8') as csv_file:
        reader = csv.DictReader(csv_file)

        if reader.fieldnames is None:
            print("Unable to detect columns in CSV file")
            sys.exit()

        if len(reader.fieldnames) < MIN_FIELDS:
            print(
                f'CSV file must have at least {MIN_FIELDS} columns but has {len(reader.fieldnames)}'
                + '.\nUnable to proceed.'
            )
            sys.exit()

        selected = pick(
            reader.fieldnames,
            'Select at least three fields (ID, reviewer comment, author response)',
            multiselect=True,
            min_selection_count=MIN_FIELDS,
        )
        selected_fields = [s[0] for s in selected]

        cmd = custom_tex_command(selected_fields)

        reviewer_ids = []
        for row in reader:
            rev_id = extract_rev_id(row[selected_fields[0]])
            if rev_id not in reviewer_ids:
                reviewer_ids.append(rev_id)

            box_cmd += create_box(row, selected_fields, reviewer_ids)

        color_cmd = color_command(reviewer_ids)

    final = ''
    with open(TEMPLATE_FILE, 'r', encoding='utf8') as template:
        for line in template:
            replaced = line.replace(PLACEHOLDER_CMD, color_cmd + cmd)
            replaced = replaced.replace(PLACEHOLDER_CCOMMENT, box_cmd)
            final += replaced
    args.tex_file.parent.mkdir(parents=True, exist_ok=True)
    with open(args.tex_file, 'w', encoding='utf8') as target:
        target.write(final)
    print(f'Done! Created LaTeX document: file://{os.getcwd()}/{args.tex_file}')
