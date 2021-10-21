import { Injectable } from '@angular/core';
import {Apollo, ApolloBase, gql, QueryRef} from "apollo-angular";
import {Observable, Subscription} from "rxjs";
import {DatePipe} from "@angular/common";
import {ChatSubscriptionGqlService} from "./chat-subscription-gql";
import {SubscriptionClient} from "subscriptions-transport-ws";
import {HttpLink} from "apollo-angular/http";
import {WebSocketLink} from "@apollo/client/link/ws";
import {InMemoryCache, split} from "@apollo/client/core";
import {getMainDefinition} from "@apollo/client/utilities";
import {AuthLibService} from "auth-lib";
import {Room} from "../data/room";

@Injectable({
  providedIn: 'root',
})
export class ApiService {

  messages: any[] = [];
  chatMembers: string[] = [];
  selectedRoom: Room;

  private query: QueryRef<any> | undefined;
  private apollo!: ApolloBase;
  private webSocketClient!: SubscriptionClient;

  constructor(private apolloProvider: Apollo,
              private auth: AuthLibService,
              private chatSubGql: ChatSubscriptionGqlService,
              private datePipe: DatePipe,
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
    const roomName = this.selectedRoom;

    const mutation = gql`
    mutation createDm($input: CreateDmRequest!){
      createDm(input: $input)
      {
        id,
        msg
      }
    }
    `;

    return this.apollo.mutate<any>({
      mutation: mutation,
      variables: {
        msg: msg,
        username: userName,
        roomName: roomName,
      },
    });
  }

  subscribeToChat(): Observable<any> {
    return this.apollo.subscribe({
      query: this.chatSubGql.document,
      variables: {
        roomName: this.selectedRoom.name
      }
    })
  }

  unsubscribeToChat() {
    this.webSocketClient.unsubscribeAll();
    this.webSocketClient.close(true);
  }

  getMessagesFromRoom(): Observable<any> {
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
      query: query,
      variables: {
        roomId: this.selectedRoom._id,
      },
    });

    return this.query.valueChanges;
  }

  getRoomsByUser(user: string): Observable<any> {
    const query = gql`
    query getRoomsByUser($userName: String!) {
      getRoomsByUser(userName: $userName)
      {
        name,
        member,
        _id
      }
    }`;

    this.query = this.apollo.watchQuery<any>({
      query: query,
      variables: {
        userName: user
      }
    });

    return this.query.valueChanges;
  }

  getRoom(roomName: string): Observable<any> {
    const query = gql`
    query getRoom($name: String!) {
      getRoom(name: $name)
      {
        name,
        member,
        messages{
          msg
        },
      }
    }`;

    this.query = this.apollo.watchQuery<any>({
      query: query,
      variables: {
        name: roomName
      }
    });

    return this.query.valueChanges;
  }

}
