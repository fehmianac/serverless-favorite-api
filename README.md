# Serverless Favorite API

The Serverless Favorite API is a RESTful API designed to manage favorite items in a microservices architecture. It is built specifically for use on AWS Services, offering easy integration with your existing AWS infrastructure.

## Features

The API provides the following main functionalities:

1. **Add Favorite**: Add a favorite item for a specific user.
   - URL: `/user/{userId}/favorite`
   - Method: POST
   - Request Body: JSON object containing the `itemId` parameter

2. **Delete Favorite**: Remove a favorite item for a specific user.
   - URL: `/user/{userId}/favorite/{itemId}`
   - Method: DELETE

3. **Check Favorites**: Check if a list of items are marked as favorites for a specific user.
   - URL: `/user/{userId}/favorite?itemIds=123,1234`
   - Method: GET
   - Query Parameters: `itemIds` (comma-separated list of item IDs)

4. **Get Favorites**: Retrieve a paginated list of favorite items for a specific user.
   - URL: `/user/{userId}/favorite`
   - Method: GET
   - Query Parameters: `nextToken` (pagination token), `limit` (number of items per page)

## Getting Started

To use the Serverless Favorite API, follow these steps:

## Examples

Here are some examples demonstrating how to use the Serverless Favorite API:

- Adding a favorite item:
  ```
  POST /user/{userId}/favorite
  Request Body: {"itemId": "your-item-id"}
  ```

- Deleting a favorite item:
  ```
  DELETE /user/{userId}/favorite/{itemId}
  ```

- Checking favorites:
  ```
  GET /user/{userId}/favorite?itemIds=123,1234
  ```

- Getting paginated favorites:
  ```
  GET /user/{userId}/favorite?nextToken=your-pagination-token&limit=10
  ```

## Contributing

Contributions to the Serverless Favorite API are welcome! If you find any issues or have suggestions for improvements, please feel free to submit a pull request or open an issue in the project repository.

## License

The Serverless Favorite API is released under the [MIT License](LICENSE).
