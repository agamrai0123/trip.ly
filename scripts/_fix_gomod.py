"""Remove workspace module require entries from service go.mod files."""
import re
import os

ROOT = r'D:\Learn\trip.ly\backend'

service_dirs = [
    'services/auth-service',
    'services/api-gateway',
    'services/trip-service',
    'services/user-service',
    'services/collaboration-service',
    'services/notification-service',
    'services/search-service',
]

workspace_modules = [
    'github.com/agamrai0123/wanderplan/pkg',
    'github.com/agamrai0123/wanderplan/proto',
]

for rel in service_dirs:
    gomod = os.path.join(ROOT, rel, 'go.mod')
    if not os.path.exists(gomod):
        print('skip:', rel)
        continue
    with open(gomod, 'r', encoding='utf-8') as f:
        content = f.read()
    original = content
    for mod in workspace_modules:
        content = re.sub(r'\t' + re.escape(mod) + r'[^\n]*\n', '', content)
    if content != original:
        with open(gomod, 'w', encoding='utf-8') as f:
            f.write(content)
        print('updated:', rel)
    else:
        print('no change:', rel)
