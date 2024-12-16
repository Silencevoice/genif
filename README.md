# Store Generic Interface
This is mainly a coding exercise, so do not expect it to be production ready.

## The story so far
I recently found how to define an interface with a generic parameter and then use a struct with the same parameter than implements the interface generically (which I though was not possible).

## The idea
I wanted to define a generic `Store` interface that would define the behaviour of a repository (maybe `Storer` would be a better name, or simple `Repository`).

```go
type Store[T any] interface {
	GetByID(ctx context.Context, id string) (*T, error)
	GetMultipleByID(ctx context.Context, ids []string) ([]*T, error)
	GetAll(ctx context.Context) ([]*T, error)
	Insert(ctx context.Context, id string, entity *T) (*T, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, id string, entity *T) error
}
```
Where `T` would be entity to store or retrieve. I left out the specific `FindBy...` methods because those are meant to be implemented by the structures.  Also, I am asuming that the key is a string always. Will probably fix later.

### Memory implementation

Then, I implemented an *in-memory* generic implementation that would use a data map and a mutex to grant access:
```go
type MemStore[T any] struct {
	sync.RWMutex
	data map[string]T
}

func NewMemStore[T any]() *MemStore[T] {
	return &MemStore[T]{
		data: make(map[string]T),
	}
}
```

Then, that generic `MemStore[T]` implements all the methods of the `Store[T]` interface above so it could be used with any entity. 
```go
func TestInsert(t *testing.T) {
	ctx := context.Background()
	store := NewMemStore[TestEntity]()

	t.Run("Insert entity", func(t *testing.T) {
		entity := TestEntity{ID: "1", Value: "test-value"}
		inserted, err := store.Insert(ctx, entity.ID, &entity)
		require.NoError(t, err)
		assert.Equal(t, entity, *inserted)
	})

	t.Run("Insert duplicate entity", func(t *testing.T) {
		entity := TestEntity{ID: "1", Value: "test-value-duplicate"}
		_, err := store.Insert(ctx, entity.ID, &entity)
		assert.Error(t, err)
		assert.Equal(t, "already existing key", err.Error())
	})
}
```

### Mongo implementation
Then, I decided to do the same but the persistance would be a MongoDb database:
```go
type MongoStore[T any] struct {
	collection *mongo.Collection
}

func NewMongoStore[T any](db *mongo.Database, collectionName string) *MongoStore[T] {
	return &MongoStore[T]{
		collection: db.Collection(collectionName),
	}
}
```

## What about specific queries?
We will need a way to make special queries to find by other fields or to update the database