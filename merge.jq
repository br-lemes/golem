.[0] * .[1]
| walk(
    if type == "object" and has("responses") then
      .responses |= del(.["451"])
    else .
    end
  )
