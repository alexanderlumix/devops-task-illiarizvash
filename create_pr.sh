#!/bin/bash

# Create Pull Request Script
echo "🚀 Creating Pull Request..."

# Repository details
REPO="alexanderlumix/devops-task-illiarizvash"
BASE_BRANCH="main"
HEAD_BRANCH="dev"

# PR details
TITLE="Translate project documentation to English"
BODY="## Summary

This PR translates all Russian documentation to English and renames files with Russian names.

### Changes Made:
- ✅ Translated all documentation files from Russian to English
- ✅ Renamed \`ОТЧЕТ_ПО_ВЫПОЛНЕННОЙ_РАБОТЕ_И_НАСТРОЙКАМ.md\` to \`TASK_RESULT_REPORT.md\`
- ✅ Updated all README files and documentation
- ✅ Maintained project structure and formatting

### Files Updated:
- \`local-development/README.md\` - Local development documentation
- \`local-development/SCRIPTS_SUMMARY.md\` - Scripts summary
- \`problem-solving/action-items.md\` - Action items for critical issues
- \`problem-solving/final-status-report.md\` - Final status report
- \`problem-solving/quick-checklist.md\` - Quick checklist
- \`problem-solving/ISSUES_SUMMARY.md\` - Issues summary
- \`TASK_RESULT_REPORT.md\` - Task result report (renamed)

### Benefits:
- 🌍 International accessibility
- 📚 Better documentation for global teams
- 🔍 Easier code review process
- 📖 Consistent English documentation

All functionality remains unchanged, only documentation language updated."

# Create PR using GitHub CLI
echo "📝 Creating PR with title: $TITLE"
echo "🔗 Repository: $REPO"
echo "📤 From: $HEAD_BRANCH"
echo "📥 To: $BASE_BRANCH"

# Try to create PR
if command -v gh &> /dev/null; then
    echo "🔧 Using GitHub CLI..."
    gh pr create \
        --repo "$REPO" \
        --title "$TITLE" \
        --body "$BODY" \
        --base "$BASE_BRANCH" \
        --head "$HEAD_BRANCH" \
        --web
else
    echo "❌ GitHub CLI not available"
    echo "🔗 Please create PR manually using this URL:"
    echo "https://github.com/$REPO/compare/$BASE_BRANCH...$HEAD_BRANCH"
fi

echo "✅ PR creation script completed!" 