import { Injectable } from '@angular/core';
import {Apollo, ApolloBase, gql, QueryRef} from "apollo-angular";
import {Observable} from "rxjs";
import {ChatSubscriptionGqlService} from "./chat-subscription-gql";
import {SubscriptionClient} from "subscriptions-transport-ws";
import {HttpLink} from "apollo-angular/http";
import {WebSocketLink} from "@apollo/client/link/ws";
import {InMemoryCache, split} from "@apollo/client/core";
import {getMainDefinition} from "@apollo/client/utilities";
import {AuthLibService} from "auth-lib";

@Injectable({
  providedIn: 'root',
})
export class ApiService {

  messages: any[] = [];
  chatMembers: string[] = [];
  selectedRoom: any;

  private query: QueryRef<any> | undefined;
  private apollo!: ApolloBase;
  private webSocketClient!: SubscriptionClient;

  constructor(private apolloProvider: Apollo,
              private auth: AuthLibService,
              private chatSubGql: ChatSubscriptionGqlService,
              private httpLink: HttpLink) {
    this.start();
  }

  private start() {
    const http = this.httpLink.create({
      uri: 'http://localhost:8081/query'
    });

    this.webSocketClient = new SubscriptionClient('ws://localhost:8081/query', {
      reconnect: true,
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
      cache: new InMemoryCache(),
      link: link
    });

    this.apollo = this.apolloProvider.use('chat');
  }

  writeDm(msg: string): Observable<any> {
    const userName = this.auth.userName;
    const roomName = this.selectedRoom.name;
    const roomId = this.selectedRoom.id;

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
        userName: userName,
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

  getRoom(): Observable<any> {
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
      variables: {
        roomName: this.selectedRoom.name,
        roomID: this.selectedRoom.id,
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
        owner
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
        owner
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
