#!/usr/bin/env bash
set -e

echo "========== BUILDING & INSTALLING OUTPUTCOPY =========="
go install .

GO_BIN="$(go env GOPATH)/bin"
SHELL_RC="${HOME}/.bashrc"

if [ -f "${HOME}/.zshrc" ] && [ -n "$ZSH_VERSION" ]; then
    SHELL_RC="${HOME}/.zshrc"
fi

if ! grep -q "$GO_BIN" "$SHELL_RC" 2>/dev/null; then
    echo "" >> "$SHELL_RC"
    echo "# Added by outputcopy installer" >> "$SHELL_RC"
    echo "export PATH=\"\$PATH:$GO_BIN\"" >> "$SHELL_RC"
    echo "Configured PATH in $SHELL_RC"
else
    echo "PATH already configured in $SHELL_RC"
fi

echo "========== INSTALLATION COMPLETE =========="
