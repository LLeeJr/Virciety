export class Room {
  __typename: string;
  _id: string;
  name: string;
  member: string[];

  constructor(typename: string, id: string, name: string, member: string[]) {
    this.__typename = typename;
    this._id = id;
    this.name = name;
    this.member = member;
  }
}
