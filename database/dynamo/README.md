## Usage

### Initiate the client
```
client := CreateNewQupDynamo( dynamodb.New(sess) )
```
#### Supported tags
| Tag    | description                                       |
|--------|---------------------------------------------------|
| key    | Primary key of the record                         |
| idx    | Indication the field belongs to a SI              |
| unique | Indicates the field should be unique in the table |
| range  | Indicates the field is a range key                |

#### Retrieve a single item
```
type ExampleRecord struct {
	Id         string `dynamo:"id,key"`
	ExternalId string `dynamo:"external_id,idx|gsi-name"`
	GroupKey   string `dynamo:"group_key,idx|gsi-name"`
}

record := ExampleRecord{}
recordId := 123
tableName := "myDynamoTable"

err := client.Retrieve(tableName, recordId, &record)

if err != nil {
   // Do error handling
}
```

### Query for records
```
builder, err := client.Query(tableName, ExampleRecord{})

if err != nil {
   // Do error handling
}

result, err := builder.Equals("external_id", "12345").Execute()

if err != nil {
   // Do error handling
}

item := result.First(&record)
```

### Save an Item
```
item := ExampleRecord{
  Id:  "1234",
  ExternalId: "45678",
  GroupKey: "group123",
}

err := client.Save(tableName, item)

if err != nil {
    // Do error handling
}
```

### Delete an Item
```
err := client.Delete(tableName, "1234")

if err != nil {
    // Do error handling
}
```

### Use a transaction
```
transaction, err := client.Transaction(tableName, ExampleRecord{})
newItem := ExampleRecord{
  Id:  "1234",
  ExternalId: "new-external-id",
  GroupKey: "group123",
}

if err != nil {
    // Do error handling
}

errs := transaction.Delete("1234").Save(newItem).Commit()

if errs != nil {
    // Do error handling
}
```

### Scan items
```
items := make([]ExampleRecord, 10)

err := client.Scan(tableName, &items, 10)

if err != nil {
    // Do error handling
}
```
