import {gql} from "apollo-angular";

// ------------------------- Queries, Mutations and Subscriptions

export const GET_EVENTS = gql`
  query getEvents {
    getEvents {
      upcomingEvents {
        id
        title
        host
        location
        description
        startDate
        endDate
        members
      }
      ongoingEvents {
        id
        title
        host
        location
        description
        startDate
        endDate
        members
      }
      pastEvents {
        id
        title
        host
        location
        description
        startDate
        endDate
        members
      }
    }
  }
`

export const CREATE_EVENT = gql`
    mutation createEvent($title: String!, $host: String!, $description: String!, $startDate: String!, $endDate: String!, $location: String!) {
      createEvent(newEvent: {title: $title, host: $host, description: $description, startDate: $startDate, endDate: $endDate, location: $location}) {
        event {
          id
          title
          host
          location
          members
          description
          startDate
          endDate
        }
        type
      }
    }
  `;

export const EDIT_EVENT = gql`
    mutation editEvent($eventID: String!, $title: String!, $description: String!, $members: [String!]!, $startDate: String!, $endDate: String!, $location: String!) {
      editEvent(edit: {eventID: $eventID, title: $title, description: $description, members: $members, startDate: $startDate, endDate: $endDate, location: $location})
    }
  `;

export const SUBSCRIBE_EVENT = gql`
    mutation subscribeEvent($eventID: String!, $title: String!, $description: String!, $members: [String!]!, $startDate: String!, $endDate: String!, $location: String!) {
      subscribeEvent(subscribe: {eventID: $eventID, title: $title, description: $description, members: $members, startDate: $startDate, endDate: $endDate, location: $location})
    }
  `;

export const ATTEND_EVENT = gql`
    mutation attendEvent($eventID: String!, $title: String!, $description: String!, $members: [String!]!, $startDate: String!, $endDate: String!, $location: String!) {
      attendEvent(attend: {eventID: $eventID, title: $title, description: $description, members: $members, startDate: $startDate, endDate: $endDate, location: $location})
    }
  `;

export const REMOVE_EVENT = gql`
    mutation removeEvent($remove: String!) {
      removeEvent(remove: $remove)
    }
  `;

// ------------------------- Queries, Mutations and Subscriptions end
