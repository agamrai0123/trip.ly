"""
Renames all wanderplan/* module paths to github.com/agamrai0123/wanderplan/*.
Processes go.mod files and all .go source files under backend/.
"""
import os
import re

ROOT = r'D:\Learn\trip.ly\backend'

# Map of old module prefix -> new module prefix
REPLACEMENTS = [
    ('wanderplan/pkg/config',          'github.com/agamrai0123/wanderplan/pkg/config'),
    ('wanderplan/pkg/database',        'github.com/agamrai0123/wanderplan/pkg/database'),
    ('wanderplan/pkg/errors',          'github.com/agamrai0123/wanderplan/pkg/errors'),
    ('wanderplan/pkg/grpc',            'github.com/agamrai0123/wanderplan/pkg/grpc'),
    ('wanderplan/pkg/jwt',             'github.com/agamrai0123/wanderplan/pkg/jwt'),
    ('wanderplan/pkg/kafka',           'github.com/agamrai0123/wanderplan/pkg/kafka'),
    ('wanderplan/pkg/logger',          'github.com/agamrai0123/wanderplan/pkg/logger'),
    ('wanderplan/pkg/middleware',      'github.com/agamrai0123/wanderplan/pkg/middleware'),
    ('wanderplan/pkg/response',        'github.com/agamrai0123/wanderplan/pkg/response'),
    # proto gen
    ('wanderplan/proto/gen/wanderplan/v1', 'github.com/agamrai0123/wanderplan/proto/gen/wanderplan/v1'),
    ('wanderplan/proto/gen',           'github.com/agamrai0123/wanderplan/proto/gen'),
    ('wanderplan/proto',               'github.com/agamrai0123/wanderplan/proto'),
    # pkg umbrella (must come after sub-paths)
    ('wanderplan/pkg',                 'github.com/agamrai0123/wanderplan/pkg'),
    # services
    ('wanderplan/auth-service/internal', 'github.com/agamrai0123/wanderplan/services/auth-service/internal'),
    ('wanderplan/auth-service',        'github.com/agamrai0123/wanderplan/services/auth-service'),
    ('wanderplan/api-gateway/internal', 'github.com/agamrai0123/wanderplan/services/api-gateway/internal'),
    ('wanderplan/api-gateway',         'github.com/agamrai0123/wanderplan/services/api-gateway'),
    ('wanderplan/trip-service/internal', 'github.com/agamrai0123/wanderplan/services/trip-service/internal'),
    ('wanderplan/trip-service',        'github.com/agamrai0123/wanderplan/services/trip-service'),
    ('wanderplan/user-service/internal', 'github.com/agamrai0123/wanderplan/services/user-service/internal'),
    ('wanderplan/user-service',        'github.com/agamrai0123/wanderplan/services/user-service'),
    ('wanderplan/collaboration-service/internal', 'github.com/agamrai0123/wanderplan/services/collaboration-service/internal'),
    ('wanderplan/collaboration-service', 'github.com/agamrai0123/wanderplan/services/collaboration-service'),
    ('wanderplan/notification-service/internal', 'github.com/agamrai0123/wanderplan/services/notification-service/internal'),
    ('wanderplan/notification-service', 'github.com/agamrai0123/wanderplan/services/notification-service'),
    ('wanderplan/search-service/internal', 'github.com/agamrai0123/wanderplan/services/search-service/internal'),
    ('wanderplan/search-service',      'github.com/agamrai0123/wanderplan/services/search-service'),
]

def replace_in_file(path, replacements):
    with open(path, 'r', encoding='utf-8', errors='replace') as f:
        content = f.read()
    original = content
    for old, new in replacements:
        content = content.replace(old, new)
    if content != original:
        with open(path, 'w', encoding='utf-8') as f:
            f.write(content)
        return True
    return False

changed = 0
for dirpath, dirnames, filenames in os.walk(ROOT):
    # Skip .git and vendor directories
    dirnames[:] = [d for d in dirnames if d not in ('.git', 'vendor')]
    for filename in filenames:
        if filename.endswith('.go') or filename == 'go.mod' or filename == 'go.work':
            filepath = os.path.join(dirpath, filename)
            if replace_in_file(filepath, REPLACEMENTS):
                print(f'updated: {filepath[len(ROOT)+1:]}')
                changed += 1

print(f'\nTotal files updated: {changed}')
