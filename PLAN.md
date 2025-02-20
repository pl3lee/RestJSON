# Plan
## What this project is about
Something like json-server, but a web hosted version. Users can just create a json file on the app, then it automatically creates an API out of the JSON, which the user can make HTTP requests to from anywhere.

## Main Technologies
- Frontend: Vite React, allows for live editing and auto saving of JSON file
- Backend: Go, handles user authentication, endpoints for creating, editing JSON. Also has an API for users to access their JSON file
- Database: PostgreSQL, stores user information, and stores file metadata
- Object storage: AWS S3, stores the JSON files themselves
- Cache: Redis, to minimize S3 reads.

There will be 2 Go services:
1. frontend-api: specifically used to serve the frontend webapp.
2. public-api: used to serve the general public, where users can make HTTP requests from anywhere and treat the JSON as database.

## Flow
### User logs in through web frontend
1. User logins to Google/Github using OAuth
2. Frontend receives authorization code
3. Frontend makes HTTP request to frontend-api with authorization code
4. frontend-api verifies authorization code by sending request to Google/Github
5. Google/Github returns user information
6. frontend-api creates access token in the form of JWT, and refresh token (stored in database), and sends tokens back using cookies
7. Frontend sends both tokens automatically for every request

### User writing/updating JSON file on web frontend
1. User creates/updates JSON file
2. Frontend sends the updated or new JSON payload to frontend-api
3. If creating, then frontend-api generates a unique file ID and corresponding S3 key, stores metadata in database, and uploads to S3
4. If updating, then frontend-api uses metadata from database to upload to S3
5. Invalidate Redis cache entry for updated file

### User accessing JSON file through public-api (GET)
1. public-api receives user request to view or retrieve JSON file
2. In the request, a file ID is included, and an API key is included in Authorization header
3. public-api checks database to verify that user associated with API key has access to the file
4. public-api checks Redis for the JSON file
5. If in Redis, then return the cached JSON data
6. If not, then query S3, store in Redis, and return JSON data

### User updating JSON file through public-api (POST, PUT)
1. public-api receives user request to update JSON file
2. public-api checks API key to ensure the user has permission
3. public-api fetches JSON file from Redis or S3, and modifies it based on the user request
4. Uploads the updated JSON file back to S3, invalidates Redis cache entry
5. Sends response back to user of the updated JSON

