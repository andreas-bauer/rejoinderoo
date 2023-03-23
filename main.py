import argparse
from pathlib import Path
import csv
from pick import pick

CSV_FILE = 'input.csv'
TEMPLATE_FILE = 'template.tex'
MIN_FIELDS = 3
PLACEHOLDER_CMD = '%%%%%custom-command%%%%%'
PLACEHOLDER_CCOMMENT = '%%%%%ccomments%%%%%'


def custom_tex_command(fieldnames: list[str]) -> str:
    amount = len(fieldnames)
    colors = '\\colorlet{revColorDefault}{black!15!white}\n\n'
    cmd = colors + f'''\\newcommand{{\\ccomment}}[{amount}]{{
    \\begin{{tcolorbox}}[title=#1, colback=white, coltitle=black, colbacktitle=revColorDefault]
    \\textbf{{{fieldnames[1]}:}} #2 \\tcblower
    \\textbf{{{fieldnames[2]}:}} #3
    '''

    for i in range(3, amount):
        cmd += f'  \\\\ \\textbf{{{fieldnames[i]}:}} #{i+1} \n'

    cmd += '\\end{tcolorbox}\n}\n'
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


def create_box(row, fieldnames: list[str]) -> str:
    box_cmd = '\n\n\\ccomment'

    for field in fieldnames:
        tex = to_latex(row[field])
        box_cmd += f'{{\n{tex}\n}}'
    return box_cmd


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description='desc')
    parser.add_argument('comment_file', type=Path, help='The comment csv file')
    parser.add_argument('tex_file', type=Path, help='The TeX output file')

    args = parser.parse_args()
    box_cmd = ''
    cmd = ''

    with open(args.comment_file, 'r', encoding='utf8') as csv_file:
        reader = csv.DictReader(csv_file)

        if len(reader.fieldnames) < MIN_FIELDS:
            print(f'CSV file must have at least {MIN_FIELDS} columns but has {len(reader.fieldnames)}.\nUnable to proceed.')
            exit()

        selected = pick(reader.fieldnames,
                        'Select at least three fields (ID, reviewer comment, author response)',
                        multiselect=True, min_selection_count=MIN_FIELDS)
        selected_fields = [s[0] for s in selected]

        cmd = custom_tex_command(selected_fields)

        for row in reader:
            box_cmd += create_box(row, selected_fields)

    final = ''
    with open(TEMPLATE_FILE, 'r', encoding='utf8') as template:
        for line in template:
            replaced = line.replace(PLACEHOLDER_CMD, cmd)
            replaced = replaced.replace(PLACEHOLDER_CCOMMENT, box_cmd)
            final += replaced
    args.tex_file.parent.mkdir(parents=True, exist_ok=True)
    with open(args.tex_file, 'w', encoding='utf8') as target:
        target.write(final)
    print('Done! Created LaTeX file:', args.tex_file)
