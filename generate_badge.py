import sys
import re
import subprocess

def generate_svg(coverage_percentage):
    color = "red" if coverage_percentage < 80 else "yellow" if coverage_percentage < 90 else "green"
    
    svg = f'''<svg xmlns="http://www.w3.org/2000/svg" width="120" height="20">
    <rect width="120" height="20" rx="3" fill="#555" />
    <rect width="{coverage_percentage * 1.2}" height="20" rx="3" fill="{color}" />
    <text x="60" y="15" font-size="11" fill="#fff" text-anchor="middle">{coverage_percentage}%</text>
    </svg>'''
    
    return svg

def get_coverage_percentage():
    result = subprocess.run(
        ["go", "tool", "cover", "-func=coverage.out"],
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        text=True
    )

    match = re.search(r"total:\s+(\d+\.\d+)%", result.stdout)
    if match:
        return float(match.group(1))
    else:
        raise ValueError("Coverage percentage not found in go tool cover output")

if __name__ == "__main__":
    try:
        coverage_percentage = get_coverage_percentage()
        svg = generate_svg(coverage_percentage)
        with open("coverage_badge.svg", "w") as f:
            f.write(svg)
        print("Badge generated successfully!")
    except Exception as e:
        print(f"Error: {e}")
        sys.exit(1)
