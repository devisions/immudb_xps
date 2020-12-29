
## Notes

Used this simple example to test the behaviour of `getByIndex` when storing bytes ([]byte) of a custom type encoded using gob.

Result:

- immuclient fails:
  ```
  $ immuclient getByIndex 6
    proto: cannot parse reserved wire type
  $
  ```
- `si, _ := client.ByIndex(ctx, 6)` returns `nil`
