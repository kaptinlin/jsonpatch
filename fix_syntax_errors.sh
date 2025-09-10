#!/bin/bash

# Script to fix syntax errors introduced by the previous script

cd op

# Fix composite literal syntax errors
for file in *.go; do
    if [[ -f "$file" ]]; then
        # Fix malformed composite literals
        sed -i '' -E 's/&Op([A-Za-z]+Operation)\{([A-Za-z]+Operation)\{/\&\2\{/g' "$file"
        sed -i '' -E 's/&Op([A-Za-z]+Operation)\{([A-Za-z]+)\{/\&\2\{/g' "$file"
    fi
done

echo "Fixed syntax errors in operation files"