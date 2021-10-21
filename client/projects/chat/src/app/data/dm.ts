
export class Dm {
  chatroomId: string;
  createdAt: string;
  createdBy: string;
  msg: string;
  __typename: string;

  constructor(chatroomId: string, createdAt: string, createdBy: string, msg: string, typename: string) {
    this.chatroomId = chatroomId;
    this.createdAt = createdAt;
    this.createdBy = createdBy;
    this.msg = msg;
    this.__typename = typename;
  }
}
