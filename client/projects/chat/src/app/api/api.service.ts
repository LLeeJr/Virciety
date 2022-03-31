import { Injectable } from '@angular/core';
import {Apollo, ApolloBase, gql, QueryRef} from "apollo-angular";
import {Observable, Subject} from "rxjs";
import {ChatSubscriptionGqlService} from "./chat-subscription-gql";
import {SubscriptionClient} from "subscriptions-transport-ws";
import {HttpLink} from "apollo-angular/http";
import {WebSocketLink} from "@apollo/client/link/ws";
import {ApolloLink, from, InMemoryCache, split} from "@apollo/client/core";
import {getMainDefinition} from "@apollo/client/utilities";
import {setContext} from "@apollo/client/link/context";
import {environment} from "../../environments/environment";
import {onError} from "@apollo/client/link/error";

@Injectable({
  providedIn: 'root',
})
export class ApiService {

  messages: any[] = [];
  chatMembers: string[] = [];
  selectedRoom: any;
  errorState: Subject<string> = new Subject<string>();

  private query: QueryRef<any> | undefined;
  private apollo!: ApolloBase;
  private webSocketClient!: SubscriptionClient;

  constructor(private apolloProvider: Apollo,
              private chatSubGql: ChatSubscriptionGqlService,
              private httpLink: HttpLink) {
    this.start();
  }

  private start() {
    let errorLink = onError(({graphQLErrors, networkError }) => {
      if (networkError) {
        let msg = `Chat backend is currently offline, try again later!`;
        this.errorState.next(msg);
      }
    });

    const basic = setContext((operation, context) => ({
      headers: {
        Accept: 'charset=utf-8'
      }
    }));

    const auth = setContext((operation, context) => {
      const token = localStorage.getItem('token');

      if (token === null) {
        return {};
      } else {
        return {
          headers: {
            Authorization: `JWT ${token}`
          }
        };
      }
    });

    const http = ApolloLink.from([basic, auth, errorLink, this.httpLink.create({ uri: environment.chatAPI})]);
    const cache = new InMemoryCache({
      typePolicies: {
        Query: {
          fields: {
            getRoomsByUser: {
              merge: false,
            },
          },
        },
      },
    });

    this.webSocketClient = new SubscriptionClient(environment.chatWS, {
      reconnect: true,
      reconnectionAttempts: 3,
    });
    const ws = new WebSocketLink(this.webSocketClient);

    const link = split(
      // split based on operation type
      ({query}) => {
        let definition = getMainDefinition(query);
        return (
          definition.kind === 'OperationDefinition' && definition.operation === 'subscription'
        );
      },
      ws,
      http,
    );

    this.apolloProvider.createNamed('chat', {
      cache: cache,
      link: from([errorLink, link]),
    });

    this.apollo = this.apolloProvider.use('chat');
  }

  writeDm(msg: string, roomName: string, roomId: string, username: string): Observable<any> {

    const mutation = gql`
    mutation createDm($msg: String!, $userName: String!, $roomName: String!, $roomID: String!){
      createDm(msg: $msg, userName: $userName, roomName: $roomName, roomID: $roomID)
      {
        chatroomId,
        createdAt,
        createdBy,
        msg,
      }
    }
    `;

    return this.apollo.mutate<any>({
      mutation: mutation,
      variables: {
        msg: msg,
        userName: username,
        roomName: roomName,
        roomID: roomId,
      },
    });
  }

  subscribeToChat(roomName: string): Observable<any> {
    return this.apollo.subscribe({
      query: this.chatSubGql.document,
      variables: {
        roomName: roomName
      }
    })
  }

  unsubscribeToChat() {
    this.webSocketClient.close(true);
  }

  getMessagesFromRoom(roomId: string): Observable<any> {
    const query = gql`
    query getMessagesFromRoom($roomId: String!){
      getMessagesFromRoom(roomId: $roomId)
      {
        chatroomId,
        createdAt,
        createdBy,
        msg,
      }
    }
    `;

    this.query = this.apollo.watchQuery<any>({
      fetchPolicy: 'network-only',
      query: query,
      variables: {
        roomId: roomId,
      },
    });

    return this.query.valueChanges;
  }

  getRoomsByUser(user: string): Observable<any> {
    const query = gql`
    query getRoomsByUser($userName: String!) {
      getRoomsByUser(userName: $userName)
      {
        id,
        member,
        name,
        owner,
        isDirect
      }
    }`;

    this.query = this.apollo.watchQuery<any>({
      query: query,
      fetchPolicy: "network-only",
      variables: {
        userName: user
      }
    });

    return this.query.valueChanges;
  }

  getRoom(roomName: string, roomId: string): Observable<any> {
    const query = gql`
    query getRoom($roomName: String!, $roomID: String!) {
      getRoom(roomName: $roomName, roomID: $roomID)
      {
        id,
        member,
        name,
        owner
      }
    }`;

    this.query = this.apollo.watchQuery<any>({
      query: query,
      fetchPolicy: "network-only",
      variables: {
        roomName: roomName,
        roomID: roomId,
      }
    });

    return this.query.valueChanges;
  }

  getDirectRoom(user1: string, user2: string): Observable<any> {
    const query = gql`
    query getDirectRoom($user1: String!, $user2: String!) {
      getDirectRoom(user1: $user1, user2: $user2)
      {
        id,
        member,
        name,
        owner,
        isDirect
      }
    }`;

    this.query = this.apollo.watchQuery<any>({
      query: query,
      variables: {
        user1: user1,
        user2: user2,
      }
    });

    return this.query.valueChanges;
  }

  createRoom(member: string[], name: string, owner: string, isDirect: boolean) {
    const mutation = gql`
    mutation createRoom($input: CreateRoom!){
      createRoom(input: $input)
      {
        id,
        member,
        name,
        owner,
        isDirect
      }
    }
    `;

    const input = {
      member: member,
      name: name,
      owner: owner,
      isDirect: isDirect,
    }

    return this.apollo.mutate<any>({
      mutation: mutation,
      variables: {
        input: input,
      },
    });
  }

  deleteRoom(roomName: string, roomId: string, owner: string) {
    const mutation = gql`
    mutation deleteRoom($remove: RemoveRoom!){
      deleteRoom(remove: $remove)
    }
    `;

    const remove = {
      id: roomId,
      roomName: roomName,
      userName: owner
    }

    return this.apollo.mutate<any>({
      mutation: mutation,
      variables: {
        remove: remove,
      }
    })
  }

  leaveChat(roomId: string, username: string, owner: string | undefined): Observable<any> {
    const mutation = gql`
    mutation leaveChat($roomId: String!, $userName: String!, $owner: String) {
      leaveChat(roomId: $roomId, userName: $userName, owner: $owner) {
        id,
      }
    }
    `;

    return this.apollo.mutate<any>({
      mutation: mutation,
      variables: {
        roomId: roomId,
        userName: username,
        owner: owner,
      },
    });
  }
}
