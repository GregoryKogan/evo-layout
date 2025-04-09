import json
import os

PROBLEM = "TSP"
ALGO = "SSGA"

script_path = os.path.dirname(os.path.realpath(__file__))
logfile = os.path.join(script_path, f"{PROBLEM}_{ALGO}.jsonl")

csv_lines = [f"elapsed;{ALGO}_objective\n"]

with open(logfile) as f:
    f.readline()
    lines = f.readlines()

for line in lines:
    data = json.loads(line)
    csv_lines.append(
        f'{int(data["elapsed"]) / 1e9};{data["solution"]["objectives"][0]}\n'
        # f'{int(data["elapsed"]) / 1e9};{data["solution"]["fitness"]}\n'
    )

output = os.path.join(script_path, f"{PROBLEM}_{ALGO}.csv")
with open(output, "w") as f:
    f.writelines(csv_lines)
