# Add "zz_type" to fields that have a "type" key.
# "zz_type" contains the underlying type of "type" or, 
# if no underlying type is defined, "zz_type" = "type".

# underlying types
{
  "Int":      "int64",
  "Float":    "float64",
  "String":   "string",
  "BoolAV":   "bool",
  "FixQuality": "int64"
} as $types |

.items |= map(    
  .fields |= map(
    if has("type") then
      .type as $t |
      . + {"zz_type": ($types | if has($t) then .[$t] else $t end) }
    else
      .
    end
  )
)