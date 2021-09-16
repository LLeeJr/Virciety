import {Injectable} from '@angular/core';
import {Subscription, gql} from 'apollo-angular';


@Injectable({
  providedIn: 'root'
})
export class ChatSubscriptionGqlService extends Subscription {
  document = gql`
  subscription dmAdded {
    dmAdded {
      id,
      msg
    }
  }
  `;
}
