import {Injectable} from "@angular/core";
import {gql, Subscription} from "apollo-angular";


@Injectable({
  providedIn: 'root',
})
export class NotificationSubscriptionGql extends Subscription {
  document = gql`
  subscription notifAdded($userName: String!){
    notifAdded(userName: $userName) {
      id,
      receiver,
      text,
      event
    }
  }
  `;
}
