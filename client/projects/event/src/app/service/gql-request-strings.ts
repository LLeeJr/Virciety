import {gql} from "apollo-angular";

// ------------------------- Queries, Mutations and Subscriptions

export const GET_EVENTS = gql`
  query getEvents($username: String) {
    getEvents(username: $username) {
      upcomingEvents {
        id
        title
        host
        location
        description
        startDate
        endDate
        subscribers
        attendees
        currentlyAttended
      }
      ongoingEvents {
        id
        title
        host
        location
        description
        startDate
        endDate
        subscribers
        attendees
        currentlyAttended
      }
      pastEvents {
        id
        title
        host
        location
        description
        startDate
        endDate
        subscribers
        attendees
        currentlyAttended
      }
    }
  }
`

export const GET_EVENT = gql`
    query getEvent($eventID: String!) {
      getEvent(eventID: $eventID) {
        id
        title
        host
        location
        description
        startDate
        endDate
        subscribers
        attendees
        currentlyAttended
      }
    }
  `;

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

export const NOTIFY_HOST_OF_EVENT = gql`
  query notifyHostOfEvent($username: String!, $eventID: String!, $reportedBy: String!) {
    notifyHostOfEvent(username: $username, eventID: $eventID, reportedBy: $reportedBy)
  }
`;

export const NOTIFY_CONTACT_PERSONS = gql`
  query notifyContactPersons($username: String!, $eventID: String!) {
    notifyContactPersons(username: $username, eventID: $eventID)
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
          subscribers
          description
          startDate
          endDate
          attendees
        }
        type
      }
    }
  `;

export const EDIT_EVENT = gql`
    mutation editEvent($eventID: String!, $title: String!, $description: String!, $subscribers: [String!]!, $startDate: String!, $endDate: String!, $location: String!, $attendees: [String!]!) {
      editEvent(edit: {eventID: $eventID, title: $title, description: $description, subscribers: $subscribers, startDate: $startDate, endDate: $endDate, location: $location, attendees: $attendees}) {
        type
      }
    }
  `;

export const SUBSCRIBE_EVENT = gql`
    mutation subscribeEvent($eventID: String!, $title: String!, $description: String!, $subscribers: [String!]!, $startDate: String!, $endDate: String!, $location: String!, $attendees: [String!]!) {
      subscribeEvent(subscribe: {eventID: $eventID, title: $title, description: $description, subscribers: $subscribers, startDate: $startDate, endDate: $endDate, location: $location, attendees: $attendees})
    }
  `;

export const ATTEND_EVENT = gql`
    mutation attendEvent($eventID: String!, $title: String!, $description: String!, $subscribers: [String!]!, $startDate: String!, $endDate: String!, $location: String!, $attendees: [String!]!, $username: String!, $left: Boolean!) {
      attendEvent(attend: {eventID: $eventID, title: $title, description: $description, subscribers: $subscribers, startDate: $startDate, endDate: $endDate, location: $location, attendees: $attendees}, username: $username, left: $left)
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
