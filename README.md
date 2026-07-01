# blogs

GraphQL API built with [gqlgen](https://gqlgen.com) — a code-first Go library for GraphQL servers.

## Schema

```graphql
type Post {
  id: ID!
  title: String!
  body: String!
  author: User!
  comments: [Comment!]!
}

type Comment {
  id: ID!
  body: String!
  author: User!
  post: Post!
}

type User {
  id: ID!
  name: String!
}

type Query {
  posts: [Post!]!
  post(id: ID!): Post
}

type Mutation {
  createPost(input: NewPost!): Post!
  createComment(input: NewComment!): Comment!
}
```

## Quickstart

```bash
go run ./...
```

Open http://localhost:8080 in your browser.

## Example queries

**Create a post:**

```graphql
mutation {
  createPost(input: { title: "Hello World", body: "First post!", authorId: "1" }) {
    id
    title
    author { name }
  }
}
```

**Add a comment:**

```graphql
mutation {
  createComment(input: { body: "Great post!", authorId: "2", postId: "1" }) {
    id
    body
    author { name }
    post { title }
  }
}
```

**List posts with comments:**

```graphql
query {
  posts {
    id
    title
    body
    author { name }
    comments {
      body
      author { name }
    }
  }
}
```

**Get a single post:**

```graphql
query {
  post(id: "1") {
    title
    author { name }
  }
}
```

## Development

- `make generate` — regenerate code from schema
- `make build` — compile binary
- `make run` — start server
