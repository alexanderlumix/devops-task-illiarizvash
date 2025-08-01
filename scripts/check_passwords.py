#!/usr/bin/env python3
"""
Password Detection Script for Pre-commit Hook
This script scans files for potential hardcoded passwords and credentials.
"""

import re
import sys
import os
from pathlib import Path
from typing import List, Dict, Tuple

# Patterns for detecting passwords and credentials
PASSWORD_PATTERNS = [
    # Common password patterns
    r'password\s*[:=]\s*["\']([^"\']{6,})["\']',
    r'pass\s*[:=]\s*["\']([^"\']{6,})["\']',
    r'pwd\s*[:=]\s*["\']([^"\']{6,})["\']',
    r'secret\s*[:=]\s*["\']([^"\']{6,})["\']',
    r'token\s*[:=]\s*["\']([^"\']{6,})["\']',
    r'key\s*[:=]\s*["\']([^"\']{6,})["\']',
    r'auth\s*[:=]\s*["\']([^"\']{6,})["\']',
    
    # MongoDB connection strings with passwords
    r'mongodb://[^:]+:([^@]+)@',
    r'mongodb\+srv://[^:]+:([^@]+)@',
    
    # Hardcoded credentials
    r'["\'](admin|root|user|test|demo|password|secret|123456|qwerty|admin123)["\']',
    
    # API keys and tokens
    r'["\']([a-zA-Z0-9]{32,})["\']',  # Long alphanumeric strings
    r'["\']([a-zA-Z0-9]{20,})["\']',  # Medium alphanumeric strings
    
    # Common weak passwords
    r'["\'](password|123456|qwerty|admin|root|user|test|demo|guest|welcome)["\']',
]

# Files to exclude from scanning
EXCLUDED_FILES = {
    '.secrets.baseline',
    '.env.example',
    'passwords_and_credentials.txt',
    'README.md',
    'node_modules',
    '.git',
    '__pycache__',
    '.pytest_cache',
    '.venv',
    'venv',
}

# Directories to exclude
EXCLUDED_DIRS = {
    'docs',
    'tests',
    '.git',
    'node_modules',
    '__pycache__',
    '.pytest_cache',
    '.venv',
    'venv',
}

def is_excluded_file(filepath: str) -> bool:
    """Check if file should be excluded from scanning."""
    path = Path(filepath)
    
    # Check if file is in excluded directories
    for part in path.parts:
        if part in EXCLUDED_DIRS:
            return True
    
    # Check if file is in excluded files
    if path.name in EXCLUDED_FILES:
        return True
    
    # Check if file is documentation
    if path.suffix == '.md' and 'docs' in path.parts:
        return True
    
    return False

def scan_file_for_passwords(filepath: str) -> List[Tuple[int, str, str]]:
    """Scan a single file for password patterns."""
    issues = []
    
    try:
        with open(filepath, 'r', encoding='utf-8') as f:
            content = f.read()
            
        lines = content.split('\n')
        
        for line_num, line in enumerate(lines, 1):
            for pattern in PASSWORD_PATTERNS:
                matches = re.finditer(pattern, line, re.IGNORECASE)
                for match in matches:
                    # Extract the matched password/credential
                    if match.groups():
                        credential = match.group(1)
                    else:
                        credential = match.group(0)
                    
                    # Skip if it's a variable name or comment
                    if '=' in line and credential in line.split('=')[0]:
                        continue
                    
                    # Skip if it's a variable reference (like {ADMIN_PASS})
                    if credential.startswith('{') and credential.endswith('}'):
                        continue
                    
                    if line.strip().startswith('#'):
                        continue
                    
                    issues.append((line_num, line.strip(), credential))
                    
    except Exception as e:
        print(f"Error reading file {filepath}: {e}", file=sys.stderr)
    
    return issues

def main():
    """Main function to scan files for passwords."""
    issues_found = []
    
    # Get list of files to scan from command line arguments
    files_to_scan = sys.argv[1:] if len(sys.argv) > 1 else []
    
    if not files_to_scan:
        print("No files provided to scan", file=sys.stderr)
        return 1
    
    for filepath in files_to_scan:
        if not os.path.exists(filepath):
            print(f"File not found: {filepath}", file=sys.stderr)
            continue
            
        if is_excluded_file(filepath):
            continue
            
        issues = scan_file_for_passwords(filepath)
        if issues:
            issues_found.append((filepath, issues))
    
    # Report findings
    if issues_found:
        print("🚨 PASSWORD DETECTION ALERT 🚨")
        print("=" * 50)
        print("Potential hardcoded passwords/credentials found:")
        print()
        
        for filepath, issues in issues_found:
            print(f"📁 File: {filepath}")
            for line_num, line, credential in issues:
                print(f"   Line {line_num}: {line}")
                print(f"   ⚠️  Potential credential: {credential}")
                print()
        
        print("=" * 50)
        print("❌ COMMIT BLOCKED: Hardcoded credentials detected!")
        print()
        print("💡 Recommendations:")
        print("   - Use environment variables for sensitive data")
        print("   - Create .env files for local development")
        print("   - Use secret management services in production")
        print("   - Add .env to .gitignore")
        print()
        print("Example:")
        print("   ❌ password = 'secret123'")
        print("   ✅ password = os.getenv('DB_PASSWORD')")
        
        return 1
    else:
        print("✅ No hardcoded passwords detected")
        return 0

if __name__ == "__main__":
    sys.exit(main()) 