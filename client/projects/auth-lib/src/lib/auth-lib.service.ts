import { Injectable } from '@angular/core';
import {Observable, Subject} from "rxjs";
import {Apollo, ApolloBase, gql, QueryRef} from "apollo-angular";
import {HttpLink} from "apollo-angular/http";
import {InMemoryCache} from "@apollo/client/core";
import {KeycloakService} from "keycloak-angular";

@Injectable({
  providedIn: 'root'
})
export class AuthLibService {

  private apollo: ApolloBase;
  private query: QueryRef<any>;
  private _userName = new Subject<string>()
  _activeId = new Subject<string>()

  userName: string = '';
  activeId: string = '';


  constructor(private apolloProvider: Apollo,
              private httpLink: HttpLink,
              private keycloak: KeycloakService) {
    this.start();
    // this.keycloak.
  }

  private start() {
    this.apollo = this.apolloProvider.use('user');
    if (this.apollo === undefined) {
      const http = this.httpLink.create({
        uri: "http://localhost:8085/query"
      })
      this.apolloProvider.createNamed('user', {
        cache: new InMemoryCache(),
        link: http,
      })
      this.apollo = this.apolloProvider.use('user');
    }
  }

  getUserName(): Subject<string> {
    return this._userName;
  }

  login(firstName: string | undefined, lastName: string | undefined, username: string | undefined) {
    // console.log('Login: ', userName);
    this._userName.next(username);
    this.userName = username!;

    this.getUserByName(username!).subscribe(value => {
      if (value && value.data && value.data.getUserByName) {
        this._activeId.next(value.data.getUserByName.id);
        this.activeId = value.data.getUserByName.id;
      }
    }, () => {
      this.createUser(firstName!, lastName!, username!).subscribe(value => {
        if (value && value.data && value.data.createUser) {
          this._activeId.next(value.data.getUserByName.id);
          this.activeId = value.data.createUser.id;
        }
      });
    });
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

  createUser(firstName: string, lastName: string, username: string) {
    const mutation = gql`
    mutation createUser($input: UserData!) {
      createUser(input: $input) {
        id,
        username
      }
    }`;

    const input = {
      username: username,
      firstName: firstName,
      lastName: lastName,
    };

    return this.apollo.mutate<any>({
      mutation: mutation,
      variables: {
        input: input
      }
    });
  }
}
