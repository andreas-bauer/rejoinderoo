import csv
from pick import pick

csv_file_name = 'test.csv'
MIN_FIELDS = 3


def custom_tex_command(selected):
    amount = len(selected)
    cmd = f'''\\newcommand{{\\ccomment}}[{amount}]{{
    \\begin{{tcolorbox}}[title=#1, colback=white, coltitle=black, colbacktitle=black!15!white]
    \\textbf{{{selected[1][0]}:}} #2 \\tcblower
    \\textbf{{{selected[2][0]}:}} #3
    '''

    for i in range(3, amount):
        cmd += f'  \\\\ \\textbf{{{selected[i][0]}:}} #{i+1} \n'

    cmd += '\\end{tcolorbox}\n}'
    return cmd


try:
    with open(csv_file_name, 'r', encoding='utf8') as csv_file:
        reader = csv.DictReader(csv_file)

        if len(reader.fieldnames) < MIN_FIELDS:
            print(f'CSV file must have at least {MIN_FIELDS} columns but has {len(reader.fieldnames)}.\nUnable to proceed.')
            exit()

        selected = pick(reader.fieldnames,
            'Select at least three fields (ID, reviewer comment, author response)',
            multiselect=True, min_selection_count=MIN_FIELDS)

        cmd = custom_tex_command(selected)
        print(cmd)

except FileNotFoundError:
    print("Unable to find", csv_file_name)
