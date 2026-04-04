import sys

with open('internal/tools/tools.go', 'r') as f:
    lines = f.readlines()
    for i in range(208, 220):
        line = lines[i]
        print(f'Line {i+1}: {repr(line)}')
        if line.strip():
            first_char = line[0]
            print(f'  First char: {repr(first_char)} (code: {ord(first_char)})')