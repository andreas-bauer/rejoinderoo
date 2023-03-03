import sys
import csv
from pick import pick

CSV_FILE = 'input.csv'
TEMPLATE_FILE = 'template.tex'
TARGET_FILE = 'out.tex'
MIN_FIELDS = 3
PLACEHOLDER_CMD = '%%%%%custom-command%%%%%\n'
PLACEHOLDER_CCOMMENT = '%%%%%ccomments%%%%%\n'


def custom_tex_command(fieldnames):
    amount = len(fieldnames)
    cmd = f'''\\newcommand{{\\ccomment}}[{amount}]{{
    \\begin{{tcolorbox}}[title=#1, colback=white, coltitle=black, colbacktitle=black!15!white]
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


def create_box(row, fieldnames):
    box_cmd = '\n\n\\ccomment'

    for field in fieldnames:
        tex = to_latex(row[field])
        box_cmd += f'{{\n{tex}\n}}'
    return box_cmd


cmd = ''
box_cmd = ''

has_arguments = len(sys.argv) > 1
if has_arguments:
    CSV_FILE = sys.argv[1]

try:
    with open(CSV_FILE, 'r', encoding='utf8') as csv_file:
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

except FileNotFoundError:
    print('Unable to find', CSV_FILE)

try:
    final = ''
    with open(TEMPLATE_FILE, 'r', encoding='utf8') as template:
        for line in template:
            if line == PLACEHOLDER_CMD:
                line = cmd

            if line == PLACEHOLDER_CCOMMENT:
                line = box_cmd

            final += line

    with open(TARGET_FILE, 'w', encoding='utf8') as target:
        target.write(final)
        print('Done! Created LaTeX file:', TARGET_FILE)

except FileNotFoundError:
    print('Unable to find', TEMPLATE_FILE)
