.components.schemas as $allSchemas
| ($allSchemas
   | to_entries
   | map(select(.value.type == "string") | .key)) as $stringSchemas
| {
    openapi: .openapi,
    info: .info,
    paths: {},
    components: {
      schemas: (
        $allSchemas
        | del(.ValidationError, .HTTPValidationError)
        | walk(
            if type == "object" and .["$ref"]? then
              (.["$ref"] | split("/") | last) as $name
              | if ($stringSchemas | index($name)) != null
                then $allSchemas[$name]
                else .
                end
            else .
            end
          )
        | with_entries(select(.key as $name | ($stringSchemas | index($name)) == null))
        | walk(if type == "object" then del(.enum) else . end)
      )
    }
  }
