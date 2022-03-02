export class Notification {

  event: string;
  id: string;
  params: {
    key: string,
    value: string,
  }[];
  read: boolean;
  receiver: string;
  route: string;
  text: string;
  timestamp: string;
}
