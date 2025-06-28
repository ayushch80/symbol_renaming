#!/usr/bin/env python3
"""
rename_symbols.py

A universal Python script that takes:
  1) a binary or text file (e.g. a Mach-O, ELF, object dump, or plain text),
  2) a list of symbol names to rename (one per line),
and produces a patched output file where each listed symbol is replaced
by a newly generated name of the same length.

Usage:
  python rename_symbols.py \
      --input-file input.bin \
      --output-file patched.bin \
      --symbols-file symbols.txt

Output:
  - patched.bin        # The file with symbols renamed in-place
  - mapping.json       # JSON file with the mapping of original → new names
"""

import argparse
import os
import sys
import json
import random
import string

def parse_args():
    p = argparse.ArgumentParser(
        description="Rename selected symbols in a file, in-place, using same-length random names."
    )
    p.add_argument(
        "-i", "--input-file", required=True,
        help="Path to the input file (binary or text)."
    )
    p.add_argument(
        "-o", "--output-file", required=True,
        help="Path where the patched file will be written."
    )
    p.add_argument(
        "-s", "--symbols-file", required=True,
        help="Path to a text file containing one symbol name per line."
    )
    return p.parse_args()

def load_symbols(path):
    """Load symbols to rename from a text file, one per line."""
    with open(path, "r", encoding="utf-8") as f:
        syms = [line.strip() for line in f if line.strip()]
    # Remove duplicates while preserving order
    seen = set()
    unique = []
    for s in syms:
        if s not in seen:
            seen.add(s)
            unique.append(s)
    return unique

def random_name(length):
    """Generate a random name of exactly `length` characters.
       Must start with letter or underscore (to be valid as a symbol)."""
    if length < 1:
        raise ValueError("Symbol length must be >= 1")
    first_chars = string.ascii_letters + "_"
    other_chars = string.ascii_letters + string.digits + "_"
    name = [random.choice(first_chars)]
    for _ in range(length-1):
        name.append(random.choice(other_chars))
    return "".join(name)

def build_mapping(symbols):
    """Build a mapping original -> new, preserving exact length."""
    mapping = {}
    for sym in symbols:
        L = len(sym)
        new = random_name(L)
        # Avoid collisions
        while new in mapping.values() or new == sym:
            new = random_name(L)
        mapping[sym] = new
    return mapping

def patch_file(in_path, out_path, mapping):
    """Read the input file as binary, replace occurrences, write patched output."""
    data = open(in_path, "rb").read()
    for orig, new in mapping.items():
        ob = orig.encode("utf-8")
        nb = new.encode("utf-8")
        if len(ob) != len(nb):
            raise RuntimeError(f"Length mismatch for '{orig}': {len(ob)} vs {len(nb)}")
        count = data.count(ob)
        if count:
            print(f"  Replacing {count} occurrence(s) of '{orig}' → '{new}'")
            data = data.replace(ob, nb)
    with open(out_path, "wb") as f:
        f.write(data)

def main():
    args = parse_args()
    if not os.path.isfile(args.input_file):
        print(f"Error: input file not found: {args.input_file}", file=sys.stderr)
        sys.exit(1)
    if not os.path.isfile(args.symbols_file):
        print(f"Error: symbols file not found: {args.symbols_file}", file=sys.stderr)
        sys.exit(1)

    symbols = load_symbols(args.symbols_file)
    if not symbols:
        print("No symbols to rename.", file=sys.stderr)
        sys.exit(1)

    mapping = build_mapping(symbols)

    print(f"Patching {len(symbols)} symbols in '{args.input_file}' → '{args.output_file}'")
    patch_file(args.input_file, args.output_file, mapping)

    # Save mapping for later reference
    with open("mapping.json", "w", encoding="utf-8") as m:
        json.dump(mapping, m, indent=2)
    print("Mapping written to mapping.json")

if __name__ == "__main__":
    main()