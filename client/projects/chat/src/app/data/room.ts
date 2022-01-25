export class Room {
  __typename: string;
  id: string;
  name: string;
  member: string[];
  owner: string;

  constructor(typename: string, id: string, name: string, member: string[], owner: string) {
    this.__typename = typename;
    this.id = id;
    this.name = name;
    this.member = member;
    this.owner = owner;
  }
}
