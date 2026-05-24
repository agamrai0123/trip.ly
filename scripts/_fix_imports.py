"""Fix corrupted import paths caused by duplicate rename script runs.

Any path like:
  github.com/agamrai0123/[github.com/agamrai0123/[...]]wanderplan/X
should become:
  github.com/agamrai0123/wanderplan/X

Also ensures proto module name and service go.mod require entries are correct.
"""

import re
import os

BACKEND = r'D:\Learn\trip.ly\backend'

# Regex to collapse any repeated github.com/agamrai0123/ prefix
# Matches: github.com/agamrai0123/(github.com/agamrai0123/)*wanderplan/
COLLAPSE_PATTERN = re.compile(
    r'github\.com/agamrai0123/(?:github\.com/agamrai0123/)*wanderplan/'
)
CORRECT_PREFIX = 'github.com/agamrai0123/wanderplan/'

def fix_content(content: str) -> str:
    return COLLAPSE_PATTERN.sub(CORRECT_PREFIX, content)

updated = []
skipped = []

for root, dirs, files in os.walk(BACKEND):
    # Skip hidden dirs and vendor
    dirs[:] = [d for d in dirs if not d.startswith('.') and d != 'vendor']
    for fname in files:
        if not (fname.endswith('.go') or fname == 'go.mod' or fname == 'go.work'):
            continue
        fpath = os.path.join(root, fname)
        try:
            with open(fpath, 'r', encoding='utf-8') as f:
                content = f.read()
        except Exception as e:
            print(f'  ERROR reading {fpath}: {e}')
            continue
        fixed = fix_content(content)
        if fixed != content:
            with open(fpath, 'w', encoding='utf-8') as f:
                f.write(fixed)
            rel = os.path.relpath(fpath, BACKEND)
            updated.append(rel)
        else:
            skipped.append(os.path.relpath(fpath, BACKEND))

print(f'\nUpdated {len(updated)} files:')
for f in updated:
    print(f'  {f}')
print(f'\nNo change needed in {len(skipped)} files')
