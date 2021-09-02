import { Injectable } from '@angular/core';
import {Apollo, gql, QueryRef} from "apollo-angular";
import {Observable} from "rxjs";
import {DatePipe} from "@angular/common";

@Injectable({
  providedIn: 'root',
})
export class ApiService {

  messages : any[] = [];

  private query: QueryRef<any> | undefined;

  constructor(private apollo: Apollo,
              private datePipe: DatePipe) {
  }

  writeDm(msg: string, from: string, to: string): Observable<any> {
    const date = new Date();
    const formDate =this.datePipe.transform(date, 'yyyy-MM-dd HH:mm:ss');
    const id = `${from}__${formDate}__${to}`;

    const createDmRequest = {
      id: id,
      msg: msg,
    };

    const mutation = gql`
    mutation createDm($input: CreateDmRequest!){
      createDm(input: $input)
      {
        id,
        msg
      }
    }
    `;

    return this.apollo.mutate<any>({
      mutation: mutation,
      variables: {
        input: createDmRequest
      }
    });
  }

  getChat(user1: string, user2: string): Observable<any> {
    const query = gql`
    query getChat($input: GetChatRequest!){
      getChat(input: $input)
      {
        id,
        msg
      }
    }
    `;

    const getChatRequest = {
      user1: user1,
      user2: user2,
    }

    this.query = this.apollo.watchQuery<any>({
      query: query,
      variables: {
        input: getChatRequest
      },
    });

    return this.query.valueChanges;
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
