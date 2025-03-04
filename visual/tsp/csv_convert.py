import json
import os

script_path = os.path.dirname(os.path.realpath(__file__))
logfile = os.path.join(script_path, "TSP.jsonl")

csv_lines = ["elapsed,generation,fitness\n"]

with open(logfile) as f:
    f.readline()
    lines = f.readlines()

max_fitness = json.loads(lines[-1])["solution"]["fitness"]
for line in lines:
    data = json.loads(line)
    if "took" in data:
        max_fitness = data["solution"]["fitness"]
        continue
    csv_lines.append(
        f'{int(data["elapsed"]) / 1e9},{data["generation"]},{data["solution"]["fitness"] / max_fitness * 100}\n'
    )

output = os.path.join(script_path, "tsp.csv")
with open(output, "w") as f:
    f.writelines(csv_lines)
