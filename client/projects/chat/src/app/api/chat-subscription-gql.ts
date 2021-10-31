import {Injectable} from '@angular/core';
import {Subscription, gql} from 'apollo-angular';


@Injectable({
  providedIn: 'root'
})
export class ChatSubscriptionGqlService extends Subscription {
  document = gql`
  subscription dmAdded($roomName: String!) {
    dmAdded(roomName: $roomName) {
      chatroomId,
      createdAt,
      createdBy,
      msg,
    }
  }
  `;
}
