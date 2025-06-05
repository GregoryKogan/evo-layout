import json
import os
import matplotlib.pyplot as plt
from matplotlib.ticker import ScalarFormatter


def plot_intersections_over_time(log_dir):
    # Setup plot
    plt.figure(figsize=(12, 8))
    colors = plt.cm.tab10.colors

    CUTOFF_TIME = 10

    # Plot each algorithm's data
    files = [f for f in os.listdir(log_dir) if f.endswith(".jsonl")]
    for i, filename in enumerate(files):
        algo_name = filename.split(".")[0].split("_")[1]
        filepath = os.path.join(log_dir, filename)

        times, intersections = [], []
        with open(filepath) as f:
            for line in f:
                try:
                    data = json.loads(line)
                    if "elapsed" in data and "solution" in data:
                        time_sec = data["elapsed"] / 1e9
                        inter = data["solution"].get("intersections", float("inf"))

                        # Stop collecting data after cutoff time
                        if time_sec > CUTOFF_TIME:
                            break

                        times.append(time_sec)
                        intersections.append(inter)
                except json.JSONDecodeError:
                    continue

        # Plot algorithm data
        plt.plot(
            times,
            intersections,
            "o-",
            markersize=5,
            linewidth=1.5,
            label=algo_name,
            color=colors[i % len(colors)],
            alpha=0.8,
        )

    # Configure plot
    plt.xlabel("Время (секунды)", fontsize=12)
    plt.ylabel("Количество пересечений рёбер", fontsize=12)
    plt.grid(True, alpha=0.2)
    plt.ylim(bottom=0)  # Start y-axis at 0

    # Format x-axis for large time values
    if CUTOFF_TIME > 1000:
        plt.ticklabel_format(axis="x", style="sci", scilimits=(0, 0))
        plt.gca().xaxis.get_offset_text().set_fontsize(10)

    plt.legend(title="Алгоритмы", fontsize=10)
    plt.tight_layout()
    assets_dir = "/Users/gregorykogan/Documents/МИФИ/6 семестр/УИР/ПЗ/assets"
    plt.savefig(f"{assets_dir}/g-25-comp.png", dpi=300)
    plt.show()


if __name__ == "__main__":
    log_directory = "./visual/competition/logs"
    plot_intersections_over_time(log_directory)
