import { Injectable } from '@angular/core';
import {Apollo, ApolloBase, gql, QueryRef} from "apollo-angular";
import {AuthLibService} from "auth-lib";
import {Observable} from "rxjs";

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
}
