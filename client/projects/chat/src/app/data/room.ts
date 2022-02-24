export class Room {
  __typename: string;
  id: string;
  name: string;
  member: string[];
  owner: string;
  isDirect: boolean;

  constructor(typename: string, id: string, name: string, member: string[], owner: string, isDirect: boolean) {
    this.__typename = typename;
    this.id = id;
    this.name = name;
    this.member = member;
    this.owner = owner;
    this.isDirect = isDirect;
  }
}
