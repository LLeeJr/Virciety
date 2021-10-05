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

@Injectable({
  providedIn: 'root',
})
export class ApiService {

  messages: any[] = [];
  chatPartner = '';
  selectedRoom = '';

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
    return this.apollo.subscribe({query: this.chatSubGql.document})
  }

  unsubscribeToChat() {
    this.webSocketClient.unsubscribeAll();
    this.webSocketClient.close(true);
  }

  getChat(): Observable<any> {
    const user1 = this.auth.userName;
    const user2 = this.chatPartner;
    const query = gql`
    query getChat($input: GetChatRequest!){
      getChat(input: $input)
      {
        id,
        msg
      }
    }
    `;

    const getChatRequest = {
      user1: user1,
      user2: user2,
    }

    this.query = this.apollo.watchQuery<any>({
      query: query,
      variables: {
        input: getChatRequest
      },
    });

    return this.query.valueChanges;
  }

  getOpenChats(userName: string): Observable<any> {
    const query = gql`
    query getOpenChats($userName: String!){
      getOpenChats(userName: $userName)
      {
        withUser,
        preview
      }
    }
    `;

    this.query = this.apollo.watchQuery<any>({
        query: query,
        variables: {
          userName: userName,
        },
      });

    return this.query.valueChanges;
  }

  getRoomsByUser(): Observable<any> {
    const userName = this.auth.userName;
    const query = gql`
    query getRoom($userName: String!) {
      getRoom(userName: $userName)
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
        userName: userName
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
