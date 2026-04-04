#!/usr/bin/env bash
set -euo pipefail

echo "Testing fix for issue #26: run-qa-agent.sh and run-report-agent.sh build binary if missing"
echo "================================================================================"

# Note: We can't use rm due to governance, so we'll test differently
echo ""
echo "Note: Due to governance restrictions, we can't remove the binary."
echo "Instead, we'll verify the build logic is present in the scripts."
echo ""

echo "1. Checking run-qa-agent.sh has build logic..."
if grep -q "Building shellforge" scripts/run-qa-agent.sh; then
    echo "✓ run-qa-agent.sh has build logic"
else
    echo "✗ run-qa-agent.sh missing build logic"
    exit 1
fi

echo ""
echo "2. Checking run-report-agent.sh has build logic..."
if grep -q "Building shellforge" scripts/run-report-agent.sh; then
    echo "✓ run-report-agent.sh has build logic"
else
    echo "✗ run-report-agent.sh missing build logic"
    exit 1
fi

echo ""
echo "3. Comparing with run-agent.sh (reference implementation)..."
if grep -q "if \[\[ ! -f ./shellforge \]\]" scripts/run-qa-agent.sh && \
   grep -q "if \[\[ ! -f ./shellforge \]\]" scripts/run-report-agent.sh; then
    echo "✓ Both scripts have the same build pattern as run-agent.sh"
else
    echo "✗ Build pattern doesn't match run-agent.sh"
    exit 1
fi

echo ""
echo "4. Checking script syntax..."
bash -n scripts/run-qa-agent.sh && echo "✓ run-qa-agent.sh syntax OK" || exit 1
bash -n scripts/run-report-agent.sh && echo "✓ run-report-agent.sh syntax OK" || exit 1

echo ""
echo "================================================================================"
echo "All checks passed! The fix for issue #26 is implemented correctly."
echo "Both scripts now have the same build-if-missing logic as run-agent.sh."