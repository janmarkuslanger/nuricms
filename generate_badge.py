import sys
import re
import subprocess

def generate_svg(coverage_percentage):
    color = "red" if coverage_percentage < 80 else "yellow" if coverage_percentage < 90 else "green"
    label = "coverage"
    label_width = 60
    value_width = 60
    total_width = label_width + value_width

    svg = f'''<svg xmlns="http://www.w3.org/2000/svg" width="{total_width}" height="20">
    <linearGradient id="smooth" x2="0" y2="100%">
      <stop offset="0" stop-color="#fff" stop-opacity=".7"/>
      <stop offset=".1" stop-color="#aaa" stop-opacity=".1"/>
      <stop offset=".9" stop-color="#000" stop-opacity=".3"/>
      <stop offset="1" stop-color="#000" stop-opacity=".5"/>
    </linearGradient>
    <rect rx="3" width="{total_width}" height="20" fill="#555"/>
    <rect rx="3" x="{label_width}" width="{value_width}" height="20" fill="{color}"/>
    <rect rx="3" width="{total_width}" height="20" fill="url(#smooth)"/>
    <g fill="#fff" text-anchor="middle" font-family="Verdana,Geneva,DejaVu Sans,sans-serif" font-size="11">
      <text x="{label_width/2}" y="15" fill="#010101" fill-opacity=".3">{label}</text>
      <text x="{label_width/2}" y="14">{label}</text>
      <text x="{label_width + value_width/2}" y="15" fill="#010101" fill-opacity=".3">{coverage_percentage:.1f}%</text>
      <text x="{label_width + value_width/2}" y="14">{coverage_percentage:.1f}%</text>
    </g>
    </svg>'''
    
    return svg

def get_coverage_percentage():
    result = subprocess.run(
        ["go", "tool", "cover", "-func=coverage.out"],
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        text=True
    )
    
    print(result.stdout)

    match = re.search(r"total:.*?([\d\.]+)%", result.stdout)
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
