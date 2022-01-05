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
  profilePictureId: string;
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
  user: any;

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
}
