# Add "zz_i" containing the index to all fields.
# Add "zz_xarg" to all fields with a base name (no dot in the name).
# "zz_xarg" contains an array of sub field indices (same base name with dot postfix)
.items |= map(
  .fields |= (
    # add zz_i
    reduce range(length) as $zz_i (.; .[$zz_i] += {$zz_i})
    # group by base name
    | group_by(.name | split(".")[0])
    # sort each sub-array by base name then by zz_i order (note: sort_by arg is an implicit array)
    | (.[] |= sort_by((.name | contains(".")), .zz_i))
    # add zz_xarg to each base name field containing an array of sub field indices
    | (.[] |= .[0] + {"zz_xarg": .[1:] | map(.zz_i)})
    # keep fields ordered
    | sort_by(.zz_i)
  )
)