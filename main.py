import csv
from pick import pick

csv_file_name = 'test.csv'
MIN_FIELDS = 3


def custom_tex_command(fieldnames):
    amount = len(fieldnames)
    cmd = f'''\\newcommand{{\\ccomment}}[{amount}]{{
    \\begin{{tcolorbox}}[title=#1, colback=white, coltitle=black, colbacktitle=black!15!white]
    \\textbf{{{fieldnames[1]}:}} #2 \\tcblower
    \\textbf{{{fieldnames[2]}:}} #3
    '''

    for i in range(3, amount):
        cmd += f'  \\\\ \\textbf{{{fieldnames[i]}:}} #{i+1} \n'

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
        selected_fields = [s[0] for s in selected]

        cmd = custom_tex_command(selected_fields)
        print(cmd)

except FileNotFoundError:
    print("Unable to find", csv_file_name)
