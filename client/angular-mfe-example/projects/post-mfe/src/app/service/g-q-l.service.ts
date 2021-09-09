import { Injectable } from '@angular/core';
import {Apollo, ApolloBase, gql} from "apollo-angular";
import {Observable} from "rxjs";
import {WebSocketLink} from "@apollo/client/link/ws";
import {InMemoryCache, split} from "@apollo/client/core";
import {getMainDefinition} from "@apollo/client/utilities";
import {HttpLink} from "apollo-angular/http";
import {SubscriptionClient} from "subscriptions-transport-ws";

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

  private UPLOAD_FILE = gql`
    mutation upload($file: String!) {
        upload(file: $file) {
          id
          content
          contentType
        }
      }
  `;

  private CREATE_POST = gql`
      mutation CreatePost($newPost: CreatePostRequest!) {
        createPost(newPost: $newPost) {
          id
          description
          data
          likedBy
          comments
        }
      }
  `

  private POST_SUBSCRIPTION = gql`
    subscription postCreated {
      postCreated
    }
  `;

  // ------------------------- Queries, Mutations and Subscriptions end

  private apollo: ApolloBase;
  private readonly webSocketClient: SubscriptionClient;

  constructor(private apolloProvider: Apollo,
              private httpLink: HttpLink) {
    const http = httpLink.create({
      uri: 'http://localhost:8083/query',
    });

    this.webSocketClient = new SubscriptionClient(`ws://localhost:8083/query`, {
      reconnect: true
    });

    const ws = new WebSocketLink(this.webSocketClient);

    const link = split(
      // split based on operation type
      ({query}) => {
        const data = getMainDefinition(query);
        return (
          data.kind === 'OperationDefinition' && data.operation === 'subscription'
        );
      },
      ws,
      http,
    );

    this.apolloProvider.createNamed('post', {
      cache: new InMemoryCache(),
      link: link,
    });

    this.apollo = this.apolloProvider.use('post');
  }

  public closeWebSocket(): void {
    this.webSocketClient.unsubscribeAll();
    this.webSocketClient.close(true);
    console.log(`Closed Websocket? ${this.webSocketClient.status}`);
  }

  getPosts(): Observable<any> {
    return this.apollo
      .watchQuery({
        query: this.GET_POSTS,
      })
      .valueChanges;
  }

  postCreated(): Observable<any> {
    return this.apollo.subscribe({query: this.POST_SUBSCRIPTION});
  }

  upload(fileBase64: string): Observable<any> {
    return this.apollo.mutate({
      mutation: this.UPLOAD_FILE,
      variables: {
        file: fileBase64
      }
    });
  }
}
