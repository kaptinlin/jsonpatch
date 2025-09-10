#!/bin/bash

# Add non-stuttering type aliases for all Op*Operation types

# Array of files and their corresponding operations
declare -A ops=(
    ["op/ends.go"]="EndsOperation"
    ["op/extend.go"]="ExtendOperation" 
    ["op/flip.go"]="FlipOperation"
    ["op/in.go"]="InOperation"
    ["op/inc.go"]="IncOperation"
    ["op/less.go"]="LessOperation"
    ["op/matches.go"]="MatchesOperation"
    ["op/merge.go"]="MergeOperation"
    ["op/more.go"]="MoreOperation"
    ["op/move.go"]="MoveOperation"
    ["op/not.go"]="NotOperation"
    ["op/or.go"]="OrOperation"
    ["op/remove.go"]="RemoveOperation"
    ["op/replace.go"]="ReplaceOperation"
    ["op/split.go"]="SplitOperation"
    ["op/starts.go"]="StartsOperation"
    ["op/strdel.go"]="StrDelOperation"
    ["op/strins.go"]="StrInsOperation"
    ["op/test.go"]="TestOperation"
    ["op/test_string.go"]="TestStringOperation"
    ["op/test_string_len.go"]="TestStringLenOperation"
    ["op/test_type.go"]="TestTypeOperation"
    ["op/type.go"]="TypeOperation"
    ["op/undefined.go"]="UndefinedOperation"
)

for file in "${!ops[@]}"; do
    alias_name="${ops[$file]}"
    op_name="Op${alias_name}"
    
    echo "Processing $file -> $alias_name"
    
    # Find the line with the struct definition and add the alias after it
    if grep -q "type $op_name struct" "$file"; then
        # Find the closing brace of the struct and add the alias after it
        awk -v alias="$alias_name" -v op="$op_name" '
            /^type '"$op_name"' struct/ { in_struct = 1; print; next }
            in_struct && /^}$/ { 
                print
                print ""
                print "// " alias " is a non-stuttering alias for " op "."
                print "type " alias " = " op
                in_struct = 0
                next
            }
            { print }
        ' "$file" > "$file.tmp" && mv "$file.tmp" "$file"
    fi
done

# Handle interfaces.go - add alias for OpResult
echo "Processing op/interfaces.go -> Result"
file="op/interfaces.go"
if grep -q "type OpResult" "$file"; then
    awk '
        /^type OpResult/ { 
            print
            getline
            print
            print "// Result is a non-stuttering alias for OpResult."
            print "type Result[T internal.Document] = OpResult[T]"
            next
        }
        { print }
    ' "$file" > "$file.tmp" && mv "$file.tmp" "$file"
fi

echo "Done!"