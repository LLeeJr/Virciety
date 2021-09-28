
export class Dm {
  __typename: string;
  id: string;
  msg: string;

  constructor(typename: string, id: string, msg: string) {
    this.__typename = typename;
    this.id = id;
    this.msg = msg;
  }
}
