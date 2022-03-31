import {Inject, Injectable} from '@angular/core';
import {Observable, Subject} from "rxjs";
import {Apollo, ApolloBase, gql, QueryRef} from "apollo-angular";
import {HttpLink} from "apollo-angular/http";
import {from, InMemoryCache} from "@apollo/client/core";
import {onError} from "@apollo/client/link/error";

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
  private _userName = new Subject<string | undefined>();
  _activeId = new Subject<string>();
  errorState: Subject<string> = new Subject<string>();
  error: any = null;

  userName: string = '';
  activeId: string = '';
  private environment: any;


  constructor(private apolloProvider: Apollo,
              private httpLink: HttpLink,
              @Inject('environment') environment: any) {
    this.environment = environment;
    this.start();
  }

  private start() {
    let errorLink = onError(({graphQLErrors, networkError }) => {
      if (networkError) {
        let msg = `User backend is currently offline!`;
        this.errorState.next(msg);
        this.error = msg;
      }
    });
    this.apollo = this.apolloProvider.use('user');
    if (this.apollo === undefined) {
      const http = this.httpLink.create({
        uri: this.environment.userAPI
      })
      this.apolloProvider.createNamed('user', {
        cache: new InMemoryCache(),
        link: from([errorLink, http]),
      });
      this.apollo = this.apolloProvider.use('user');
    }
  }

  getUserName(): Subject<string | undefined> {
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

  addFollow(id: string, username: string, toAdd: string): Observable<any> {
    const mutation = gql`
    mutation addFollow($id: ID!, $username: String!, $toAdd: String!){
      addFollow(id: $id, username: $username, toAdd: $toAdd) {
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
        followers,
        profilePictureId
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
