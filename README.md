# RestJSON

RestJSON is an open-source tool that combines the power of a JSON editor with automatic REST API generation. Create and edit your JSON data through a convenient online interface, and instantly get a fully functional REST API without any additional configuration.

## Motivation
- Quickly create a mock API for testing and development
- Prototype and experiment with API endpoints

## Features

- **Online JSON Editor**: Create and modify your JSON data structure with real-time validation
- **Instant API Generation**: Your JSON automatically becomes a RESTful API endpoint
- **CRUD Operations**: Full support for GET, POST, PUT, PATCH, DELETE operations
- (TODO) **Filtering & Sorting**: Query your data with powerful filtering options
- **No Backend Required**: Everything works out of the box


### Accessing your API

Once your JSON file is created, it's automatically available as a REST API. Access your data using these endpoints:

```
# Access entire JSON file
GET    /public/{fileId}

# For resources (objects)
GET    /public/{fileId}/{resource}         # Get all resource data
PUT    /public/{fileId}/{resource}         # Replace entire resource
PATCH  /public/{fileId}/{resource}         # Partially update resource

# For array items within resources
GET    /public/{fileId}/{resource}/{id}    # Get a specific item by ID
POST   /public/{fileId}/{resource}         # Create a new item in the array
PUT    /public/{fileId}/{resource}/{id}    # Replace an item completely
PATCH  /public/{fileId}/{resource}/{id}    # Partially update an item
DELETE /public/{fileId}/{resource}/{id}    # Delete an item
```

Note: API access requires authentication with an API key, which can be obtained in the Account page.

## Examples

### Sample JSON

```json
{
  "posts": [
    { "id": 1, "title": "First post", "author": "John" },
    { "id": 2, "title": "Second post", "author": "Jane" }
  ],
  "comments": [
    { "id": 1, "body": "Great post!", "postId": 1 },
    { "id": 2, "body": "I agree!", "postId": 1 }
  ]
}
```

### API Endpoints Created

- `GET /public/{fileId}/posts` - Get all posts
- `GET /public/{fileId}/posts/1` - Get post with id 1
- `POST /public/{fileId}/posts` - Create a new post
- `PUT /public/{fileId}/posts/1` - Replace post with id 1
- `PATCH /public/{fileId}/posts/1` - Update specific fields of post with id 1
- `DELETE /public/{fileId}/posts/1` - Delete post with id 1

## Built With
- Frontend: React with Vite
- Backend: Go
- Database: PostgreSQL
- File Storage: AWS S3
- Cache: Redis

## Credits

RestJSON is inspired by:

- [JSON Server](https://github.com/typicode/json-server) 
- [json-bucket](https://github.com/Nico-Mayer/json-bucket)

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

