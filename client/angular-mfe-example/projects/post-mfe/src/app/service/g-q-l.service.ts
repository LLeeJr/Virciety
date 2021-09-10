import { Injectable } from '@angular/core';
import {Apollo, ApolloBase, gql} from "apollo-angular";
import {Observable} from "rxjs";
import {InMemoryCache} from "@apollo/client/core";
import {HttpLink} from "apollo-angular/http";

@Injectable({
  providedIn: 'root'
})
export class GQLService {
  // ------------------------- Queries, Mutations and Subscriptions

  private GET_POSTS = gql`
          {
            getPosts {
              id
              data {
                id
                content
                contentType
              }
              description
              comments
              likedBy
            }
          }`;

  private CREATE_POST = gql`
    mutation createPost($username: String!, $description: String!, $data: String!) {
      createPost(newPost: {username: $username, description: $description, data: $data}) {
        id
        description
        data {
          id
          content
          contentType
        }
        likedBy
        comments
      }
    }
  `;

  // ------------------------- Queries, Mutations and Subscriptions end

  private apollo: ApolloBase;

  constructor(private apolloProvider: Apollo,
              private httpLink: HttpLink) {
    const http = httpLink.create({
      uri: 'http://localhost:8083/query',
    });

    this.apolloProvider.createNamed('post', {
      cache: new InMemoryCache(),
      link: http,
    });

    this.apollo = this.apolloProvider.use('post');
  }

  getPosts(): Observable<any> {
    return this.apollo
      .watchQuery({
        query: this.GET_POSTS,
      })
      .valueChanges;
  }

  createPost(fileBase64: string, description: string, username: string): Observable<any> {
    return this.apollo.mutate({
      mutation: this.CREATE_POST,
      variables: {
        username: username,
        description: description,
        data: fileBase64
      }
    });
  }
}
