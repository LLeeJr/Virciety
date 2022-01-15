import {gql} from "apollo-angular";

// ------------------------- Queries, Mutations and Subscriptions

export const GET_EVENTS = gql`
  query getEvents {
    getEvents {
      id
      title
      host
      location
      description
      startDate
      endDate
    }
  }
`

export const CREATE_EVENT = gql`
    mutation createEvent($title: String!, $host: String!, $description: String!, $startDate: String!, $endDate: String!, $location: String!) {
      createEvent(newEvent: {title: $title, host: $host, description: $description, startDate: $startDate, endDate: $endDate, location: $location}) {
        id
        title
        host
        location
        description
        startDate
        endDate
      }
    }
  `;

// ------------------------- Queries, Mutations and Subscriptions end
