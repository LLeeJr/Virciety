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

@Injectable({
  providedIn: 'root',
})
export class ApiService {

  messages: any[] = [];
  chatPartner = '';

  private query: QueryRef<any> | undefined;
  private apollo!: ApolloBase;
  private webSocketClient!: SubscriptionClient;

  constructor(private apolloProvider: Apollo,
              private datePipe: DatePipe,
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

  writeDm(msg: string, from: string, to: string): Observable<any> {
    const date = new Date();
    const formDate =this.datePipe.transform(date, 'yyyy-MM-dd HH:mm:ss');
    const id = `${from}__${formDate}__${to}`;

    const createDmRequest = {
      id: id,
      msg: msg,
    };

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
        input: createDmRequest
      }
    });
  }

  subscribeToChat(): Observable<any> {
    return this.apollo.subscribe({query: this.chatSubGql.document})
  }

  unsubscribeToChat() {
    this.webSocketClient.unsubscribeAll();
    this.webSocketClient.close(true);
  }

  getChat(user1: string, user2: string): Observable<any> {
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

  getDms(): Observable<any> {
    const DMS_QUERY = gql`
    query {
      getDms {
        id,
        msg
      }
    }
  `;

    this.query = this.apollo.watchQuery({
      query: DMS_QUERY,
      variables: {}
    });

    return this.query.valueChanges;
  }
}
