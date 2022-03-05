import {gql} from "apollo-angular";

// ------------------------- Queries, Mutations and Subscriptions

export const GET_POSTS = gql`
    query getPosts($id: String!, $fetchLimit: Int!, $filter: String) {
      getPosts(id: $id, fetchLimit: $fetchLimit, filter: $filter) {
        id
        data {
          name
          contentType
        }
        description
        username
        comments
        likedBy
      }
    }
  `;

export const GET_DATA = gql`
    query getData($fileID: String!) {
      getData(fileID: $fileID)
    }
  `;

export const GET_POST_COMMENTS = gql`
    query getPostComments($id: String!) {
      getPostComments(id: $id) {
        comments{
          id
          postID
          comment
          createdBy
        },
        userIdMap{
          key,
          value
        },
      }
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
        username
        likedBy
        comments
      }
    }
  `;

export const LIKE_POST = gql`
    mutation likePost($id: String!, $description: String!, $newLikedBy: [String!]!, $postOwner: String!, $comments: [String!]!, $liked: Boolean!, $likedBy: String!) {
      likePost(like: {id: $id, description: $description, newLikedBy: $newLikedBy, postOwner: $postOwner, comments: $comments, liked: $liked, likedBy: $likedBy})
    }
  `;

export const EDIT_POST = gql`
    mutation editPost($id: String!, $newDescription: String!, $likedBy: [String!]!, $comments: [String!]!) {
      editPost(edit: {id: $id, newDescription: $newDescription, likedBy: $likedBy, comments: $comments})
    }
  `;

export const REMOVE_POST = gql`
    mutation removePost($id: String!, $fileID: String!) {
      removePost(remove: {id: $id, fileID: $fileID})
    }
  `;

export const ADD_COMMENT = gql`
    mutation addComment($comment: AddCommentRequest!) {
      addComment(comment: $comment) {
        id
        postID
        comment
        createdBy
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
        username
        likedBy
        comments
      }
    }
  `;

// ------------------------- Queries, Mutations and Subscriptions end
