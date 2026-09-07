def without_451:
  walk(
    if type == "object" and has("responses") then
      .responses |= del(.["451"])
    else .
    end
  );

.[0] as $standard_doc
| .[1] as $sandbox_doc
| ($standard_doc.paths | without_451) as $standard
| ($sandbox_doc.paths | without_451) as $sandbox
| [
    ($standard | keys[]) as $path
    | select($sandbox | has($path))
    | select($standard[$path] != $sandbox[$path])
    | $path
  ] as $divergent
| if ($divergent | length) > 0 then
    error(
      "common paths differ between standard.json and sandbox.json: "
      + ($divergent | join(", "))
    )
  else
    true
  end
