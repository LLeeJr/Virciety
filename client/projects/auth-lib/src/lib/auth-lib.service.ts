import { Injectable } from '@angular/core';
import {Observable, Subject} from "rxjs";
import {Apollo, ApolloBase, gql, QueryRef} from "apollo-angular";
import {HttpLink} from "apollo-angular/http";
import {InMemoryCache} from "@apollo/client/core";
import {KeycloakService} from "keycloak-angular";

export interface User {
  firstName: string,
  follows: string[],
  followers: string[],
  id: string,
  lastName: string,
  profilePictureId: string;
  username: string,
  __typename: string,
}
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
        followers,
        profilePictureId
      }
    }
    `;

    this.query = this.apollo.watchQuery<any>({
      fetchPolicy: 'network-only',
      query: query,
      variables: {
        id: id,
      },
    });

    return this.query.valueChanges;
  }

  addProfilePicture(contentType: string, fileBase64: any, username: string): Observable<any> {
    const mutation = gql`
    mutation addProfilePicture($input: AddProfilePicture!) {
        addProfilePicture(input: $input)
      }
    `;

    const input = {
      username: username,
      data: fileBase64,
    };

    return this.apollo.mutate({
      mutation: mutation,
      variables: {
        input: input,
      },
    })
  }

  getProfilePicture(fileID: string): Observable<any> {
    const query = gql`
    query getProfilePicture($fileID: String!) {
      getProfilePicture(fileID: $fileID)
    }
    `;

    this.query = this.apollo.watchQuery<any>({
      query: query,
      variables: {
        fileID: fileID,
      },
    });

    return this.query.valueChanges;
  }

  removeProfilePicture(username: string, fileID: string): Observable<any> {
    const mutation = gql`
    mutation removeProfilePicture($remove: RemoveProfilePicture!) {
      removeProfilePicture(remove: $remove)
    }
    `;

    const remove = {
      username: username,
      fileID: fileID,
    };

    return this.apollo.mutate({
      mutation: mutation,
      variables: {
        remove: remove,
      },
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
        follows,
        followers,
        profilePictureId
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
