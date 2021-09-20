import {gql} from "apollo-angular";

// ------------------------- Queries, Mutations and Subscriptions

export const GET_POSTS = gql`
    query getPosts($id: String!, $fetchLimit: Int!) {
      getPosts(id: $id, fetchLimit: $fetchLimit) {
        id
        data {
          name
          contentType
        }
        description
        comments
        likedBy
      }
    }
  `;

export const GET_DATA = gql`
    query getData($id: String!) {
      getData(id: $id)
    }
  `;

export const CREATE_POST = gql`
    mutation createPost($username: String!, $description: String!, $data: String!) {
      createPost(newPost: {username: $username, description: $description, data: $data}) {
        id
        description
        data {
          name
          contentType
        }
        likedBy
        comments
      }
    }
  `;

export const NEW_POST_CREATED = gql`
    subscription newPostCreated {
      newPostCreated {
        id
        description
        data {
          name
          contentType
        }
        likedBy
        comments
      }
    }
  `;

// ------------------------- Queries, Mutations and Subscriptions end
