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
        attending
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
        attending
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
        attending
      }
    }
  }
`

export const USER_DATA_EXISTS = gql`
  query userDataExists($username: String!) {
    userDataExists(username: $username) {
      username
      firstname
      lastname
      street
      housenumber
      postalcode
      city
      email
    }
  }
`;

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
    mutation editEvent($eventID: String!, $title: String!, $description: String!, $members: [String!]!, $startDate: String!, $endDate: String!, $location: String!, $attending: [String!]!) {
      editEvent(edit: {eventID: $eventID, title: $title, description: $description, members: $members, startDate: $startDate, endDate: $endDate, location: $location, attending: $attending}) {
        type
      }
    }
  `;

export const SUBSCRIBE_EVENT = gql`
    mutation subscribeEvent($eventID: String!, $title: String!, $description: String!, $members: [String!]!, $startDate: String!, $endDate: String!, $location: String!, $attending: [String!]!) {
      subscribeEvent(subscribe: {eventID: $eventID, title: $title, description: $description, members: $members, startDate: $startDate, endDate: $endDate, location: $location, attending: $attending})
    }
  `;

export const ATTEND_EVENT = gql`
    mutation attendEvent($eventID: String!, $title: String!, $description: String!, $members: [String!]!, $startDate: String!, $endDate: String!, $location: String!, $attending: [String!]!, $username: String!, $left: Boolean!) {
      attendEvent(attend: {eventID: $eventID, title: $title, description: $description, members: $members, startDate: $startDate, endDate: $endDate, location: $location, attending: $attending}, username: $username, left: $left)
    }
  `;

export const REMOVE_EVENT = gql`
    mutation removeEvent($remove: String!) {
      removeEvent(remove: $remove)
    }
  `;

export const ADD_USER_DATA = gql`
    mutation addUserData($username: String!, $firstname: String!, $lastname: String!, $street: String!, $housenumber: String!, $postalcode: String!, $city: String!, $email: String!) {
      addUserData(userData: {username: $username, firstname: $firstname, lastname: $lastname, street: $street, housenumber: $housenumber, postalcode: $postalcode, city: $city, email: $email}) {
        username
        firstname
        lastname
        street
        housenumber
        postalcode
        city
        email
      }
    }
  `;

// ------------------------- Queries, Mutations and Subscriptions end
