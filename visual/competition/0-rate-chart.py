import pandas as pd
import matplotlib.pyplot as plt
import numpy as np

plt.style.use("seaborn-v0_8")

# Data provided in markdown table
data = """
| Алгоритм      | ПЛ25    | ПЛ50    | ПЛ100   | ПЛ200  |
|---------------|---------|---------|---------|--------|
| SGA           | 90%     | 34%     | 0%      | 0%     |
| SSGA          | 94%     | 38%     | 0%      | 0%     |
| NSGA2         | 94%     | 24%     | 0%      | 0%     |
| SPEA2         | 92%     | 2%      | 0%      | 0%     |
| FR            | 0.7%    | 0.0%    | 0.0%    | 0.0%   |
| SSGA-FR       | 3.65%   | 0.00%   | 0.00%   | 0.00%  |
| FR-NSGA2      | 100.00% | 100.00% | 100.00% | 44.35% |
| FR-SSGA-NSGA2 | 73.74%  | 26.83%  | 0.00%   | 0.00%  |
"""

# Convert markdown table to a list of lists, then to a pandas DataFrame
lines = data.strip().split("\n")
header = [h.strip() for h in lines[0].split("|") if h.strip()]
values = []
for line in lines[2:]:
    if line.strip():
        values.append([v.strip() for v in line.split("|") if v.strip()])

df = pd.DataFrame(values, columns=header)

# Convert percentage columns to numeric
for col in df.columns[1:]:  # Skip 'Алгоритм' column
    df[col] = df[col].str.replace("%", "").astype(float)

# Set 'Алгоритм' as index for easier plotting
df = df.set_index("Алгоритм")

# Plotting
fig, ax = plt.subplots(figsize=(15, 7))

bar_width = 0.1
# Positions of the bars on the x-axis
r = np.arange(len(df.columns))

for i, (algorithm, row) in enumerate(df.iterrows()):
    # Adjust bar positions for grouping
    offset = bar_width * i
    ax.bar(r + offset, row.values, width=bar_width, label=algorithm)

# Add labels, title, and legend
ax.set_ylabel("Доля распутанных решений (%)", fontsize=20, fontweight="bold")
ax.set_xticks(r + bar_width * (len(df.index) - 1) / 2)  # Center x-ticks
ax.set_xticklabels(
    [
        "25 вершин",
        "50 вершин",
        "100 вершин",
        "200 вершин",
    ],
    fontsize=20,
    fontweight="bold",
)

fig.legend(
    loc="lower center",
    ncol=8,
    prop={"weight": "bold", "size": 16},
)  # Adjust ncol as needed
ax.grid(axis="y", linestyle="--", alpha=0.7)
plt.tight_layout(rect=[0, 0.05, 1, 0.95])

assets_dir = "/Users/gregorykogan/Documents/МИФИ/6 семестр/УИР/ПЗ/assets"
plt.savefig(f"{assets_dir}/0-rate-final-competition.png", dpi=300)
plt.show()
