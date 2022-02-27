import { Injectable } from '@angular/core';
import {setContext} from "@apollo/client/link/context";
import {ApolloLink, InMemoryCache, split} from "@apollo/client/core";
import {Apollo, ApolloBase, gql, QueryRef} from "apollo-angular";
import {AuthLibService} from "auth-lib";
import {HttpLink} from "apollo-angular/http";
import {SubscriptionClient} from "subscriptions-transport-ws";
import {WebSocketLink} from "@apollo/client/link/ws";
import {getMainDefinition} from "@apollo/client/utilities";
import {NotificationSubscriptionGql} from "./notification-subscription-gql";
import {Observable} from "rxjs";

const base = 'localhost:8082/query';

@Injectable({
  providedIn: 'root'
})
export class ApiService {

  private webSocketClient!: SubscriptionClient;
  private apollo!: ApolloBase;
  private query: QueryRef<any>;

  constructor(private apolloProvider: Apollo,
              private auth: AuthLibService,
              private httpLink: HttpLink,
              private notifSubGql: NotificationSubscriptionGql) {
    this.start();
  }

  private start() {
    const basic = setContext((operation) => ({
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

    const http = ApolloLink.from([basic, auth, this.httpLink.create({ uri: `http://${base}`})]);
    const cache = new InMemoryCache();

    this.webSocketClient = new SubscriptionClient(`ws://${base}`, {
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

    this.apolloProvider.createNamed('notification', {
      cache: cache,
      link: link
    });

    this.apollo = this.apolloProvider.use('notification');
  }

  subscribeToNotifications(userName: string): Observable<any> {
    return this.apollo.subscribe({
      query: this.notifSubGql.document,
      variables: {
        userName: userName,
      },
    });
  }

  getNotifs(userName: string): Observable<any> {
    const query = gql`
    query getNotifsByReceiver($receiver: String!) {
      getNotifsByReceiver(receiver: $receiver) {
        id,
        event,
        text,
        receiver
      }
    }
    `;

    this.query = this.apollo.watchQuery<any>({
      fetchPolicy: 'network-only',
      query: query,
      variables: {
        receiver: userName,
      },
    });

    return this.query.valueChanges;
  }
}
