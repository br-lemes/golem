.components.schemas
| to_entries
| map(select(.value.enum? and (.value.enum | all(.[]; type == "string"))))
| map({key: .key, value: .value.enum})
| from_entries
