import { Injectable } from '@angular/core';
import {Apollo, ApolloBase, gql, QueryRef} from "apollo-angular";
import {AuthLibService} from "auth-lib";
import {Observable} from "rxjs";

export interface User {
  firstName: string,
  follows: string[],
  followers: string[],
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
              private auth: AuthLibService) {
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
        follows,
        followers
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
        follows,
        followers
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

  addFollow(id: string, username: string, toAdd: string): Observable<any> {
    const mutation = gql`
    mutation addFollow($id: ID!, $username: String!, $toAdd: String!){
      addFollow(id: $id, username: $username, toAdd: $toAdd) {
        id,
        username,
        firstName,
        lastName,
        follows,
        followers
      }
    }
    `;

    return this.apollo.mutate<any>({
      mutation: mutation,
      variables: {
        id: id,
        username: username,
        toAdd: toAdd,
      },
    });
  }

  removeFollow(id: string, username: string, toRemove: string): Observable<any> {
    const mutation = gql`
    mutation removeFollow($id: ID!, $username: String!, $toRemove: String!){
      removeFollow(id: $id, username: $username, toRemove: $toRemove) {
        id,
        username,
        firstName,
        lastName,
        follows,
        followers
      }
    }
    `;

    return this.apollo.mutate<any>({
      mutation: mutation,
      variables: {
        id: id,
        username: username,
        toRemove: toRemove,
      },
    });
  }
}
