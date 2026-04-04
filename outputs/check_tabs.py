import sys

with open('internal/tools/tools.go', 'r') as f:
    lines = f.readlines()
    for i, line in enumerate(lines[212:214], 213):
        print(f'Line {i}: {repr(line)}')
        if 'rel, _ := filepath.Rel' in line:
            print(f'  First char code: {ord(line[0]) if line else "empty"}')
            print(f'  Line starts with tab: {line.startswith(chr(9))}')
            print(f'  Line starts with space: {line.startswith(" ")}')
            # Show the exact string
            print(f'  Exact line: |{line}|')