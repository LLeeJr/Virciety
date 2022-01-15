import {gql} from "apollo-angular";

// ------------------------- Queries, Mutations and Subscriptions

export const GET_EVENTS = gql`
  query getEvents {
    getEvents {
      id
      host
      location
      description
      startDate
      endDate
    }
  }
`

export const CREATE_EVENT = gql`
    mutation createEvent($host: String!, $description: String!, $startDate: String!, $endDate: String!, $location: String!) {
      createEvent(newEvent: {host: $host, description: $description, startDate: $startDate, endDate: $endDate, location: $location}) {
        id
        host
        location
        description
        startDate
        endDate
      }
    }
  `;

// ------------------------- Queries, Mutations and Subscriptions end
