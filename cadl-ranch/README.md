## @azure-tools/typespec-go

### step

- add ./script/tspcompile.js cadlRanch list(cadl-ranch-specs resourcemanager path)
```
  'managed_identity': ['azure/resource-manager/models/common-types/managed-identity'],
  'resources': ['azure/resource-manager/models/resources'],
```
- `node ./scripts/tspcompile.js --filter` (previous step added cadlRanch)

- `node ./scripts/cadl-ranch.js --start/--stop` cadl-ranch server

- write cadl-ranch test, ref: `https://github.com/Azure/cadl-ranch/tree/main/packages/cadl-ranch-specs/http/azure/resource-manager/models`