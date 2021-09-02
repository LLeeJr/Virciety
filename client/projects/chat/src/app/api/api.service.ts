import { Injectable } from '@angular/core';
import {Apollo, gql, QueryRef} from "apollo-angular";
import {Observable} from "rxjs";

@Injectable({
  providedIn: 'root'
})
export class ApiService {

  messages : any[] = [];

  private query: QueryRef<any> | undefined;

  constructor(private apollo: Apollo) {
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
