import { Injectable } from '@angular/core';
import {setContext} from "@apollo/client/link/context";
import {ApolloLink, from, InMemoryCache, split} from "@apollo/client/core";
import {Apollo, ApolloBase, gql, QueryRef} from "apollo-angular";
import {HttpLink} from "apollo-angular/http";
import {SubscriptionClient} from "subscriptions-transport-ws";
import {WebSocketLink} from "@apollo/client/link/ws";
import {getMainDefinition} from "@apollo/client/utilities";
import {NotificationSubscriptionGql} from "./notification-subscription-gql";
import {Observable, Subject} from "rxjs";
import {onError} from "@apollo/client/link/error";
import {environment} from "../../environments/environment";

@Injectable({
  providedIn: 'root'
})
export class ApiService {

  private webSocketClient!: SubscriptionClient;
  private apollo!: ApolloBase;
  private query: QueryRef<any>;
  errorState: Subject<string> = new Subject<string>();

  constructor(private apolloProvider: Apollo,
              private httpLink: HttpLink,
              private notifSubGql: NotificationSubscriptionGql) {
    this.start();
  }

  private start() {
    let errorLink = onError(({graphQLErrors, networkError }) => {
      if (networkError) {
        let msg = `Notifications are currently offline!`;
        this.errorState.next(msg);
      }
    });
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

    const http = ApolloLink.from([basic, auth, this.httpLink.create({ uri: environment.notifsAPI})]);
    const cache = new InMemoryCache();

    this.webSocketClient = new SubscriptionClient(environment.notifsWS, {
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

    this.apolloProvider.createNamed('notification', {
      cache: cache,
      link: from([errorLink, link]),
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

  getNotifications(userName: string): Observable<any> {
    const query = gql`
    query getNotifsByReceiver($receiver: String!) {
      getNotifsByReceiver(receiver: $receiver) {
        id,
        event,
        read,
        receiver,
        text,
        timestamp,
        params{
          key,
          value
        },
        route
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

  setReadStatus(id: string, status: boolean): Observable<any> {
    const mutation = gql`
    mutation setReadStatus($id: String!, $status: Boolean!) {
      setReadStatus(id: $id, status: $status) {
        id,
        event,
        read,
        receiver,
        text,
        timestamp,
        params{
          key,
          value
        },
        route
      }
    }
    `;

    return this.apollo.mutate<any>({
      mutation: mutation,
      variables: {
        id: id,
        status: status,
      },
    });
  }
}
