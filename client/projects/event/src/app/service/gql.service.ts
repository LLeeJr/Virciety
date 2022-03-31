import { Injectable } from '@angular/core';
import {Apollo, ApolloBase} from "apollo-angular";
import {HttpLink} from "apollo-angular/http";
import {DocumentNode, from, InMemoryCache} from "@apollo/client/core";
import {
  ADD_USER_DATA,
  ATTEND_EVENT,
  CREATE_EVENT,
  EDIT_EVENT,
  GET_EVENT,
  GET_EVENTS,
  REMOVE_EVENT,
  SUBSCRIBE_EVENT,
  USER_DATA_EXISTS
} from "./gql-request-strings";
import {Event} from "../model/event";
import {EventDate} from "../event/event.component";
import {UserData} from "../model/userData";
import {onError} from "@apollo/client/link/error";
import {Subject} from "rxjs";
import {environment} from "../../environments/environment";

@Injectable({
  providedIn: 'root'
})
export class GQLService {
  private apollo: ApolloBase;
  errorState: Subject<string> = new Subject<string>();

  constructor(private apolloProvider: Apollo,
              private httpLink: HttpLink) {
    let errorLink = onError(({graphQLErrors, networkError }) => {
      if (networkError) {
        let msg = `Event backend is currently offline, try again later!`;
        this.errorState.next(msg);
      }
    });

    const http = httpLink.create({
      uri: environment.eventAPI,
    });

    try {
      this.apolloProvider.createNamed('event', {
        cache: new InMemoryCache(),
        link: from([errorLink, http]),
      });
    } catch (e) {
      console.error('Error when creating apollo client \'event\'', e);
    }

    this.apollo = this.apolloProvider.use('event');
  }

  createEvent(title: string, host: string, startDate: string, endDate: string, location: string = '', description: string = '') {
    return this.apollo.mutate({
      mutation: CREATE_EVENT,
      variables: {
        title: title,
        host: host,
        description: description,
        startDate: startDate,
        endDate: endDate,
        location: location,
      }
    })
  }

  getEvents(username: string) {
    return this.apollo.watchQuery({
      query: GET_EVENTS,
      variables: {
        username: username,
      }
    }).valueChanges
  }

  removeEvent(eventID: string) {
    return this.apollo.mutate({
        mutation: REMOVE_EVENT,
        variables: {
          remove: eventID,
        }
      }
    );
  }

  editEvent(event: Event, eventDate: EventDate) {
    return this.apollo.mutate({
        mutation: EDIT_EVENT,
        variables: {
          eventID: event.id,
          title: event.title,
          description: event.description,
          subscribers: event.subscribers,
          startDate: eventDate.startDate,
          endDate: eventDate.endDate,
          location: event.location,
          attendees: event.attendees
        }
      }
    );
  }

  subscribeEvent(event: Event) {
    return this.apollo.mutate({
        mutation: SUBSCRIBE_EVENT,
        variables: {
          eventID: event.id,
          title: event.title,
          description: event.description,
          subscribers: event.subscribers,
          startDate: event.startDate,
          endDate: event.endDate,
          location: event.location,
          attendees: event.attendees,
        }
      }
    );
  }

  attendEvent(event: Event, left: boolean, username: string) {
    return this.apollo.mutate({
        mutation: ATTEND_EVENT,
        variables: {
          eventID: event.id,
          title: event.title,
          description: event.description,
          subscribers: event.subscribers,
          startDate: event.startDate,
          endDate: event.endDate,
          location: event.location,
          attendees: event.attendees,
          left: left,
          username: username
        }
      }
    );
  }

  userDataExists(username: string) {
    return this.apollo.mutate({
      mutation: USER_DATA_EXISTS,
      variables: {
        username: username,
      }
    });
  }

  addUserData(userData: UserData) {
    return this.apollo.mutate({
      mutation: ADD_USER_DATA,
      variables: {
        username: userData.username,
        firstname: userData.firstname,
        lastname: userData.lastname,
        street: userData.address.street,
        housenumber: userData.address.housenumber,
        postalcode: userData.address.postalcode,
        city: userData.address.city,
        email: userData.email
      }
    });
  }

  notify(query: DocumentNode, username: string, eventID: string, reportedBy?: string) {
    return this.apollo.query({
      query: query,
      variables: {
        username: username,
        eventID: eventID,
        reportedBy: reportedBy,
      }
    });
  }

  getEventByID(eventID: string) {
    return this.apollo.query({
      query: GET_EVENT,
      variables: {
        eventID: eventID,
      }
    });
  }
}
