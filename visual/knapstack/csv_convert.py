import json
import os

script_path = os.path.dirname(os.path.realpath(__file__))
logfile = os.path.join(script_path, "knapsack.jsonl")

csv_lines = ["generation,fitness\n"]

with open(logfile) as f:
    f.readline()
    for line in f:
        data = json.loads(line)
        csv_lines.append(f'{data["generation"]},{data["solution"]["fitness"]}\n')

output = os.path.join(script_path, "knapsack.csv")
with open(output, "w") as f:
    f.writelines(csv_lines)
