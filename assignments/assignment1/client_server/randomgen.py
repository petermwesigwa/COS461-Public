#usr/bin/env python

from sys import argv, stderr

def main():
    if len(argv) != 1:
        print("Usage: randomgen [Integer]", file=stderr)