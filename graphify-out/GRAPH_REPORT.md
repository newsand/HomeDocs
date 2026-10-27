# Graph Report - /home/fils/Documents/pessoais/ID-Center  (2026-10-27)

## Corpus Check
- 15 files · ~27,228 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 91 nodes · 198 edges · 9 communities detected
- Extraction: 62% EXTRACTED · 38% INFERRED · 0% AMBIGUOUS · INFERRED: 75 edges (avg confidence: 0.8)
- Token cost: 0 input · 0 output

## Community Hubs (Navigation)
- [[_COMMUNITY_Community 0|Community 0]]
- [[_COMMUNITY_Community 1|Community 1]]
- [[_COMMUNITY_Community 2|Community 2]]
- [[_COMMUNITY_Community 3|Community 3]]
- [[_COMMUNITY_Community 4|Community 4]]
- [[_COMMUNITY_Community 5|Community 5]]
- [[_COMMUNITY_Community 6|Community 6]]
- [[_COMMUNITY_Community 7|Community 7]]
- [[_COMMUNITY_Community 8|Community 8]]

## God Nodes (most connected - your core abstractions)
1. `ShowPerson()` - 14 edges
2. `CreatePerson()` - 14 edges
3. `app` - 13 edges
4. `run()` - 12 edges
5. `SearchPeople()` - 10 edges
6. `UpdatePerson()` - 10 edges
7. `openTestDB()` - 8 edges
8. `columnNames()` - 7 edges
9. `TestCreatePerson()` - 7 edges
10. `TestDeletePerson()` - 7 edges

## Surprising Connections (you probably didn't know these)
- `run()` --calls--> `Migrate()`  [INFERRED]
  /home/fils/Documents/pessoais/ID-Center/cmd/idcenter/main.go → /home/fils/Documents/pessoais/ID-Center/internal/person/db.go
- `run()` --calls--> `New()`  [INFERRED]
  /home/fils/Documents/pessoais/ID-Center/cmd/idcenter/main.go → /home/fils/Documents/pessoais/ID-Center/internal/http/server.go
- `run()` --calls--> `TestUpdatePerson()`  [INFERRED]
  /home/fils/Documents/pessoais/ID-Center/cmd/idcenter/main.go → /home/fils/Documents/pessoais/ID-Center/internal/person/person_test.go
- `run()` --calls--> `TestDeletePerson()`  [INFERRED]
  /home/fils/Documents/pessoais/ID-Center/cmd/idcenter/main.go → /home/fils/Documents/pessoais/ID-Center/internal/person/person_test.go
- `run()` --calls--> `TestPages()`  [INFERRED]
  /home/fils/Documents/pessoais/ID-Center/cmd/idcenter/main.go → /home/fils/Documents/pessoais/ID-Center/internal/http/server_test.go

## Communities

### Community 0 - "Community 0"
Cohesion: 0.26
Nodes (12): CreatePerson(), replaceScans(), extensionOf(), FilePath(), newID(), removeNames(), removeStored(), writeUpload() (+4 more)

### Community 1 - "Community 1"
Cohesion: 0.42
Nodes (2): groupsFrom(), app

### Community 2 - "Community 2"
Cohesion: 0.36
Nodes (9): Open(), main(), run(), names(), openTestDB(), TestCreatePerson(), TestSearchPeople(), TestShowPerson() (+1 more)

### Community 3 - "Community 3"
Cohesion: 0.33
Nodes (7): DeletePerson(), TestDeletePerson(), likeNeedle(), SearchPeople(), loadScans(), scanPerson(), ShowPerson()

### Community 4 - "Community 4"
Cohesion: 0.22
Nodes (8): CNH, CPF, Input, Passport, Person, RG, Upload, VoterID

### Community 5 - "Community 5"
Cohesion: 0.36
Nodes (7): columnNames(), columnValues(), scanDest(), searchColumns(), insertPerson(), Migrate(), updatePersonRow()

### Community 6 - "Community 6"
Cohesion: 0.29
Nodes (7): catalog(), text(), formField, formGroup, ListItem, ViewField, ViewGroup

### Community 7 - "Community 7"
Cohesion: 0.43
Nodes (6): TestCatalogKeepsCPFOnce(), filePaths(), multipartBody(), newRouter(), serve(), TestPages()

### Community 8 - "Community 8"
Cohesion: 0.38
Nodes (6): New(), personFromInput(), readFiles(), readParts(), requestID(), requestLog()

## Knowledge Gaps
- **13 isolated node(s):** `Upload`, `CNH`, `CPF`, `RG`, `VoterID` (+8 more)
  These have ≤1 connection - possible missing edges or undocumented components.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `CreatePerson()` connect `Community 0` to `Community 1`, `Community 2`, `Community 3`, `Community 5`?**
  _High betweenness centrality (0.254) - this node is a cross-community bridge._
- **Why does `ValidName()` connect `Community 0` to `Community 2`, `Community 4`?**
  _High betweenness centrality (0.191) - this node is a cross-community bridge._
- **Why does `ShowPerson()` connect `Community 3` to `Community 0`, `Community 1`, `Community 2`, `Community 5`?**
  _High betweenness centrality (0.190) - this node is a cross-community bridge._
- **Are the 11 inferred relationships involving `ShowPerson()` (e.g. with `DeletePerson()` and `columnNames()`) actually correct?**
  _`ShowPerson()` has 11 INFERRED edges - model-reasoned connections that need verification._
- **Are the 13 inferred relationships involving `CreatePerson()` (e.g. with `ValidName()` and `writeUpload()`) actually correct?**
  _`CreatePerson()` has 13 INFERRED edges - model-reasoned connections that need verification._
- **Are the 10 inferred relationships involving `run()` (e.g. with `Open()` and `Migrate()`) actually correct?**
  _`run()` has 10 INFERRED edges - model-reasoned connections that need verification._
- **Are the 8 inferred relationships involving `SearchPeople()` (e.g. with `columnNames()` and `searchColumns()`) actually correct?**
  _`SearchPeople()` has 8 INFERRED edges - model-reasoned connections that need verification._