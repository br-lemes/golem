(.components.schemas) as $allSchemas
| ($allSchemas
   | to_entries
   | map(select(.value.type == "string") | .key)) as $stringSchemas
| {
  openapi: .openapi,
  info: .info,
  paths: (
    .paths
    | walk(
        if type == "object" and .["$ref"]? then
          (.["$ref"] | split("/") | last) as $name
          | if ($stringSchemas | index($name)) != null
            then {type: "string"}
            else .
            end
        else .
        end
      )
    | with_entries(
        .value |= with_entries(
          select(.key == "parameters" or .key == "get")
          | if .key == "get" then
              .value |= (del(.responses, .requestBody, .security) | .parameters = (.parameters // []))
            else .
            end
        )
        | select(.value | has("get"))
      )
  ),
  components: {schemas: {}}
}
