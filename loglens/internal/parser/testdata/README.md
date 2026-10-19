# Parser fixtures

These records are entirely synthetic. Field names reflect a small local OpenClaw sample; they are not a complete upstream schema contract.

`valid.jsonl` contains ordinary events. `edge-cases.jsonl` deliberately mixes malformed JSON, wrong types, missing fields, invalid timestamps, and unknown severity. It is not a JSON document; read it by line.

Define expected outcomes before writing assertions. Generate partial writes, unterminated records, and oversized input in reader tests later rather than treating them as parser responsibilities. Never replace these fixtures with real private logs.

## References

[The art of writing parser](https://arielortiz.info/pycon2025/#_create_and_implement_parsing_functions_for_each_nonterminal)
