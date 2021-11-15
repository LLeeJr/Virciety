import { Injectable } from '@angular/core';
import {Apollo, ApolloBase, gql, QueryRef} from "apollo-angular";
import {AuthLibService} from "auth-lib";
import {HttpLink} from "apollo-angular/http";
import {Observable} from "rxjs";

export interface User {
  firstName: string,
  follows: string[],
  id: string,
  lastName: string,
  username: string,
  __typename: string,
}

@Injectable({
  providedIn: 'root'
})
export class ApiService {

  private query: QueryRef<any>;
  private apollo: ApolloBase;

  activeId: string;

  constructor(private apolloProvider: Apollo,
              private auth: AuthLibService,
              private httpLink: HttpLink) {
    this.start();
    this.auth._activeId.subscribe((id: string) => {
      if (id) {
        this.activeId = id;
      }
    });
  }

  private start() {
    this.apollo = this.apolloProvider.use('user');
  }

  findUsersWithName(name: string): Observable<any> {
    const query = gql`
    query findUsersWithName($name: String!) {
      findUsersWithName(name: $name)
      {
        username,
        firstName,
        lastName,
      }
    }
    `;

    this.query = this.apollo.watchQuery<any>({
      query: query,
      variables: {
        name: name,
      },
    });

    return this.query.valueChanges;
  }

  getUserByName(name: string): Observable<any> {
    const query = gql`
    query getUserByName($name: String!) {
      getUserByName(name: $name)
      {
        id,
        username,
        firstName,
        lastName,
        follows
      }
    }
    `;

    this.query = this.apollo.watchQuery<any>({
      query: query,
      variables: {
        name: name,
      },
    });

    return this.query.valueChanges;
  }

  getUserByID(id: string): Observable<any> {
    const query = gql`
    query getUserByID($id: ID!) {
      getUserByID(id: $id)
      {
        id,
        username,
        firstName,
        lastName,
        follows
      }
    }
    `;

    this.query = this.apollo.watchQuery<any>({
      query: query,
      variables: {
        id: id,
      },
    });

    return this.query.valueChanges;
  }

  addFollow(id: string, toAdd: string): Observable<any> {
    const mutation = gql`
    mutation addFollow($id: ID!, $toAdd: String!){
      addFollow(id: $id, toAdd: $toAdd)
    }
    `;

    return this.apollo.mutate<any>({
      mutation: mutation,
      variables: {
        id: id,
        toAdd: toAdd,
      },
    });
  }

  removeFollow(id: string, toRemove: string): Observable<any> {
    const mutation = gql`
    mutation removeFollow($id: ID!, $toRemove: String!){
      removeFollow(id: $id, toRemove: $toRemove)
    }
    `;

    return this.apollo.mutate<any>({
      mutation: mutation,
      variables: {
        id: id,
        toRemove: toRemove,
      },
    });
  }
}
