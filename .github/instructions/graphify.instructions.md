---
applyTo: "**"
---

# Prefer the knowledge graph over grep

This repo has a `graphify` code graph. To find a definition, trace a call, or
scope a change, query it before grepping or reading blind — it's faster and
cheaper. Fall back to direct reads for a function's actual body or non-code text.

```bash
G="$(go env GOPATH)/bin/graphify"
$G query <regex>      # where a name is defined (file:line)
$G explain <node>     # a node + what calls it / it calls / imports
$G affected <file>    # a change's blast radius
$G update .           # refresh the graph (incremental)
```

If `graphify-out/graph.json` is missing, build it first with `$G build .`.
