import matplotlib.pyplot as plt
import numpy as np

# Solarize_Light2
# _classic_test_patch
# _mpl-gallery
# _mpl-gallery-nogrid
# bmh
# classic
# dark_background
# fast
# fivethirtyeight
# ggplot
# grayscale
# petroff10
# seaborn-v0_8
# seaborn-v0_8-bright
# seaborn-v0_8-colorblind
# seaborn-v0_8-dark
# seaborn-v0_8-dark-palette
# seaborn-v0_8-darkgrid
# seaborn-v0_8-deep
# seaborn-v0_8-muted
# seaborn-v0_8-notebook
# seaborn-v0_8-paper
# seaborn-v0_8-pastel
# seaborn-v0_8-poster
# seaborn-v0_8-talk
# seaborn-v0_8-ticks
# seaborn-v0_8-white
# seaborn-v0_8-whitegrid
# tableau-colorblind10

plt.style.use("seaborn-v0_8")

# Data provided
data = {
    "Алгоритм": [
        "SGA",
        "SSGA",
        "NSGA2",
        "SPEA2",
        "FR",
        "SSGA-FR",
        "FR-NSGA2",
        "FR-SSGA-NSGA2",
    ],
    "ПЛ25": [1.02, 0.68, 0.64, 0.84, 17.04, 17.87, 0.00, 3.84],
    "ПР25": [81.66, 79.42, 85.60, 97.52, 170.72, 159.89, 90.33, 82.11],
    "ПЛ50": [29.74, 26.96, 31.56, 46.96, 50.56, 53.97, 0.00, 31.87],
    "ПР50": [431.44, 388.78, 441.78, 450.56, 649.68, 711.33, 379.33, 355.78],
    "ПЛ100": [366.50, 304.02, 288.20, 385.00, 142.78, 137.78, 0.00, 75.67],
    "ПР100": [2139.67, 2021.22, 2280.44, 2509.89, 2569.40, 2542.56, 1621.22, 1582.78],
    "ПЛ200": [2539.22, 2181.33, 2105.44, 2433.22, 357.46, 370.78, 43.56, 861.00],
    "ПР200": [
        10571.22,
        10277.78,
        10326.22,
        10996.11,
        10239.40,
        10222.56,
        6575.22,
        7095.89,
    ],
}

algorithms = data["Алгоритм"]

# Define metric groups
metric_groups = {
    "25": ["ПЛ25", "ПР25"],
    "50": ["ПЛ50", "ПР50"],
    "100": ["ПЛ100", "ПР100"],
    "200": ["ПЛ200", "ПР200"],
}

# Set up the figure and subplots
fig, axes = plt.subplots(1, 4, figsize=(20, 8), sharey=False)  # Adjusted figure height

# Iterate through each metric group and create a subplot
for i, (postfix, metrics) in enumerate(metric_groups.items()):
    ax = axes[i]
    x = np.arange(len(metrics))  # the label locations
    width = 0.11  # the width of the bars

    # Get data for each algorithm for the current metric group
    sga_data = [data[metric][0] for metric in metrics]
    ssga_data = [data[metric][1] for metric in metrics]
    nsga2_data = [data[metric][2] for metric in metrics]
    spea2_data = [data[metric][3] for metric in metrics]
    fr_data = [data[metric][4] for metric in metrics]
    ssga_fr_data = [data[metric][5] for metric in metrics]
    fr_nsga2_data = [data[metric][6] for metric in metrics]
    fr_ssga_nsga2_data = [data[metric][7] for metric in metrics]

    # Create bars for each algorithm
    ax.bar(x - 3.5 * width, sga_data, width, label="SGA")
    ax.bar(x - 2.5 * width, ssga_data, width, label="SSGA")
    ax.bar(x - 1.5 * width, nsga2_data, width, label="NSGA2")
    ax.bar(x - 0.5 * width, spea2_data, width, label="SPEA2")
    ax.bar(x + 0.5 * width, fr_data, width, label="FR")
    ax.bar(x + 1.5 * width, ssga_fr_data, width, label="SSGA-FR")
    ax.bar(x + 2.5 * width, fr_nsga2_data, width, label="FR-NSGA2")
    ax.bar(x + 3.5 * width, fr_ssga_nsga2_data, width, label="FR-SSGA-NSGA2")

    # Add labels, title, and custom x-axis tick labels
    if i == 0:  # Only add y-label to the first subplot
        ax.set_ylabel(
            "Количество пересечений ребер (меньше - лучше)",
            fontsize=20,
            fontweight="bold",
        )
    ax.set_title(f"Графы с {postfix} вершинами", fontsize=20, fontweight="bold")
    ax.set_xticks(x)
    ax.set_xticklabels(["Планарные", "Произвольные"], fontsize=15, fontweight="bold")

    # Add a grid for better readability
    ax.yaxis.grid(True, linestyle="--", alpha=0.7)

# Place legend below the subplots
handles, labels = axes[0].get_legend_handles_labels()
fig.legend(
    handles,
    labels,
    loc="lower center",
    ncol=8,
    prop={"weight": "bold", "size": 16},
)  # Adjust ncol as needed

# Adjust layout to prevent labels from overlapping
fig.tight_layout(rect=[0, 0.05, 1, 0.95])  # Adjust rect to make space for the legend


assets_dir = "/Users/gregorykogan/Documents/МИФИ/6 семестр/УИР/ПЗ/assets"
plt.savefig(f"{assets_dir}/final-competition.png", dpi=300)
plt.show()
