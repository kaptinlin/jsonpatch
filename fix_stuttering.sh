#!/bin/bash

# Script to fix stuttering type names in all operation files

operations=(
    "ends:EndsOperation"
    "test_string_len:TestStringLenOperation"  
    "merge:MergeOperation"
    "matches:MatchesOperation"
    "type:TypeOperation"
    "test:TestOperation"
    "in:InOperation"
    "remove:RemoveOperation"
    "or:OrOperation"
    "extend:ExtendOperation"
    "test_string:TestStringOperation"
    "strdel:StrDelOperation"
    "strins:StrInsOperation"
    "test_type:TestTypeOperation"
    "flip:FlipOperation"
    "move:MoveOperation"
    "starts:StartsOperation"
    "undefined:UndefinedOperation"
    "split:SplitOperation"
    "not:NotOperation"
    "more:MoreOperation"
    "inc:IncOperation"
    "replace:ReplaceOperation"
    "less:LessOperation"
)

for op_pair in "${operations[@]}"; do
    IFS=':' read -r filename new_name <<< "$op_pair"
    old_name="Op${new_name}"
    file_path="op/${filename}.go"
    
    if [[ -f "$file_path" ]]; then
        echo "Processing $file_path: $old_name -> $new_name"
        
        # Read file and check if it has the stuttering pattern
        if grep -q "type $old_name struct" "$file_path"; then
            # Create a temporary file with the fixes
            sed_script="
                # Replace the struct definition and alias
                s|// $old_name represents|// $new_name represents|g
                s|type $old_name struct|type $new_name struct|g
                s|// $new_name is a non-stuttering alias for $old_name\.|// $old_name is a backward-compatible alias for $new_name.|g
                s|type $new_name = $old_name|type $old_name = $new_name|g
                # Replace all function signatures and method receivers
                s|\*$old_name|\*$new_name|g
                s|&$old_name{|&$new_name{|g
            "
            
            # Apply the transformations
            sed -E "$sed_script" "$file_path" > "${file_path}.tmp" && mv "${file_path}.tmp" "$file_path"
        else
            echo "  No stuttering pattern found in $file_path"
        fi
    else
        echo "  File $file_path does not exist"
    fi
done

echo "Done processing all operations!"